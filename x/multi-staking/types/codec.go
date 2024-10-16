package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateValidator{}, "multi-staking/MsgCreateValidator")
	legacy.RegisterAminoMsg(cdc, &MsgEditValidator{}, "multi-staking/MsgEditValidator")
	legacy.RegisterAminoMsg(cdc, &MsgDelegate{}, "multi-staking/MsgDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgUndelegate{}, "multi-staking/MsgUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgBeginRedelegate{}, "multi-staking/MsgBeginRedelegate")
	legacy.RegisterAminoMsg(cdc, &MsgCancelUnbonding{}, "multi-staking/MsgCancelUnbondDelegation")
	legacy.RegisterAminoMsg(cdc, &MsgVoteWeighted{}, "multi-staking/MsgVoteWeighted")
	legacy.RegisterAminoMsg(cdc, &MsgVote{}, "multi-staking/MsgVote")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "multi-staking/MsgSetWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawReward{}, "multi-staking/MsgWithdrawReward")
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidator{},
		&MsgEditValidator{},
		&MsgDelegate{},
		&MsgUndelegate{},
		&MsgBeginRedelegate{},
		&MsgCancelUnbonding{},
		&MsgVoteWeighted{},
		&MsgVote{},
		&MsgSetWithdrawAddress{},
		&MsgWithdrawReward{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}
