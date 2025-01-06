package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _, _ sdk.Msg = &MsgUpdateParams{}, &MsgCreateAccount{}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return msg.Params.ValidateBasic()
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(msg.Members) == 0 {
		return ErrMissingMembers
	}
	if msg.Threshold <= 0 {
		return ErrZeroThreshold
	}
	membersMap := map[string]struct{}{} // to check for duplicates
	for i := range msg.Members {
		if _, ok := membersMap[msg.Members[i].Address]; ok {
			return ErrDuplicateMember
		}

		membersMap[msg.Members[i].Address] = struct{}{}

		if msg.Members[i].Weight == 0 {
			return ErrZeroMemberWeight
		}
		_, err := sdk.AccAddressFromBech32(msg.Members[i].Address)
		if err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", err)
		}
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams.
func (msg MsgCreateAccount) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{authority}
}
