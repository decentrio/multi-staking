package keeper_test

import (
	"math/big"

	"github.com/cosmos/evm/contracts"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Exercise the application-wired bank, ERC20, and EVM keepers without mocks.
func (suite *KeeperTestSuite) TestUnescrowCoinToBurnIntegration() {
	ctx := suite.ctx
	validators, err := suite.app.StakingKeeper.GetAllValidators(ctx)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(validators)
	proposer, err := validators[0].GetConsAddr()
	suite.Require().NoError(err)
	header := ctx.BlockHeader()
	header.ProposerAddress = proposer
	ctx = ctx.WithBlockHeader(header)
	erc20 := suite.app.ERC20Keeper
	evm := suite.app.EVMKeeper
	bank := suite.app.BankKeeper
	owner := suite.CreateAndFundAccount(sdk.NewCoins(sdk.NewInt64Coin("aatom", 1)))
	ownerHex := common.BytesToAddress(owner)
	compiled := contracts.ERC20MinterBurnerDecimalsContract
	args, err := compiled.ABI.Pack("", "Unescrow integration token", "UIT", uint8(18))
	suite.Require().NoError(err)
	bytecode := append(append([]byte{}, compiled.Bin...), args...)
	nonce, err := suite.app.AccountKeeper.GetSequence(ctx, owner)
	suite.Require().NoError(err)
	contract := crypto.CreateAddress(ownerHex, nonce)
	_, err = evm.CallEVMWithData(ctx, ownerHex, nil, bytecode, true, nil)
	suite.Require().NoError(err)
	suite.Require().True(evm.IsContract(ctx, contract))

	_, err = erc20.RegisterERC20(ctx, &erc20types.MsgRegisterERC20{
		Signer:         authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Erc20Addresses: []string{contract.Hex()},
	})
	suite.Require().NoError(err)
	denom := erc20types.CreateDenom(contract.Hex())
	pairID := erc20.GetTokenPairID(ctx, denom)
	suite.Require().NotEmpty(pairID)

	// Produce wrapped coins through an actual ERC20 mint and conversion.
	_, err = evm.CallEVM(ctx, compiled.ABI, ownerHex, contract, true, nil, "mint", ownerHex, big.NewInt(307))
	suite.Require().NoError(err)
	_, err = erc20.ConvertERC20(ctx, &erc20types.MsgConvertERC20{
		ContractAddress: contract.Hex(), Sender: ownerHex.Hex(), Receiver: owner.String(), Amount: math.NewInt(307),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewInt64Coin(denom, 307), bank.GetSupply(ctx, denom))
	suite.Require().Equal(big.NewInt(307), erc20.BalanceOf(ctx, compiled.ABI, contract, erc20types.ModuleAddress))
	suite.Require().NoError(suite.msKeeper.EscrowCoinFrom(ctx, owner, sdk.NewInt64Coin(denom, 300)))

	withdrawal := sdk.NewInt64Coin(denom, 100)
	assertBalances := func(remainingEscrow int64) {
		suite.T().Helper()
		suite.Require().Equal(sdk.NewInt64Coin(denom, 7), bank.GetBalance(ctx, owner, denom), "preserve unrelated wrapped coins")
		suite.Require().Equal(sdk.NewInt64Coin(denom, remainingEscrow), bank.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), denom))
		suite.Require().True(bank.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), denom).IsZero())
		suite.Require().Equal(sdk.NewInt64Coin(denom, remainingEscrow+7), bank.GetSupply(ctx, denom))
	}
	assertBalances(300)

	// A live contract returns tokens and burns the wrapped coins exactly once.
	txCtx, commit := ctx.CacheContext()
	suite.Require().NoError(suite.msKeeper.UnescrowCoinTo(txCtx, owner, withdrawal))
	commit()
	assertBalances(200)
	suite.Require().Equal(big.NewInt(100), erc20.BalanceOf(ctx, compiled.ABI, contract, ownerHex))
	_, found := erc20.GetTokenPair(ctx, pairID)
	suite.Require().True(found, "successful conversion must retain the pair")

	// Model a disappeared contract using the real EVM account deletion path.
	// This fixture does not depend on SELFDESTRUCT opcode rules for a given fork.
	suite.Require().NoError(evm.DeleteAccount(ctx, contract))
	suite.Require().Nil(evm.GetAccountWithoutBalance(ctx, contract))
	_, found = erc20.GetTokenPair(ctx, pairID)
	suite.Require().True(found, "the stale pair must exist before the first burn")

	for _, remainingEscrow := range []int64{100, 0} {
		// The first withdrawal deletes the stale pair and burns the payout.
		// The second burns directly from escrow because the pair is already gone.
		txCtx, commit := ctx.CacheContext()
		suite.Require().NoError(suite.msKeeper.UnescrowCoinTo(txCtx, owner, withdrawal))
		commit()
		assertBalances(remainingEscrow)
		_, found = erc20.GetTokenPair(ctx, pairID)
		suite.Require().False(found)
		suite.Require().Empty(erc20.GetTokenPairID(ctx, denom))
		suite.Require().Empty(erc20.GetTokenPairID(ctx, contract.Hex()))
	}
}
