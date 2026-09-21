package keeper

import (
	"context"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetOrCreateMultiStakingLock(ctx context.Context, lockID types.LockID, denom string) (types.MultiStakingLock, error) {
	multiStakingLock, found, err := k.GetMultiStakingLock(ctx, lockID)
	if err != nil {
		return types.MultiStakingLock{}, err
	}
	if !found {
		multiStakingLock = types.NewMultiStakingLock(lockID, types.MultiStakingCoin{Denom: denom, Amount: math.ZeroInt()})
	}
	return multiStakingLock, nil
}

func (k Keeper) EscrowCoinFrom(ctx context.Context, fromAcc sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, sdk.NewCoins(coin))
}

func (k Keeper) UnescrowCoinTo(ctx context.Context, toAcc sdk.AccAddress, coin sdk.Coin) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	coins := sdk.NewCoins(coin)
	tokenID := k.erc20keeper.GetTokenPairID(sdkCtx, coin.Denom)
	if len(tokenID) == 0 && erc20types.ValidateErc20Denom(coin.Denom) == nil {
		// A previous conversion may already have deleted the dead contract's
		// token pair. Its remaining wrapped coins must not become native payouts.
		return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	}

	err := k.bankKeeper.SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, toAcc, coins)
	if err != nil {
		return err
	}
	// If the coin has an ERC20 token pair, convert it back to ERC20 tokens.
	if len(tokenID) != 0 {
		toAccHex := common.BytesToAddress(toAcc.Bytes()).Hex()
		_, err := k.erc20keeper.ConvertCoin(sdkCtx, &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: toAccHex,
			Sender:   toAcc.String(),
		})
		if err != nil {
			return err
		}
		// ConvertCoin returns nil, nil for both a successful conversion and a
		// selfdestructed contract. Only the latter deletes the pair and leaves
		// the Cosmos coins unburned in the recipient's account.
		if _, found := k.erc20keeper.GetTokenPair(sdkCtx, tokenID); !found {
			if err := k.BurnCoin(sdkCtx, toAcc, coin); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) MintCoin(ctx context.Context, toAcc sdk.AccAddress, coin sdk.Coin) error {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcc, sdk.NewCoins(coin))
	return err
}

func (k Keeper) LockCoinAndMintBondCoin(
	ctx context.Context,
	lockID types.LockID,
	fromAcc sdk.AccAddress,
	mintedTo sdk.AccAddress,
	coin sdk.Coin,
) (mintedBondCoin sdk.Coin, err error) {
	// escrow coin
	err = k.EscrowCoinFrom(ctx, fromAcc, coin)
	if err != nil {
		return sdk.Coin{}, err
	}

	// get multistaking coin's bond weight
	bondWeight, isMultiStakingCoin, err := k.GetBondWeight(ctx, coin.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	if !isMultiStakingCoin {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", coin.Denom,
		)
	}

	// update multistaking lock
	multiStakingCoin := types.NewMultiStakingCoin(coin.Denom, coin.Amount, bondWeight)
	lock, err := k.GetOrCreateMultiStakingLock(ctx, lockID, coin.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	err = lock.AddCoinToMultiStakingLock(multiStakingCoin)
	if err != nil {
		return sdk.Coin{}, err
	}

	k.SetMultiStakingLock(ctx, lock)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking coin * bond coin weight
	mintedBondAmount := multiStakingCoin.BondValue()
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	mintedBondCoin = sdk.NewCoin(bondDenom, mintedBondAmount)

	// mint bond coin to delegator account
	err = k.MintCoin(ctx, mintedTo, mintedBondCoin)
	if err != nil {
		return sdk.Coin{}, err
	}

	return mintedBondCoin, nil
}
