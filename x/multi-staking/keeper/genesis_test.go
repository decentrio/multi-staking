package keeper_test

import (
	dbm "github.com/cosmos/cosmos-db"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/test/simapp"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (suite *KeeperTestSuite) TestInitGenesisDuplicateValidatorCoinWithoutValidation() {
	valAddr := test.GenValAddress()
	genesis := types.DefaultGenesis()
	genesis.ValidatorMultiStakingCoins = []types.ValidatorMultiStakingCoin{
		{ValAddr: valAddr.String(), CoinDenom: "ario"},
		{ValAddr: valAddr.String(), CoinDenom: "arst"},
	}
	suite.Require().ErrorContains(genesis.Validate(), "duplicate validator multi staking coin")
	// Direct initialization still fails explicitly if validation was bypassed.
	suite.Require().PanicsWithError("validator multi staking coin already set: "+valAddr.String(), func() {
		suite.msKeeper.InitGenesis(suite.ctx, *genesis)
	})
	suite.Require().Equal("ario", suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, valAddr))
}

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	appState, err := suite.app.ExportAppStateAndValidators(false, []string{}) //nolint:staticcheck
	suite.NoError(err)

	encConfig := simapp.MakeEncodingConfig() //nolint:staticcheck

	configurator := evmtypes.NewEVMConfigurator()
	configurator.ResetTestConfig()

	emptyApp := simapp.NewSimApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"temp",
		simapp.FlagPeriodValue,
		encConfig,
		simapp.EmptyAppOptions{},
	)

	_, err = emptyApp.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   appState.AppState,
		},
	)
	suite.NoError(err)
	_, err = emptyApp.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: emptyApp.LastBlockHeight() + 1,
		Hash:   emptyApp.LastCommitID().Hash,
	})
	suite.NoError(err)

	_, err = emptyApp.ExportAppStateAndValidators(false, []string{})
	suite.NoError(err)
}
