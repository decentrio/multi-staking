package keeper_test

import "github.com/realio-tech/multi-staking-module/x/multi-staking/types"

func (suite *KeeperTestSuite) TestGetParams() {
	store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	store.Delete(types.ParamsKey)
	suite.Require().Equal(types.Params{}, suite.msKeeper.GetParams(suite.ctx))
	suite.Require().False(store.Has(types.ParamsKey), "reading missing params must not write defaults")

	params := types.Params{MainBondDenom: "ario"}
	suite.Require().NoError(suite.msKeeper.SetParams(suite.ctx, params))
	suite.Require().Equal(params, suite.msKeeper.GetParams(suite.ctx))

	store.Delete(types.ParamsKey)
	suite.Require().Equal(types.Params{}, suite.msKeeper.GetParams(suite.ctx), "missing params must not retain a previous value")
}
