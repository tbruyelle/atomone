package types_test

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/atomone-hub/atomone/x/gov/types/v1"
	"github.com/atomone-hub/atomone/x/multisig/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// this tests that Amino JSON MsgSubmitProposal.GetSignBytes() still works with Content as Any using the ModuleCdc
func TestMsgSubmitProposal_GetSignBytes(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
	tests := []struct {
		name              string
		msgs              []sdk.Msg
		title             string
		summary           string
		expectedSignBytes string
	}{
		{
			name:              "MsgVote",
			msgs:              []sdk.Msg{v1.NewMsgVote(addrs[0], 1, v1.OptionYes, "")},
			title:             "gov/MsgVote",
			summary:           "Proposal for a governance vote msg",
			expectedSignBytes: `{"type":"atomone/multisig/v1/MsgCreateProposal","value":{"messages":[{"type":"atomone/v1/MsgVote","value":{"option":1,"proposal_id":"1","voter":"cosmos1w3jhxap3gempvr"}}],"summary":"Proposal for a governance vote msg","title":"gov/MsgVote"}}`,
		},
		{
			name:              "MsgSend",
			msgs:              []sdk.Msg{banktypes.NewMsgSend(addrs[0], addrs[0], sdk.NewCoins())},
			title:             "bank/MsgSend",
			summary:           "Proposal for a bank msg send",
			expectedSignBytes: fmt.Sprintf(`{"type":"atomone/multisig/v1/MsgCreateProposal","value":{"messages":[{"type":"cosmos-sdk/MsgSend","value":{"amount":[],"from_address":"%s","to_address":"%s"}}],"summary":"Proposal for a bank msg send","title":"bank/MsgSend"}}`, addrs[0], addrs[0]),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgCreateProposal{
				Sender:  sdk.AccAddress{}.String(),
				Address: sdk.AccAddress{}.String(),
				Title:   tt.title,
				Summary: tt.summary,
			}
			err := msg.SetMsgs(tt.msgs)
			require.NoError(t, err)

			var bz []byte
			require.NotPanics(t, func() {
				bz = msg.GetSignBytes()
			})

			require.Equal(t, tt.expectedSignBytes, string(bz))
		})
	}
}
