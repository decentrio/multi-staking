package keeper_test

import (
	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *KeeperTestSuite) TestMatureUnbondingLookupErrors() {
	del, val := test.GenAddress(), test.GenValAddress()
	pairs := []stakingtypes.DVPair{{DelegatorAddress: del.String(), ValidatorAddress: val.String()}}
	suite.Require().NoError(suite.app.StakingKeeper.SetUBDQueueTimeSlice(suite.ctx, suite.ctx.BlockTime(), pairs))
	// Stale queue entries are harmless and must be skipped.
	unbonds, err := suite.msKeeper.GetMatureUnbondingDelegations(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Empty(unbonds)
	// A decode error must be returned, never appended as a zero-value delegation.
	store := suite.ctx.KVStore(suite.app.GetKey(stakingtypes.StoreKey))
	store.Set(stakingtypes.GetUBDKey(del, val), []byte{0xff})
	unbonds, err = suite.msKeeper.GetMatureUnbondingDelegations(suite.ctx)
	suite.Require().Error(err)
	suite.Require().Nil(unbonds)
}

func (suite *KeeperTestSuite) TestNewLockDenomAndEmptyDeletion() {
	id := types.MultiStakingLockID(test.GenAddress().String(), test.GenValAddress().String())
	lock, err := suite.msKeeper.GetOrCreateMultiStakingLock(suite.ctx, id, "ario")
	suite.Require().NoError(err)
	coin := types.NewMultiStakingCoin("ario", math.NewInt(10), math.LegacyOneDec())
	suite.Require().NoError(lock.AddCoinToMultiStakingLock(coin))
	suite.msKeeper.SetMultiStakingLock(suite.ctx, lock)
	// An existing lock must keep its denomination, regardless of the requested default.
	existing, err := suite.msKeeper.GetOrCreateMultiStakingLock(suite.ctx, id, "arst")
	suite.Require().NoError(err)
	suite.Require().Equal("ario", existing.LockedCoin.Denom)
	suite.Require().Error(existing.AddCoinToMultiStakingLock(types.NewMultiStakingCoin("arst", math.NewInt(1), math.LegacyOneDec())))
	suite.Require().NoError(lock.RemoveCoinFromMultiStakingLock(coin))
	suite.msKeeper.SetMultiStakingLock(suite.ctx, lock)
	_, found, err := suite.msKeeper.GetMultiStakingLock(suite.ctx, id)
	suite.Require().NoError(err)
	suite.Require().False(found)
}
