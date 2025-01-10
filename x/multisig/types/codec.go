package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	govtypes "github.com/atomone-hub/atomone/x/gov/types/v1"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAccount{}, &MsgCreateProposal{}, &MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateAccount{}, "atomone/multisig/v1/MsgCreateAccount")
	legacy.RegisterAminoMsg(cdc, &MsgCreateProposal{}, "atomone/multisig/v1/MsgCreateProposal")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "atomone/x/multisig/v1/MsgUpdateParams")
	cdc.RegisterConcrete(&Params{}, "atomone/multisig/v1/Params", nil)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Need to add registration in the atomone multisig amino for all modules
	// because of the MsgCreateProposal which can embed any other messages.
	banktypes.RegisterLegacyAminoCodec(amino)
	govtypes.RegisterLegacyAminoCodec(amino)
	consensustypes.RegisterLegacyAminoCodec(amino)
	crisistypes.RegisterLegacyAminoCodec(amino)
	distributiontypes.RegisterLegacyAminoCodec(amino)
	evidencetypes.RegisterLegacyAminoCodec(amino)
	minttypes.RegisterLegacyAminoCodec(amino)
	slashingtypes.RegisterLegacyAminoCodec(amino)
	stakingtypes.RegisterLegacyAminoCodec(amino)
	upgradetypes.RegisterLegacyAminoCodec(amino)
}
