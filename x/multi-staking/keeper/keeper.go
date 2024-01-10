package keeper

import (
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/tendermint/tendermint/libs/log"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	storeKey      storetypes.StoreKey
	memKey        storetypes.StoreKey
	cdc           codec.BinaryCodec
	stakingKeeper stakingkeeper.Keeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	stakingKeeper stakingkeeper.Keeper,
	bankKeeper types.BankKeeper,
	key storetypes.StoreKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      key,
		memKey:        memKey,
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetMatureUnbondingDelegations(ctx sdk.Context) []stakingtypes.UnbondingDelegation {
	var matureUnbondingDelegations []stakingtypes.UnbondingDelegation
	matureUnbonds := k.stakingKeeper.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		delAddr, valAddr, err := types.AccAddrAndValAddrFromStrings(dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		unbondingDelegation, found := k.stakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr) // ??
		if !found {
			continue
		}

		matureUnbondingDelegations = append(matureUnbondingDelegations, unbondingDelegation)
	}
	return matureUnbondingDelegations
}

func (k Keeper) GetUnbondingEntryAtCreationHeight(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, creationHeight int64) (stakingtypes.UnbondingDelegationEntry, bool) {
	ubd, found := k.stakingKeeper.GetUnbondingDelegation(ctx, delAcc, valAcc)
	if !found {
		return stakingtypes.UnbondingDelegationEntry{}, false
	}

	var unbondingEntryAtHeight stakingtypes.UnbondingDelegationEntry
	found = false
	for _, entry := range ubd.Entries {
		if entry.CreationHeight == creationHeight {
			if !found {
				found = true
				unbondingEntryAtHeight = entry
			} else {
				unbondingEntryAtHeight.Balance = unbondingEntryAtHeight.Balance.Add(entry.Balance)
			}
		}
	}

	return unbondingEntryAtHeight, found
}

func (k Keeper) BurnCoin(ctx sdk.Context, accAddr sdk.AccAddress, coin sdk.Coin) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) isValMultiStakingCoin(ctx sdk.Context, valAcc sdk.ValAddress, lockedCoin sdk.Coin) bool {
	return lockedCoin.Denom == k.GetValidatorMultiStakingCoin(ctx, valAcc)
}

func (k Keeper) AdjustUnbondAmount(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, amount math.Int) (adjustedAmount math.Int, err error) {
	delegation, found := k.stakingKeeper.GetDelegation(ctx, delAcc, valAcc)
	if !found {
		return math.Int{}, fmt.Errorf("delegation not found")
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAcc)
	if !found {
		return math.Int{}, fmt.Errorf("validator not found")
	}

	shares, err := validator.SharesFromTokens(amount)
	if err != nil {
		return math.Int{}, err
	}

	delShares := delegation.GetShares()
	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if shares.GT(delShares) {
		shares = delShares
	}

	return validator.TokensFromShares(shares).RoundInt(), nil
}
