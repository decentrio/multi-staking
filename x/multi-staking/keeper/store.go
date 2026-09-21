package keeper

import (
	"context"
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBondWeight(ctx context.Context, tokenDenom string) (math.LegacyDec, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetBondWeightKey(tokenDenom))
	if err != nil {
		return math.LegacyDec{}, false, fmt.Errorf("get bond weight for %q: %w", tokenDenom, err)
	}
	if bz == nil {
		return math.LegacyDec{}, false, nil
	}

	bondCoinWeight := &math.LegacyDec{}
	err = bondCoinWeight.Unmarshal(bz)
	if err != nil {
		return math.LegacyDec{}, false, fmt.Errorf("unable to unmarshal bond coin weight: %w", err)
	}
	return *bondCoinWeight, true, nil
}

func (k Keeper) SetBondWeight(ctx context.Context, tokenDenom string, tokenWeight math.LegacyDec) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := tokenWeight.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal bond coin weight %v", err))
	}

	err = store.Set(types.GetBondWeightKey(tokenDenom), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveBondWeight(ctx context.Context, tokenDenom string) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(types.GetBondWeightKey(tokenDenom))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetValidatorMultiStakingCoin(ctx context.Context, operatorAddr sdk.ValAddress) (string, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetValidatorMultiStakingCoinKey(operatorAddr))
	if err != nil {
		return "", false, fmt.Errorf("get multi-staking coin for validator %s: %w", operatorAddr, err)
	}
	if bz == nil {
		return "", false, nil
	}

	return string(bz), true, nil
}

func (k Keeper) SetValidatorMultiStakingCoin(ctx context.Context, operatorAddr sdk.ValAddress, bondDenom string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorMultiStakingCoinKey(operatorAddr)
	exists, err := store.Has(key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("validator multi staking coin already set: %s", operatorAddr)
	}

	return store.Set(key, []byte(bondDenom))
}

func (k Keeper) ValidatorMultiStakingCoinIterator(ctx context.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.ValidatorMultiStakingCoinKey)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Key()).String()
		denom := string(iterator.Value())
		if cb(valAddr, denom) {
			break
		}
	}
}

func (k Keeper) GetMultiStakingLock(ctx context.Context, multiStakingLockID types.LockID) (types.MultiStakingLock, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(multiStakingLockID.ToBytes())
	if err != nil {
		return types.MultiStakingLock{}, false, fmt.Errorf("get multi-staking lock: %w", err)
	}

	if bz == nil {
		return types.MultiStakingLock{}, false, nil
	}

	multiStakingLock := types.MultiStakingLock{}
	if err := k.cdc.Unmarshal(bz, &multiStakingLock); err != nil {
		return types.MultiStakingLock{}, false, fmt.Errorf("unmarshal multi-staking lock: %w", err)
	}
	return multiStakingLock, true, nil
}

func (k Keeper) SetMultiStakingLock(ctx context.Context, multiStakingLock types.MultiStakingLock) {
	if multiStakingLock.IsEmpty() {
		k.RemoveMultiStakingLock(ctx, multiStakingLock.LockID)
		return
	}

	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	err := store.Set(multiStakingLock.LockID.ToBytes(), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveMultiStakingLock(ctx context.Context, multiStakingLockID types.LockID) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(multiStakingLockID.ToBytes())
	if err != nil {
		panic(err)
	}
}

func (k Keeper) MultiStakingLockIterator(ctx context.Context, cb func(stakingLock types.MultiStakingLock) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var multiStakingLock types.MultiStakingLock
		k.cdc.MustUnmarshal(iterator.Value(), &multiStakingLock)
		if cb(multiStakingLock) {
			break
		}
	}
}

func (k Keeper) MultiStakingUnlockIterator(ctx context.Context, cb func(multiStakingUnlock types.MultiStakingUnlock) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.MultiStakingUnlockPrefix)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var multiStakingUnlock types.MultiStakingUnlock
		k.cdc.MustUnmarshal(iterator.Value(), &multiStakingUnlock)
		if cb(multiStakingUnlock) {
			break
		}
	}
}

func (k Keeper) BondWeightIterator(ctx context.Context, cb func(denom string, bondWeight math.LegacyDec) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.BondWeightKey)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key())
		bondWeight := &math.LegacyDec{}
		err := bondWeight.Unmarshal(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("unable to unmarshal bond coin weight %v", err))
		}
		if cb(denom, *bondWeight) {
			break
		}
	}
}

func (k Keeper) GetMultiStakingUnlock(ctx context.Context, multiStakingUnlockID types.UnlockID) (unlock types.MultiStakingUnlock, found bool, err error) {
	store := k.storeService.OpenKVStore(ctx)
	value, err := store.Get(multiStakingUnlockID.ToBytes())
	if err != nil {
		return unlock, false, fmt.Errorf("get multi-staking unlock: %w", err)
	}

	if value == nil {
		return unlock, false, nil
	}

	unlock = types.MultiStakingUnlock{}
	if err := k.cdc.Unmarshal(value, &unlock); err != nil {
		return types.MultiStakingUnlock{}, false, fmt.Errorf("unmarshal multi-staking unlock: %w", err)
	}

	return unlock, true, nil
}

// SetMultiStakingUnlock sets the unbonding delegation and associated index.
func (k Keeper) SetMultiStakingUnlock(ctx context.Context, unlock types.MultiStakingUnlock) {
	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&unlock)

	err := store.Set(unlock.UnlockID.ToBytes(), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) DeleteMultiStakingUnlock(ctx context.Context, unlockID types.UnlockID) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(unlockID.ToBytes())
	if err != nil {
		panic(err)
	}
}
