package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/realio-tech/multi-staking-module/testing/simapp"
	"github.com/realio-tech/multi-staking-module/testutil"

	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	MultiStakingDenomA = "ario"
	MultiStakingDenomB = "arst"
)

func (suite *KeeperTestSuite) TestCreateValidator() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error)
		expOut   sdk.Coin
		expErr   bool
	}{
		{
			name: "3001 token, weight 0.3, expect 900",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, sdk.MustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(3001))
				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(900)),
			expErr: false,
		},
		{
			name: "25 token, weight 0.5, expect 12",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomB, sdk.MustNewDecFromStr("0.5"))
				bondAmount := sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(25))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: false,
		},
		{
			name: "invalid bond token",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(25))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}
				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: true,
		},
		{
			name: "invalid validator address",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, sdk.MustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(3001))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker: "NewValidator",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  sdk.AccAddress([]byte("invalid")).String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			valCoins := sdk.NewCoins(sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000)), sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(10000)))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, valCoins)

			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper, multistakingkeeper.NewMsgServerImpl(*suite.msKeeper))
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				actualBond, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().True(found)
				suite.Require().Equal(bondAmount.Amount, lockRecord.LockedCoin.Amount)
				suite.Require().Equal(tc.expOut.Amount, actualBond.Shares.RoundInt())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEditValidator() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error)
		expErr   bool
	}{
		{
			name: "success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(testutil.GenValAddress(), stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "negative rate",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("-0.01")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "less than minimum rate",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.01")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "more than max rate",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.11")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "min self delegation more than validator tokens",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(10000)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "min self delegation more than old min delegation value",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(100)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.OneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, sdk.NewCoins(bondAmount))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			msgServer.CreateValidator(suite.ctx, &createMsg)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
			originMsg, err := tc.malleate(suite.ctx, msgServer)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				validatorInfo, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(originMsg.ValidatorAddress))
				if found {
					suite.Require().Equal(validatorInfo.Description, originMsg.Description)
					suite.Require().Equal(validatorInfo.MinSelfDelegation, &originMsg.MinSelfDelegation)
					suite.Require().Equal(validatorInfo.Commission.CommissionRates.Rate, &originMsg.CommissionRate)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDelegate() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgDelegate, error)
		expRate  sdk.Dec
		expErr   bool
	}{
		{
			name: "success and not change rate",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgDelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr, valAddr, bondAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return *delMsg, err
			},
			expRate: sdk.OneDec(),
			expErr:  false,
		},
		{
			name: "rate change from 1 to 0.75 (1000 / 1 + 2000 / 0.5 = 3000 / 0.75)",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgDelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr, valAddr, bondAmount)
				_, err := msgServer.Delegate(ctx, delMsg)

				if err != nil {
					return *delMsg, err
				}
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, sdk.MustNewDecFromStr("0.5"))
				bondAmount1 := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(2000))
				delMsg1 := stakingtypes.NewMsgDelegate(delAddr, valAddr, bondAmount1)
				_, err = msgServer.Delegate(ctx, delMsg1)
				return *delMsg, err
			},
			expRate: sdk.MustNewDecFromStr("0.75"),
			expErr:  false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgDelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))

				delMsg := stakingtypes.NewMsgDelegate(delAddr, testutil.GenValAddress(), bondAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return *delMsg, err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgDelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(1000))

				delMsg := stakingtypes.NewMsgDelegate(delAddr, valAddr, bondAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return *delMsg, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.OneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, sdk.NewCoins(userBalance))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			msgServer.CreateValidator(suite.ctx, &createMsg)

			_, err := tc.malleate(suite.ctx, msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate, lockRecord.GetBondWeight())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBeginRedelegate() {
	delAddr := testutil.GenAddress()
	valPubKey1 := testutil.GenPubKey()
	valPubKey2 := testutil.GenPubKey()

	valAddr1 := sdk.ValAddress(valPubKey1.Address())
	valAddr2 := sdk.ValAddress(valPubKey2.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error)
		expRate  []sdk.Dec
		expLock  []math.Int
		expErr   bool
	}{
		{
			name: "redelegate from val1 to val2",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr, valAddr1, valAddr2, bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return *multiStakingMsg, err
			},
			expRate: []sdk.Dec{sdk.OneDec(), sdk.OneDec()},
			expLock: []math.Int{sdk.NewInt(500), sdk.NewInt(1500)},
			expErr:  false,
		},
		{
			name: "delegate 500 more to val1 then change rate and redelegate to val2",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, sdk.MustNewDecFromStr("0.25"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))

				delMsg := stakingtypes.NewMsgDelegate(delAddr, valAddr1, bondAmount)
				_, err := msgServer.Delegate(ctx, delMsg)

				bondAmount1 := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
				redelMsg := stakingtypes.NewMsgBeginRedelegate(delAddr, valAddr1, valAddr2, bondAmount1)
				if err != nil {
					return *redelMsg, err
				}

				_, err = msgServer.BeginRedelegate(ctx, redelMsg)
				return *redelMsg, err
			},
			expRate: []sdk.Dec{sdk.MustNewDecFromStr("0.75"), sdk.MustNewDecFromStr("0.875")},
			expLock: []math.Int{sdk.NewInt(500), sdk.NewInt(2000)},
			expErr:  false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr, valAddr1, testutil.GenValAddress(), bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return *multiStakingMsg, err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(1000))

				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr, valAddr1, valAddr2, bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return *multiStakingMsg, err
			},
			expErr: true,
		},
		{
			name: "setup val3 with bond denom is arst then redelgate from val1 to val3",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgBeginRedelegate, error) {
				valPubKey3 := testutil.GenPubKey()
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))

				valAddr3 := sdk.ValAddress(valPubKey3.Address())
				createMsg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.1"),
					},
					MinSelfDelegation: sdk.NewInt(200),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr3.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey3),
					Value:             sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(1000)),
				}
				msgServer.CreateValidator(suite.ctx, &createMsg)
				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr, valAddr1, valAddr3, bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return *multiStakingMsg, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.OneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr1.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey1),
				Value:             bondAmount,
			}
			createMsg2 := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.1"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr2.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey2),
				Value:             bondAmount,
			}
			msgServer.CreateValidator(suite.ctx, &createMsg)
			msgServer.CreateValidator(suite.ctx, &createMsg2)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})

			_, err := tc.malleate(suite.ctx, msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId1 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr1.String())
				lockRecord1, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId1)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate[0], lockRecord1.GetBondWeight())
				suite.Require().Equal(tc.expLock[0], lockRecord1.LockedCoin.Amount)

				lockId2 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr2.String())
				lockRecord2, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId2)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate[1], lockRecord2.GetBondWeight())
				suite.Require().Equal(tc.expLock[1], lockRecord2.LockedCoin.Amount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUndelegate() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name      string
		malleate  func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (*stakingtypes.MsgUndelegate, error)
		expUnlock math.Int
		expLock   math.Int
		expErr    bool
	}{
		{
			name: "undelegate success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (*stakingtypes.MsgUndelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr, valAddr, bondAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return multiStakingMsg, err
			},
			expUnlock: sdk.NewInt(500),
			expLock:   sdk.NewInt(500),
			expErr:    false,
		},
		{
			name: "undelegate 250 then undelegate 500",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (*stakingtypes.MsgUndelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(250))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr, valAddr, bondAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				if err != nil {
					return multiStakingMsg, err
				}
				bondAmount1 := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg1 := stakingtypes.NewMsgUndelegate(delAddr, valAddr, bondAmount1)
				_, err = msgServer.Undelegate(ctx, multiStakingMsg1)
				return multiStakingMsg1, err
			},
			expUnlock: sdk.NewInt(750),
			expLock:   sdk.NewInt(250),
			expErr:    false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (*stakingtypes.MsgUndelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr, testutil.GenValAddress(), bondAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return multiStakingMsg, err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (*stakingtypes.MsgUndelegate, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(1000))

				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr, testutil.GenValAddress(), bondAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return multiStakingMsg, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.OneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}

			msgServer.CreateValidator(suite.ctx, &createMsg)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})

			_, err := tc.malleate(suite.ctx, msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId1 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord1, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId1)
				suite.Require().True(found)
				suite.Require().Equal(tc.expLock, lockRecord1.LockedCoin.Amount)

				unlockID := multistakingtypes.MultiStakingUnlockID(delAddr.String(), valAddr.String())
				unbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
				suite.Require().True(found)
				suite.Require().Equal(tc.expUnlock, unbondRecord.Entries[0].UnlockingCoin.Amount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCancelUnbondingDelegation() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name      string
		malleate  func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error)
		expUnlock math.Int
		expLock   math.Int
		expErr    bool
	}{
		{
			name: "undelegate success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, valAddr, ctx.BlockHeight(), bondAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				return *multiStakingMsg, err
			},
			expUnlock: sdk.NewInt(500),
			expLock:   sdk.NewInt(1000),
			expErr:    false,
		},
		// {
		// 	name: "undelegate 250 then undelegate 500",
		// 	malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error) {
		// 		bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(250))
		// 		multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, valAddr, ctx.BlockHeight(), bondAmount)
		// 		_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
		// 		if err != nil {
		// 			return *multiStakingMsg, err
		// 		}
		// 		bondAmount1 := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
		// 		multiStakingMsg1 := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, valAddr, ctx.BlockHeight(), bondAmount1)
		// 		_, err = msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg1)
		// 		return *multiStakingMsg1, err
		// 	},
		// 	expUnlock: sdk.NewInt(250),
		// 	expLock:   sdk.NewInt(1750),
		// 	expErr:    false,
		// },
		// {
		// 	name: "not found validator",
		// 	malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error) {
		// 		bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
		// 		multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, testutil.GenValAddress(), ctx.BlockHeight(), bondAmount)
		// 		_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
		// 		return *multiStakingMsg, err
		// 	},
		// 	expErr: true,
		// },
		// {
		// 	name: "not allow token",
		// 	malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error) {
		// 		bondAmount := sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(1000))

		// 		multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, valAddr, ctx.BlockHeight(), bondAmount)
		// 		_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
		// 		return *multiStakingMsg, err
		// 	},
		// 	expErr: true,
		// },
		// {
		// 	name: "not found entry at height 20",
		// 	malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (stakingtypes.MsgCancelUnbondingDelegation, error) {
		// 		bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(500))
		// 		multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr, valAddr, 20, bondAmount)
		// 		_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
		// 		return *multiStakingMsg, err
		// 	},
		// 	expErr: true,
		// },
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.OneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(2000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000))
			simapp.FundAccount(suite.app, suite.ctx, delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}

			msgServer.CreateValidator(suite.ctx, &createMsg)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})

			unbondMsg := stakingtypes.NewMsgUndelegate(delAddr, valAddr, sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(1000)))
			msgServer.Undelegate(suite.ctx, unbondMsg)

			_, err := tc.malleate(suite.ctx, msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId1 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord1, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId1)
				suite.Require().True(found)
				suite.Require().Equal(tc.expLock, lockRecord1.LockedCoin.Amount)

				unlockID := multistakingtypes.MultiStakingUnlockID(delAddr.String(), valAddr.String())
				unbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
				suite.Require().True(found)
				suite.Require().Equal(tc.expUnlock, unbondRecord.Entries[0].UnlockingCoin.Amount)
			}
		})
	}
}
