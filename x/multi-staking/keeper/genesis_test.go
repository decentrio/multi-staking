package keeper_test

import (
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/realio-tech/multi-staking-module/test/simapp"
)

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	oldCommitHash := suite.app.LastCommitID().Hash
	appState, err := suite.app.ExportAppStateAndValidators(false, []string{})
	suite.NoError(err)

	encConfig := simapp.MakeTestEncodingConfig()

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

	_ = emptyApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   appState.AppState,
		},
	)

	emptyApp.Commit()

	newAppState, err := emptyApp.ExportAppStateAndValidators(false, []string{})
	newCommitHash := emptyApp.LastCommitID().Hash
	suite.NoError(err)

	suite.Equal(appState.AppState, newAppState.AppState)
	suite.Equal(newCommitHash, oldCommitHash)
}
