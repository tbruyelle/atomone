package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
)

var _, _, _ sdk.Msg = &MsgUpdateParams{}, &MsgCreateAccount{}, &MsgCreateProposal{}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return m.Params.ValidateBasic()
}

// GetSignBytes returns the message bytes to sign over.
func (m MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{authority}
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgCreateAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Members) == 0 {
		return ErrMissingMembers
	}
	if m.Threshold <= 0 {
		return ErrZeroThreshold
	}
	membersMap := map[string]struct{}{} // to check for duplicates
	for i := range m.Members {
		if _, ok := membersMap[m.Members[i].Address]; ok {
			return ErrDuplicateMember
		}

		membersMap[m.Members[i].Address] = struct{}{}

		if m.Members[i].Weight == 0 {
			return ErrZeroMemberWeight
		}
		_, err := sdk.AccAddressFromBech32(m.Members[i].Address)
		if err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", err)
		}
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (m MsgCreateAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (m MsgCreateAccount) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{authority}
}

// GetMsgs unpacks m.Messages Any's into sdk.Msg's
func (m *MsgCreateProposal) GetMsgs() ([]sdk.Msg, error) {
	return sdktx.GetMsgs(m.Proposal.Messages, "sdk.MsgProposal")
}

// SetMsgs packs sdk.Msg's into m.Messages Any's
// NOTE: this will overwrite any existing messages
func (m *MsgCreateProposal) SetMsgs(msgs []sdk.Msg) error {
	anys, err := sdktx.SetMsgs(msgs)
	if err != nil {
		return err
	}

	m.Proposal.Messages = anys
	return nil
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgCreateProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid multisig account address: %s", err)
	}
	if m.Proposal == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("Proposal must be non-nil") //nolint:staticcheck
	}
	// TODO assert max length
	if m.Proposal.Title == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("proposal title cannot be empty") //nolint:staticcheck
	}
	// TODO assert max length
	if m.Proposal.Summary == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("proposal summary cannot be empty") //nolint:staticcheck
	}
	if len(m.Proposal.Messages) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("Proposal.Messages length must be non-nil") //nolint:staticcheck
	}
	msgs, err := m.GetMsgs()
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("unable to read proposal messages: %v", err) //nolint:staticcheck
	}
	for i, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("validation fail for proposal message %d: %v", i, err) //nolint:staticcheck
		}
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (m MsgCreateProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (m MsgCreateProposal) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{authority}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// NOTE: required because MsgCreateProposal embeds types.Anys.
func (m MsgCreateProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return sdktx.UnpackInterfaces(unpacker, m.Proposal.Messages)
}
