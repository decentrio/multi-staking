package keeper

import (
	"context"
	"errors"
	"testing"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type failingMintBank struct {
	types.BankKeeper
	err  error
	sent bool
}

func (b *failingMintBank) MintCoins(context.Context, string, sdk.Coins) error { return b.err }
func (b *failingMintBank) SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error {
	b.sent = true
	return nil
}
func TestMintCoinPropagatesErrorWithoutTransfer(t *testing.T) {
	mintErr := errors.New("mint failed")
	bank := &failingMintBank{err: mintErr}
	k := Keeper{bankKeeper: bank}
	require.ErrorIs(t, k.MintCoin(context.Background(), sdk.AccAddress(make([]byte, 20)), sdk.NewInt64Coin("stake", 10)), mintErr)
	require.False(t, bank.sent)
}

// Model bank balances and EVM writes in an SDK store so the test exercises
// CacheContext rollback, including writes made before conversion fails.
type payoutBank struct {
	types.BankKeeper
	key *storetypes.KVStoreKey
	err error
}

func (b payoutBank) SendCoinsFromModuleToAccount(ctx context.Context, _ string, _ sdk.AccAddress, coins sdk.Coins) error {
	if b.err != nil {
		return b.err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.KVStore(b.key).Set([]byte("payout"), []byte(coins.String()))
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent("payout"))
	return nil
}

type conversionKeeper struct {
	types.ERC20Keeper
	key    *storetypes.KVStoreKey
	pair   bool
	err    error
	called bool
}

func (e *conversionKeeper) GetTokenPairID(sdk.Context, string) []byte {
	if e.pair {
		return []byte{1}
	}
	return nil
}
func (e *conversionKeeper) ConvertCoin(ctx context.Context, _ *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error) {
	e.called = true
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.KVStore(e.key).Delete([]byte("payout"))
	sdkCtx.KVStore(e.key).Set([]byte("evm"), []byte("transfer"))
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent("conversion"))
	if e.err != nil {
		// Match Cosmos EVM's reset-and-consume-full-limit behavior on revert.
		sdkCtx.GasMeter().RefundGas(sdkCtx.GasMeter().GasConsumed(), "reset")
		sdkCtx.GasMeter().ConsumeGas(sdkCtx.GasMeter().Limit(), "EVM revert")
	}
	return nil, e.err
}
func TestUnescrowConversionIsolation(t *testing.T) {
	for _, tc := range []struct {
		name                             string
		pair, conversionFails, bankFails bool
	}{
		{name: "native payout"},
		{name: "successful conversion", pair: true},
		{name: "reverting conversion retains cosmos payout", pair: true, conversionFails: true},
		{name: "bank failure propagates", pair: true, bankFails: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key := storetypes.NewKVStoreKey("test")
			ctx := testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient"))
			bank := payoutBank{key: key}
			converter := &conversionKeeper{key: key, pair: tc.pair}
			if tc.conversionFails {
				converter.err = errors.New("execution reverted")
			}
			if tc.bankFails {
				bank.err = errors.New("insufficient escrow")
			}
			k := Keeper{bankKeeper: bank, erc20keeper: converter}
			err := k.UnescrowCoinTo(ctx, sdk.AccAddress(make([]byte, 20)), sdk.NewInt64Coin("token", 10))
			if tc.bankFails {
				require.ErrorIs(t, err, bank.err)
				require.False(t, converter.called)
				require.Empty(t, ctx.EventManager().Events())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.pair, converter.called)
			converted := tc.pair && !tc.conversionFails
			require.Equal(t, converted, ctx.KVStore(key).Has([]byte("evm")))
			if converted {
				require.False(t, ctx.KVStore(key).Has([]byte("payout")))
				require.Len(t, ctx.EventManager().Events(), 2)
			} else {
				require.Equal(t, []byte("10token"), ctx.KVStore(key).Get([]byte("payout")))
				require.Len(t, ctx.EventManager().Events(), 1)
			}
		})
	}
}
