package types

import (
	"context"
	time "time"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

type StakingKeeper interface {
	DequeueAllMatureUBDQueue(ctx context.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair)
	BondDenom(ctx context.Context) (res string)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetUnbondingDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (unlock stakingtypes.UnbondingDelegation, found bool)
	InitGenesis(ctx context.Context, data *stakingtypes.GenesisState) (res []abci.ValidatorUpdate)
	ExportGenesis(ctx context.Context) *stakingtypes.GenesisState
	GetParams(ctx context.Context) stakingtypes.Params
}

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// ERC20Keeper defines the conversion and registration operations used by this module.
type ERC20Keeper interface {
	GetTokenPairID(ctx sdk.Context, token string) []byte
	GetTokenDenom(ctx sdk.Context, tokenAddress common.Address) (string, error)
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
	ConvertERC20(ctx context.Context, msg *erc20types.MsgConvertERC20) (*erc20types.MsgConvertERC20Response, error)
	RegisterERC20(ctx context.Context, msg *erc20types.MsgRegisterERC20) (*erc20types.MsgRegisterERC20Response, error)
	GetTokenPair(ctx sdk.Context, id []byte) (erc20types.TokenPair, bool)
}
