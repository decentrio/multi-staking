package keeper

import (
	"bytes"
	"context"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetOrCreateMultiStakingLock(ctx context.Context, lockID types.LockID, denom string) types.MultiStakingLock {
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(lockID, types.MultiStakingCoin{Denom: denom, Amount: math.ZeroInt()})
	}
	return multiStakingLock
}

func (k Keeper) EscrowCoinFrom(ctx context.Context, fromAcc sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, sdk.NewCoins(coin))
}

func (k Keeper) UnescrowCoinTo(ctx context.Context, toAcc sdk.AccAddress, coin sdk.Coin) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcc, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	// Conversion is optional: a failed ERC20 transfer must not block unbonding.
	// Keep the Cosmos payout outside the cache so the recipient retains it.
	tokenId := k.erc20keeper.GetTokenPairID(sdkCtx, coin.Denom)
	if !bytes.Equal(tokenId, []byte{}) {
		conversionCtx, commit := sdkCtx.CacheContext()
		// CacheContext shares its gas meter. Cosmos EVM consumes the whole
		// limit on revert, which would poison EndBlock's infinite meter.
		conversionCtx = conversionCtx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		toAccHex := common.BytesToAddress(toAcc.Bytes()).Hex()
		_, err := k.erc20keeper.ConvertCoin(conversionCtx, &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: toAccHex,
			Sender:   toAcc.String(),
		})
		if err != nil {
			k.Logger(ctx).Error("ERC20 conversion failed; unbonded Cosmos coins retained", "recipient", toAcc.String(), "coin", coin.String(), "error", err)
			return nil
		}
		sdkCtx.GasMeter().ConsumeGas(conversionCtx.GasMeter().GasConsumed(), "ERC20 unbonding conversion")
		commit()
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
	bondWeight, isMultiStakingCoin := k.GetBondWeight(ctx, coin.Denom)
	if !isMultiStakingCoin {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", coin.Denom,
		)
	}

	// update multistaking lock
	multiStakingCoin := types.NewMultiStakingCoin(coin.Denom, coin.Amount, bondWeight)
	lock := k.GetOrCreateMultiStakingLock(ctx, lockID, coin.Denom)
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
