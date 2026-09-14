package keeper_test

import (
	"encoding/hex"
	"time"

	"github.com/cosmos/evm/contracts"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/vm/statedb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	lock := suite.msKeeper.GetOrCreateMultiStakingLock(suite.ctx, id, "ario")
	coin := types.NewMultiStakingCoin("ario", math.NewInt(10), math.LegacyOneDec())
	suite.Require().NoError(lock.AddCoinToMultiStakingLock(coin))
	suite.msKeeper.SetMultiStakingLock(suite.ctx, lock)
	// An existing lock must keep its denomination, regardless of the requested default.
	existing := suite.msKeeper.GetOrCreateMultiStakingLock(suite.ctx, id, "arst")
	suite.Require().Equal("ario", existing.LockedCoin.Denom)
	suite.Require().Error(existing.AddCoinToMultiStakingLock(types.NewMultiStakingCoin("arst", math.NewInt(1), math.LegacyOneDec())))
	suite.Require().NoError(lock.RemoveCoinFromMultiStakingLock(coin))
	suite.msKeeper.SetMultiStakingLock(suite.ctx, lock)
	_, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, id)
	suite.Require().False(found)
}

// Exercise the real ERC20/EVM keepers through EndBlock. The contract returns a
// valid balanceOf result but reverts on transfer, after ConvertCoin has escrowed
// the recipient's Cosmos coins in the ERC20 module.
func (suite *KeeperTestSuite) TestEndBlockERC20RevertRetainsPayout() {
	vals, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	suite.Require().NoError(err)
	val := vals[0]
	consAddr, err := val.GetConsAddr()
	suite.Require().NoError(err)
	header := suite.ctx.BlockHeader()
	header.ProposerAddress = consAddr
	suite.ctx = suite.ctx.WithBlockHeader(header)
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	suite.Require().NoError(err)
	denom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, valAddr)
	coin := sdk.NewInt64Coin(denom, 100)
	del := suite.CreateAndFundAccount(sdk.NewCoins(coin))
	_, err = suite.msgServer.Delegate(suite.ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: del.String(), ValidatorAddress: val.OperatorAddress, Amount: coin,
	})
	suite.Require().NoError(err)
	suite.ctx = suite.ctx.WithBlockHeight(2).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
	_, err = suite.msgServer.Undelegate(suite.ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: del.String(), ValidatorAddress: val.OperatorAddress, Amount: coin,
	})
	suite.Require().NoError(err)
	ubd, err := suite.app.StakingKeeper.GetUnbondingDelegation(suite.ctx, del, valAddr)
	suite.Require().NoError(err)

	contract := common.HexToAddress("0x1000000000000000000000000000000000000001")
	// if calldata selector == balanceOf: return uint256(0); otherwise revert.
	code, err := hex.DecodeString("60003560e01c6370a082311460145760006000fd5b600060005260206000f3")
	suite.Require().NoError(err)
	hash := crypto.Keccak256(code)
	suite.Require().NoError(suite.app.EVMKeeper.SetAccount(suite.ctx, contract, statedb.Account{CodeHash: hash, Balance: uint256.NewInt(0)}))
	suite.app.EVMKeeper.SetCode(suite.ctx, hash, code)
	pair := erc20types.NewTokenPair(contract, denom, erc20types.OWNER_EXTERNAL)
	suite.Require().NoError(suite.app.ERC20Keeper.SetToken(suite.ctx, pair))
	// Confirm the fixture reaches the transfer, rather than failing balanceOf.
	_, err = suite.app.EVMKeeper.CallEVM(suite.ctx, contracts.ERC20MinterBurnerDecimalsContract.ABI, erc20types.ModuleAddress, contract, false, nil, "balanceOf", common.BytesToAddress(del))
	suite.Require().NoError(err)
	balance := suite.app.ERC20Keeper.BalanceOf(suite.ctx, contracts.ERC20MinterBurnerDecimalsContract.ABI, contract, common.BytesToAddress(del))
	suite.Require().NotNil(balance)
	suite.Require().Zero(balance.Sign())
	probeCtx, _ := suite.ctx.CacheContext()
	probeCtx = probeCtx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	_, err = suite.app.EVMKeeper.CallEVM(probeCtx, contracts.ERC20MinterBurnerDecimalsContract.ABI, erc20types.ModuleAddress, contract, true, nil, "transfer", common.BytesToAddress(del), coin.Amount.BigInt())
	suite.Require().ErrorContains(err, "execution reverted")

	suite.ctx = suite.ctx.WithBlockTime(ubd.Entries[0].CompletionTime).WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
	beforeERC20 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(erc20types.ModuleAddress.Bytes()), denom)
	_, err = suite.app.EndBlocker(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(coin, suite.app.BankKeeper.GetBalance(suite.ctx, del, denom))
	suite.Require().Equal(beforeERC20, suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(erc20types.ModuleAddress.Bytes()), denom))
	_, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, types.MultiStakingUnlockID(del.String(), val.OperatorAddress))
	suite.Require().False(found)
	// The failed conversion cannot strand the entry or pay it out a second time.
	suite.ctx = suite.ctx.WithBlockHeight(3).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
	_, err = suite.app.EndBlocker(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(coin, suite.app.BankKeeper.GetBalance(suite.ctx, del, denom))
}
