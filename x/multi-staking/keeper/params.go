package keeper

import (
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the x/staking module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	err = store.Set(types.ParamsKey, bz)
	if err != nil {
		return err
	}

	return nil
}

// GetParams gets the x/staking module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params, err
	}

	if bz == nil {
		return params, fmt.Errorf("empty params bytes")
	}

	if err := k.cdc.Unmarshal(bz, &params); err != nil {
		return params, err
	}

	return params, nil
}
