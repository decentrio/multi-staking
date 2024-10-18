package keeper_test

import (
	"encoding/json"
	"fmt"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/realio-tech/multi-staking-module/test/simapp"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	suite.app.Commit()
	appState, err := suite.app.ExportAppStateAndValidators(false, []string{}, []string{})
	suite.NoError(err)

	encConfig := simapp.MakeTestEncodingConfig()

	// NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encCdc, EmptyAppOptions{})
	emptyApp := simapp.NewSimApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"temp",
		5,
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
	// a, _ := suite.app.CrisisKeeper.ConstantFee.Get(suite.ctx)
	// fmt.Println("a :", a)
	// err = emptyApp.CrisisKeeper.ConstantFee.Set(suite.ctx, a)

	// _, err = emptyApp.FinalizeBlock(&abci.RequestFinalizeBlock{
	// 	Height: emptyApp.LastBlockHeight() + 1,
	// })
	// suite.NoError(err)

	// _, err = emptyApp.Commit()
	// suite.NoError(err)

	ctx2 := emptyApp.NewContextLegacy(true, tmproto.Header{Height: emptyApp.LastBlockHeight()})

	fmt.Println("block123 ", emptyApp.LastBlockHeight())
	fmt.Println("block1234 ", suite.app.LastBlockHeight())
	c, _ := emptyApp.CrisisKeeper.ConstantFee.Get(ctx2)
	fmt.Println(" ggggg1 ", c)

	fmt.Println(" dasda ", emptyApp.LastBlockHeight())

	newAppState, err := emptyApp.ExportAppStateAndValidators2(ctx2, false, []string{}, []string{})
	suite.NoError(err)
	var genesisState simapp.GenesisState
	err = json.Unmarshal(appState.AppState, &genesisState)
	suite.NoError(err)

	var genesisState2 simapp.GenesisState
	err = json.Unmarshal(newAppState.AppState, &genesisState2)
	suite.NoError(err)

	for key, value := range genesisState {
		b := genesisState2[key]
		if b == nil {
			fmt.Println("key==", key)
			continue
		}

		if len(value) != len(b) {
			fmt.Println("key==", key)
			continue
		}

		for i := range value {
			if value[i] != b[i] {
				fmt.Println("key==", key)
				break
			}
		}
	}

	fmt.Println("hash1 ", suite.app.LastCommitID().Hash)
	fmt.Println("hash2 ", emptyApp.LastCommitID().Hash)

	var bankGenesis banktypes.GenesisState
	err = json.Unmarshal(genesisState[banktypes.ModuleName], &bankGenesis)
	suite.NoError(err)
	fmt.Println("bankGenesis1  :", bankGenesis)

	var bankGenesis2 banktypes.GenesisState
	err = json.Unmarshal(genesisState2[banktypes.ModuleName], &bankGenesis2)
	suite.NoError(err)
	fmt.Println("bankGenesis2  :", bankGenesis2)

	// var mintGenesis minttypes.GenesisState
	// err = json.Unmarshal(genesisState[minttypes.ModuleName], &mintGenesis)
	// suite.NoError(err)
	// fmt.Println("mintGenesis1  :", mintGenesis)

	// var mintGenesis2 minttypes.GenesisState
	// err = json.Unmarshal(genesisState2[minttypes.ModuleName], &mintGenesis2)
	// suite.NoError(err)
	// fmt.Println("mintGenesis2  :", mintGenesis2)

	var crisis1 crisistypes.GenesisState
	err = json.Unmarshal(genesisState2[crisistypes.ModuleName], &crisis1)
	suite.NoError(err)
	fmt.Println("crisis1  :", crisis1)

	fmt.Println("abc  :")
	// fmt.Println("genesisState2  :", genesisState2)

	suite.Equal(appState.AppState, newAppState.AppState)
}
