package types_test

import (
	"testing"

	"github.com/atomone-hub/atomone/x/multisig/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		name          string
		genState      *types.GenesisState
		expectedError string
	}{
		{
			name:     "default is valid",
			genState: types.DefaultGenesis(),
		},
		{
			name:     "valid genesis state",
			genState: &types.GenesisState{},
		},
		{
			name: "account creator is not a member",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "creator1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
			},
			expectedError: "account acc1 creator not a member",
		},
		{
			name: "duplicate account members",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{
							{Address: "member1"},
							{Address: "member1"},
						},
					},
				},
			},
			expectedError: "duplicate member member1 of account acc1",
		},
		{
			name: "duplicate accounts",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
			},
			expectedError: "duplicate account address: acc1",
		},
		{
			name: "proposal with unknown account address",
			genState: &types.GenesisState{
				Proposals: []*types.Proposal{
					{
						Id:             42,
						AccountAddress: "acc1",
					},
				},
			},
			expectedError: "proposal id 42: account acc1 does not exists",
		},
		{
			name: "proposal with proposer not a member of account",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
				Proposals: []*types.Proposal{
					{
						Id:             42,
						AccountAddress: "acc1",
						Proposer:       "prop1",
					},
				},
			},
			expectedError: "proposal id 42 proposer prop1: not a member of account acc1",
		},
		{
			name: "duplicate proposals",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
				Proposals: []*types.Proposal{
					{
						Id:             42,
						AccountAddress: "acc1",
						Proposer:       "member1",
					},
					{
						Id:             42,
						AccountAddress: "acc1",
						Proposer:       "member1",
					},
				},
			},
			expectedError: "duplicate proposal id: 42",
		},
		{
			name: "vote with unknown account address",
			genState: &types.GenesisState{
				Votes: []*types.Vote{
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "voter1",
					},
				},
			},
			expectedError: "vote {proposalId:42 voter:voter1}: account acc1 does not exists",
		},
		{
			name: "vote with voter not a member of account",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
				Votes: []*types.Vote{
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "voter1",
					},
				},
			},
			expectedError: "vote {proposalId:42 voter:voter1}: voter not a member of account acc1",
		},
		{
			name: "vote with unknown proposal id",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
				Votes: []*types.Vote{
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "member1",
					},
				},
			},
			expectedError: "vote {proposalId:42 voter:member1}: proposal does not exists",
		},
		{
			name: "duplicate votes",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{{Address: "member1"}},
					},
				},
				Proposals: []*types.Proposal{
					{
						Id:             42,
						AccountAddress: "acc1",
						Proposer:       "member1",
					},
				},
				Votes: []*types.Vote{
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "member1",
					},
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "member1",
					},
				},
			},
			expectedError: "duplicate vote: {proposalId:42 voter:member1}",
		},
		{
			name: "valid with account, proposal and votes",
			genState: &types.GenesisState{
				Accounts: []*types.Account{
					{
						Address: "acc1",
						Creator: "member1",
						Members: []types.Member{
							{Address: "member1"},
							{Address: "member2"},
						},
					},
				},
				Proposals: []*types.Proposal{
					{
						Id:             42,
						AccountAddress: "acc1",
						Proposer:       "member1",
					},
				},
				Votes: []*types.Vote{
					{
						AccountAddress: "acc1",
						ProposalId:     42,
						VoterAddress:   "member2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genState.Validate()

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}
