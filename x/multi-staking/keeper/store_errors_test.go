package keeper

import (
	"context"
	"errors"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"

	corestore "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type failingStore struct {
	corestore.KVStore
	getErr, hasErr, setErr error
	setCalled              bool
}

func (s *failingStore) Get(key []byte) ([]byte, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.KVStore.Get(key)
}

func (s *failingStore) Has(key []byte) (bool, error) {
	if s.hasErr != nil {
		return false, s.hasErr
	}
	return s.KVStore.Has(key)
}

func (s *failingStore) Set(key, value []byte) error {
	s.setCalled = true
	if s.setErr != nil {
		return s.setErr
	}
	return s.KVStore.Set(key, value)
}

type testStoreService struct{ store corestore.KVStore }

func (s testStoreService) OpenKVStore(context.Context) corestore.KVStore { return s.store }

func TestSetValidatorMultiStakingCoinStoreErrors(t *testing.T) {
	injectedErr := errors.New("store unavailable")
	for _, operation := range []string{"has", "set"} {
		t.Run(operation, func(t *testing.T) {
			db := dbm.NewMemDB()
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			store := &failingStore{KVStore: db}
			if operation == "has" {
				store.hasErr = injectedErr
			} else {
				store.setErr = injectedErr
			}
			k := Keeper{storeService: testStoreService{store}}
			valAddr := sdk.ValAddress(make([]byte, 20))

			err := k.SetValidatorMultiStakingCoin(context.Background(), valAddr, "ario")
			require.ErrorIs(t, err, injectedErr)
			require.Equal(t, operation == "set", store.setCalled, "failed existence checks must not attempt a write")
			stored, err := db.Has(types.GetValidatorMultiStakingCoinKey(valAddr))
			require.NoError(t, err)
			require.False(t, stored)
		})
	}
}

func TestGetParamsStoreError(t *testing.T) {
	injectedErr := errors.New("read unavailable")
	store := &failingStore{getErr: injectedErr}
	k := Keeper{
		storeService: testStoreService{store},
		cdc:          codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
	}
	params, err := k.GetParams(sdk.Context{})
	require.ErrorIs(t, err, injectedErr)
	require.Equal(t, types.Params{}, params)
	require.False(t, store.setCalled)
}

func TestStoreGettersPropagateReadErrors(t *testing.T) {
	injectedErr := errors.New("read unavailable")
	store := &failingStore{getErr: injectedErr}
	k := Keeper{
		storeService: testStoreService{store},
		cdc:          codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
	}
	delAddr := sdk.AccAddress(make([]byte, 20))
	valAddr := sdk.ValAddress(make([]byte, 20))
	lockID := types.MultiStakingLockID(delAddr.String(), valAddr.String())
	unlockID := types.MultiStakingUnlockID(delAddr.String(), valAddr.String())

	tests := []struct {
		name string
		get  func() (bool, error)
	}{
		{
			name: "bond weight",
			get: func() (bool, error) {
				_, found, err := k.GetBondWeight(context.Background(), "ario")
				return found, err
			},
		},
		{
			name: "validator multi staking coin",
			get: func() (bool, error) {
				_, found, err := k.GetValidatorMultiStakingCoin(context.Background(), valAddr)
				return found, err
			},
		},
		{
			name: "multi staking lock",
			get: func() (bool, error) {
				_, found, err := k.GetMultiStakingLock(context.Background(), lockID)
				return found, err
			},
		},
		{
			name: "multi staking unlock",
			get: func() (bool, error) {
				_, found, err := k.GetMultiStakingUnlock(context.Background(), unlockID)
				return found, err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			found, err := tc.get()
			require.ErrorIs(t, err, injectedErr)
			require.False(t, found)
		})
	}
}

func TestStoreGettersReturnDecodeErrors(t *testing.T) {
	db := dbm.NewMemDB()
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	k := Keeper{
		storeService: testStoreService{db},
		cdc:          codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
	}
	delAddr := sdk.AccAddress(make([]byte, 20))
	valAddr := sdk.ValAddress(make([]byte, 20))
	lockID := types.MultiStakingLockID(delAddr.String(), valAddr.String())
	unlockID := types.MultiStakingUnlockID(delAddr.String(), valAddr.String())

	tests := []struct {
		name string
		key  []byte
		get  func() (bool, error)
	}{
		{
			name: "bond weight",
			key:  types.GetBondWeightKey("ario"),
			get: func() (bool, error) {
				_, found, err := k.GetBondWeight(context.Background(), "ario")
				return found, err
			},
		},
		{
			name: "multi staking lock",
			key:  lockID.ToBytes(),
			get: func() (bool, error) {
				_, found, err := k.GetMultiStakingLock(context.Background(), lockID)
				return found, err
			},
		},
		{
			name: "multi staking unlock",
			key:  unlockID.ToBytes(),
			get: func() (bool, error) {
				_, found, err := k.GetMultiStakingUnlock(context.Background(), unlockID)
				return found, err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, db.Set(tc.key, []byte{0xff}))
			found, err := tc.get()
			require.Error(t, err)
			require.False(t, found)
			require.NoError(t, db.Delete(tc.key))
		})
	}
}
