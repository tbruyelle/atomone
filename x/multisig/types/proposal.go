package types

import (
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
)

var _ codectypes.UnpackInterfacesMessage = &Proposal{}

// NewProposal creates a new Proposal instance
func NewProposal(id uint64, accountAddr string, submitTime time.Time, proposer, title, summary string, messages []sdk.Msg) (Proposal, error) {
	msgs, err := sdktx.SetMsgs(messages)
	if err != nil {
		return Proposal{}, err
	}
	p := Proposal{
		Id:             id,
		AccountAddress: accountAddr,
		Messages:       msgs,
		SubmitTime:     &submitTime,
		Proposer:       proposer,
		Title:          title,
		Summary:        summary,
	}
	return p, nil
}

// GetMessages returns the proposal messages
func (p Proposal) GetMsgs() ([]sdk.Msg, error) {
	return sdktx.GetMsgs(p.Messages, "sdk.MsgProposal")
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Proposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return sdktx.UnpackInterfaces(unpacker, p.Messages)
}
