package types_test

import (
	"strings"
	"testing"

	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisValidateValidatorMultiStakingCoins(t *testing.T) {
	valA := test.GenValAddress().String()
	valB := test.GenValAddress().String()

	for _, tc := range []struct {
		name    string
		coins   []types.ValidatorMultiStakingCoin
		wantErr string
	}{
		{name: "default genesis"},
		{
			name: "distinct validators sharing a denom",
			coins: []types.ValidatorMultiStakingCoin{
				{ValAddr: valA, CoinDenom: "ario"},
				{ValAddr: valB, CoinDenom: "ario"},
			},
		},
		{
			name: "identical duplicate",
			coins: []types.ValidatorMultiStakingCoin{
				{ValAddr: valA, CoinDenom: "ario"},
				{ValAddr: valA, CoinDenom: "ario"},
			},
			wantErr: "duplicate validator multi staking coin",
		},
		{
			name: "nonadjacent duplicate with different denom",
			coins: []types.ValidatorMultiStakingCoin{
				{ValAddr: valA, CoinDenom: "ario"},
				{ValAddr: valB, CoinDenom: "ario"},
				{ValAddr: valA, CoinDenom: "arst"},
			},
			wantErr: "duplicate validator multi staking coin",
		},
		{
			name: "duplicate with alternate casing",
			coins: []types.ValidatorMultiStakingCoin{
				{ValAddr: valA, CoinDenom: "ario"},
				{ValAddr: strings.ToUpper(valA), CoinDenom: "ario"},
			},
			wantErr: "duplicate validator multi staking coin",
		},
		{
			name:    "invalid validator address",
			coins:   []types.ValidatorMultiStakingCoin{{ValAddr: "invalid", CoinDenom: "ario"}},
			wantErr: "invalid validator multi staking coin address",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			genesis := types.DefaultGenesis()
			genesis.ValidatorMultiStakingCoins = tc.coins
			err := genesis.Validate()
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
