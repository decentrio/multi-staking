package keeper_test

import "github.com/realio-tech/multi-staking-module/x/multi-staking/types"

func (suite *KeeperTestSuite) TestGetParams() {
	store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	store.Delete(types.ParamsKey)
	params, err := suite.msKeeper.GetParams(suite.ctx)
	suite.Require().EqualError(err, "empty params bytes")
	suite.Require().Equal(types.Params{}, params)
	suite.Require().False(store.Has(types.ParamsKey), "reading missing params must not write defaults")

	want := types.Params{MainBondDenom: "ario"}
	suite.Require().NoError(suite.msKeeper.SetParams(suite.ctx, want))
	params, err = suite.msKeeper.GetParams(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(want, params)

	store.Delete(types.ParamsKey)
	params, err = suite.msKeeper.GetParams(suite.ctx)
	suite.Require().EqualError(err, "empty params bytes")
	suite.Require().Equal(types.Params{}, params, "missing params must not retain a previous value")

	store.Set(types.ParamsKey, []byte("invalid params"))
	params, err = suite.msKeeper.GetParams(suite.ctx)
	suite.Require().Error(err)
	suite.Require().Equal(types.Params{}, params)
}
