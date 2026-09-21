package keeper_test

import (
	"time"

	"github.com/realio-tech/multi-staking-module/test"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Regression: cancelling (part of) an unbonding delegation must keep the unlock
// entry and the staking unbonding entry in sync, whatever the bond weight and
// whether or not the validator was slashed while unbonding. Before, a weight
// below one made the unlock entry shrink by the requested amount while the
// staking entry only shrank by the truncated bond value, so at maturity the
// unlock amount derived from the remaining bond exceeded the unlock entry and
// the EndBlocker panicked.
func (suite *KeeperTestSuite) TestCancelUnbondingKeepsEntriesInSync() {
	testCases := []struct {
		name         string
		weight       string
		stake        int64
		cancel       int64
		slash        string
		expUnlock    int64 // unlock entry amount right after the cancel (0 = entry removed)
		expLock      int64
		expMatureBal int64 // staker balance once the unbonding period has passed
	}{
		{name: "weight 0.4 partial cancel rounds down", weight: "0.4", stake: 10, cancel: 7, expUnlock: 5, expLock: 5, expMatureBal: 5},
		{name: "weight 0.4 cancel all leaves no dust entry", weight: "0.4", stake: 11, cancel: 11, expUnlock: 0, expLock: 11, expMatureBal: 0},
		{name: "weight 0.4 cancel all but rounding dust", weight: "0.4", stake: 11, cancel: 10, expUnlock: 0, expLock: 11, expMatureBal: 0},
		{name: "weight 1 partial cancel", weight: "1", stake: 100, cancel: 40, expUnlock: 60, expLock: 40, expMatureBal: 60},
		{name: "weight 1 cancel all", weight: "1", stake: 100, cancel: 100, expUnlock: 0, expLock: 100, expMatureBal: 0},
		{name: "weight 1 slashed then cancel all", weight: "1", stake: 100, cancel: 100, slash: "0.5", expUnlock: 0, expLock: 100, expMatureBal: 0},
		{name: "weight 1 slashed then cancel more than what is left", weight: "1", stake: 100, cancel: 60, slash: "0.5", expUnlock: 0, expLock: 100, expMatureBal: 0},
		{name: "weight 1 slashed then cancel part", weight: "1", stake: 100, cancel: 20, slash: "0.5", expUnlock: 80, expLock: 20, expMatureBal: 30},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			denom := MultiStakingDenomA
			suite.msKeeper.SetBondWeight(suite.ctx, denom, math.LegacyMustNewDecFromStr(tc.weight))

			// fresh, bonded validator so the share exchange rate is 1:1
			valPubKey := test.GenPubKey()
			valAddr := sdk.ValAddress(valPubKey.Address())
			valAcc := sdk.AccAddress(valPubKey.Address())
			suite.FundAccount(valAcc, sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(20_000_000))))
			msgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper)
			_, err := msgServer.CreateValidator(suite.ctx, &stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{Moniker: "v"},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
				},
				MinSelfDelegation: math.NewInt(100),
				DelegatorAddress:  valAcc.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             sdk.NewCoin(denom, math.NewInt(10_000_000)),
			})
			suite.Require().NoError(err)
			_, err = suite.app.EndBlocker(suite.ctx.WithBlockGasMeter(storetypes.NewInfiniteGasMeter()))
			suite.Require().NoError(err)

			coin := sdk.NewCoin(denom, math.NewInt(tc.stake))
			staker := suite.CreateAndFundAccount(sdk.NewCoins(coin))
			_, err = msgServer.Delegate(suite.ctx, &stakingtypes.MsgDelegate{
				DelegatorAddress: staker.String(), ValidatorAddress: valAddr.String(), Amount: coin,
			})
			suite.Require().NoError(err)

			// height 2: undelegate everything
			suite.ctx = suite.ctx.WithBlockHeight(2).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
			_, err = msgServer.Undelegate(suite.ctx, &stakingtypes.MsgUndelegate{
				DelegatorAddress: staker.String(), ValidatorAddress: valAddr.String(), Amount: coin,
			})
			suite.Require().NoError(err)

			// height 3: the validator is slashed for an infraction the unbonding delegation was exposed to
			suite.ctx = suite.ctx.WithBlockHeight(3).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
			if tc.slash != "" {
				val, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
				suite.Require().NoError(err)
				consAddr, err := val.GetConsAddr()
				suite.Require().NoError(err)
				power := suite.app.StakingKeeper.TokensToConsensusPower(suite.ctx, val.Tokens)
				suite.Require().NoError(suite.app.SlashingKeeper.Slash(suite.ctx, consAddr, math.LegacyMustNewDecFromStr(tc.slash), power, 2))
			}

			_, err = msgServer.CancelUnbondingDelegation(suite.ctx, stakingtypes.NewMsgCancelUnbondingDelegation(
				staker.String(), valAddr.String(), 2, sdk.NewCoin(denom, math.NewInt(tc.cancel)),
			))
			suite.Require().NoError(err)

			unlockID := multistakingtypes.MultiStakingUnlockID(staker.String(), valAddr.String())
			unlock, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
			if tc.expUnlock == 0 {
				suite.Require().False(found, "unlock entry would never mature")
			} else {
				suite.Require().True(found)
				suite.Require().Equal(math.NewInt(tc.expUnlock), unlock.Entries[0].UnlockingCoin.Amount)
			}
			lock, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, multistakingtypes.MultiStakingLockID(staker.String(), valAddr.String()))
			suite.Require().True(found)
			suite.Require().Equal(math.NewInt(tc.expLock), lock.LockedCoin.Amount)

			// pass the unbonding period
			suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(365 * 24 * time.Hour)).WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
			suite.Require().NotPanics(func() {
				_, err = suite.app.EndBlocker(suite.ctx)
			})
			suite.Require().NoError(err)

			_, found = suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
			suite.Require().False(found, "unlock entry must be settled or removed")
			suite.Require().Equal(math.NewInt(tc.expMatureBal), suite.app.BankKeeper.GetBalance(suite.ctx, staker, denom).Amount)
		})
	}
}

// A cancel amount larger than the unlock entry is still rejected.
func (suite *KeeperTestSuite) TestCancelUnbondingMoreThanUnlockEntryFails() {
	vals, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	suite.Require().NoError(err)
	val := vals[0]
	operatorAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	suite.Require().NoError(err)
	denom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, operatorAddr)

	coin := sdk.NewCoin(denom, math.NewInt(1000))
	staker := suite.CreateAndFundAccount(sdk.NewCoins(coin))
	msgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper)
	_, err = msgServer.Delegate(suite.ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: staker.String(), ValidatorAddress: val.OperatorAddress, Amount: coin,
	})
	suite.Require().NoError(err)

	suite.ctx = suite.ctx.WithBlockHeight(2).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
	_, err = msgServer.Undelegate(suite.ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: staker.String(), ValidatorAddress: val.OperatorAddress, Amount: sdk.NewCoin(denom, math.NewInt(500)),
	})
	suite.Require().NoError(err)

	_, err = msgServer.CancelUnbondingDelegation(suite.ctx, stakingtypes.NewMsgCancelUnbondingDelegation(
		staker.String(), val.OperatorAddress, 2, sdk.NewCoin(denom, math.NewInt(501)),
	))
	suite.Require().Error(err)
}
