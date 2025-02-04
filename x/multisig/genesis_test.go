package multisig_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/atomone-hub/atomone/x/multisig"
	"github.com/atomone-hub/atomone/x/multisig/testutil"
	"github.com/atomone-hub/atomone/x/multisig/types"
)

func TestGenesis(t *testing.T) {
	var (
		testAddrs    = simtestutil.CreateRandomAccounts(3)
		accAddr      = testAddrs[0].String()
		member1Addr  = testAddrs[1].String()
		member2Addr  = testAddrs[2].String()
		genesisState = types.GenesisState{
			Params: types.DefaultParams(),
			Accounts: []*types.Account{
				{
					Address: accAddr,
					Creator: member1Addr,
					Members: []types.Member{
						{
							Address: member1Addr,
							Weight:  1,
						},
						{
							Address: member2Addr,
							Weight:  2,
						},
					},
				},
			},
			Proposals: []*types.Proposal{
				{
					Id:             42,
					AccountAddress: accAddr,
					Proposer:       member1Addr,
				},
			},
			Votes: []*types.Vote{
				{
					VoterAddress:   member1Addr,
					AccountAddress: accAddr,
					ProposalId:     42,
					Vote:           types.VoteOption_VOTE_OPTION_YES,
				},
				{
					VoterAddress:   member2Addr,
					AccountAddress: accAddr,
					ProposalId:     42,
					Vote:           types.VoteOption_VOTE_OPTION_NO,
				},
			},
		}
		msgs = []sdk.Msg{
			banktypes.NewMsgSend(testAddrs[0], testAddrs[1], sdk.NewCoins(sdk.NewInt64Coin("uatone", 1))),
		}
	)
	anys, err := sdktx.SetMsgs(msgs)
	require.NoError(t, err)
	genesisState.Proposals[0].Messages = anys

	k, _, ctx := testutil.SetupMultisigKeeper(t)

	multisig.InitGenesis(ctx, *k, genesisState)
	got := multisig.ExportGenesis(ctx, *k)

	require.NotNil(t, got)
	require.Equal(t, genesisState, *got)
}
