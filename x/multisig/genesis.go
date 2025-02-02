package multisig

import (
	"github.com/atomone-hub/atomone/collections"
	"github.com/atomone-hub/atomone/x/multisig/keeper"
	"github.com/atomone-hub/atomone/x/multisig/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.Params.Set(ctx, genState.Params)

	for _, account := range genState.Accounts {
		addrBz := sdk.MustAccAddressFromBech32(account.Address)
		k.Accounts.Set(ctx, addrBz, *account)
	}

	for _, proposal := range genState.Proposals {
		addrBz := sdk.MustAccAddressFromBech32(proposal.AccountAddress)
		k.SetProposal(ctx, addrBz, proposal.Id, *proposal)
	}

	for _, vote := range genState.Votes {
		accountAddrBz := sdk.MustAccAddressFromBech32(vote.AccountAddress)
		voterAddrBz := sdk.MustAccAddressFromBech32(vote.VoterAddress)
		k.SetProposalVote(ctx, accountAddrBz, vote.ProposalId, voterAddrBz, *vote)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	var accounts []*types.Account
	err = k.Accounts.Walk(ctx, nil, func(_ []byte, acc types.Account) (bool, error) {
		accounts = append(accounts, &acc)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	var proposals []*types.Proposal
	err = k.Proposals.Walk(ctx, nil, func(_ collections.Pair[[]byte, uint64], proposal types.Proposal) (bool, error) {
		proposals = append(proposals, &proposal)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	var votes []*types.Vote
	err = k.Votes.Walk(ctx, nil,
		func(_ collections.Triple[[]byte, uint64, []byte], vote types.Vote) (bool, error) {
			votes = append(votes, &vote)
			return false, nil
		})
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:    params,
		Accounts:  accounts,
		Proposals: proposals,
		Votes:     votes,
	}
}
