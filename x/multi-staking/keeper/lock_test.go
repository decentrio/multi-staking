package keeper_test

import (
	"context"
	"errors"
	"math/big"

	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/realio-tech/multi-staking-module/test"
	keeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Use the real ERC20 conversion and bank keepers; only EVM calls are simulated.
type unescrowEVM struct {
	erc20types.EVMKeeper
	account *statedb.Account
	balance big.Int
}

func (e *unescrowEVM) GetAccountWithoutBalance(sdk.Context, common.Address) *statedb.Account {
	return e.account
}

func (e *unescrowEVM) CallEVM(_ sdk.Context, contractABI abi.ABI, _, _ common.Address, _ bool, _ *big.Int, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	var value interface{}
	switch method {
	case "balanceOf":
		value = &e.balance
	case "transfer":
		e.balance.Add(&e.balance, args[1].(*big.Int))
		value = true
	default:
		panic("unexpected EVM method: " + method)
	}
	ret, err := contractABI.Methods[method].Outputs.Pack(value)
	return &evmtypes.MsgEthereumTxResponse{Ret: ret}, err
}

type unescrowFailingBank struct {
	types.BankKeeper
	burnErr   error
	returnErr error
}

func (b unescrowFailingBank) BurnCoins(ctx context.Context, module string, coins sdk.Coins) error {
	if b.burnErr != nil {
		return b.burnErr
	}
	return b.BankKeeper.BurnCoins(ctx, module, coins)
}

func (b unescrowFailingBank) SendCoinsFromAccountToModule(ctx context.Context, sender sdk.AccAddress, module string, coins sdk.Coins) error {
	if b.returnErr != nil {
		return b.returnErr
	}
	return b.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, module, coins)
}

func (suite *KeeperTestSuite) unescrowKeepers(evm *unescrowEVM, bank types.BankKeeper) (*keeper.Keeper, erc20keeper.Keeper) {
	erc20 := erc20keeper.NewKeeper(
		suite.app.GetKey(erc20types.StoreKey), suite.app.AppCodec(),
		authtypes.NewModuleAddress(govtypes.ModuleName), suite.app.AccountKeeper,
		suite.app.BankKeeper, evm, suite.app.StakingKeeper, nil,
	)
	ms := keeper.NewKeeper(suite.app.AppCodec(), suite.app.AccountKeeper,
		suite.app.StakingKeeper, bank, erc20, nil, "", nil, nil)
	return ms, erc20
}

func (suite *KeeperTestSuite) TestUnescrowCoinTo() {
	contract := common.HexToAddress("0x1234567890123456789012345678901234567890")
	denom := erc20types.CreateDenom(contract.Hex())
	liveAccount := &statedb.Account{CodeHash: common.HexToHash("0x01").Bytes()}
	for _, tc := range []struct {
		name       string
		denom      string
		registered bool
		account    *statedb.Account
		burned     bool
		pairExists bool
	}{
		{"native coin", "unative", false, nil, false, false},
		{"missing contract", denom, true, nil, true, false},
		{"contract without code", denom, true, statedb.NewEmptyAccount(), true, false},
		{"pair already deleted", denom, false, nil, true, false},
		{"live contract returns nil response", denom, true, liveAccount, true, true},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			evm := &unescrowEVM{account: tc.account}
			ms, erc20 := suite.unescrowKeepers(evm, suite.app.BankKeeper)
			pair := erc20types.NewTokenPair(contract, tc.denom, erc20types.OWNER_EXTERNAL)
			if tc.registered {
				suite.Require().NoError(erc20.SetToken(suite.ctx, pair))
			}
			coin := sdk.NewInt64Coin(tc.denom, 100)
			initial := sdk.NewInt64Coin(tc.denom, 7)
			recipient := suite.CreateAndFundAccount(sdk.NewCoins(initial.Add(coin)))
			suite.Require().NoError(ms.EscrowCoinFrom(suite.ctx, recipient, coin))
			supplyBefore := suite.app.BankKeeper.GetSupply(suite.ctx, tc.denom)

			suite.Require().NoError(ms.UnescrowCoinTo(suite.ctx, recipient, coin))

			expectedBalance, expectedSupply := initial.Add(coin), supplyBefore
			if tc.burned {
				expectedBalance, expectedSupply = initial, supplyBefore.Sub(coin)
			}
			suite.Require().Equal(expectedBalance, suite.app.BankKeeper.GetBalance(suite.ctx, recipient, tc.denom))
			suite.Require().Equal(expectedSupply, suite.app.BankKeeper.GetSupply(suite.ctx, tc.denom))
			for _, module := range []string{types.ModuleName, erc20types.ModuleName} {
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, authtypes.NewModuleAddress(module), tc.denom)
				suite.Require().True(balance.IsZero(), balance.String())
			}
			_, exists := erc20.GetTokenPair(suite.ctx, pair.GetID())
			suite.Require().Equal(tc.pairExists, exists)
			if tc.pairExists {
				suite.Require().Equal(coin.Amount.BigInt(), &evm.balance)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUnescrowDeadERC20Repeated() {
	ms, erc20 := suite.unescrowKeepers(&unescrowEVM{}, suite.app.BankKeeper)
	contract := common.HexToAddress("0x1234567890123456789012345678901234567890")
	denom := erc20types.CreateDenom(contract.Hex())
	pair := erc20types.NewTokenPair(contract, denom, erc20types.OWNER_EXTERNAL)
	suite.Require().NoError(erc20.SetToken(suite.ctx, pair))
	coin := sdk.NewInt64Coin(denom, 100)
	deposit := suite.CreateAndFundAccount(sdk.NewCoins(coin.Add(coin)))
	suite.Require().NoError(ms.EscrowCoinFrom(suite.ctx, deposit, coin.Add(coin)))
	for i := int64(1); i >= 0; i-- {
		recipient := test.GenAddress()
		suite.Require().NoError(ms.UnescrowCoinTo(suite.ctx, recipient, coin))
		suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, recipient, denom).IsZero())
		suite.Require().Equal(sdk.NewInt64Coin(denom, i*100), suite.app.BankKeeper.GetSupply(suite.ctx, denom))
		suite.Require().Empty(erc20.GetTokenPairID(suite.ctx, denom))
	}
}

func (suite *KeeperTestSuite) TestUnescrowCoinToErrors() {
	injectedErr := errors.New("injected bank failure")
	for _, tc := range []struct {
		name       string
		registered bool
		disabled   bool
		burnErr    error
		returnErr  error
	}{
		{"burn failure after deletion", true, false, injectedErr, nil},
		{"return to module failure after deletion", true, false, nil, injectedErr},
		{"burn failure with deleted pair", false, false, injectedErr, nil},
		{"conversion disabled", true, true, nil, nil},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			bank := unescrowFailingBank{BankKeeper: suite.app.BankKeeper, burnErr: tc.burnErr, returnErr: tc.returnErr}
			ms, erc20 := suite.unescrowKeepers(&unescrowEVM{}, bank)
			contract := common.HexToAddress("0x1234567890123456789012345678901234567890")
			denom := erc20types.CreateDenom(contract.Hex())
			pair := erc20types.NewTokenPair(contract, denom, erc20types.OWNER_EXTERNAL)
			pair.Enabled = !tc.disabled
			if tc.registered {
				suite.Require().NoError(erc20.SetToken(suite.ctx, pair))
			}
			coin := sdk.NewInt64Coin(denom, 100)
			recipient := suite.CreateAndFundAccount(sdk.NewCoins(coin))
			suite.Require().NoError(suite.msKeeper.EscrowCoinFrom(suite.ctx, recipient, coin))
			err := ms.UnescrowCoinTo(suite.ctx, recipient, coin)
			if tc.disabled {
				suite.Require().ErrorIs(err, erc20types.ErrERC20TokenPairDisabled)
			} else {
				suite.Require().ErrorIs(err, injectedErr)
			}
		})
	}
}
