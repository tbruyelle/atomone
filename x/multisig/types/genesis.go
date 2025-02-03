package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// weed out duplicate accounts
	accountAddrs := make(map[string]*Account)
	for _, a := range gs.Accounts {
		memberAddrs := make(map[string]struct{})
		if _, ok := accountAddrs[a.Address]; ok {
			return fmt.Errorf("duplicate account address: %s", a.Address)
		}
		if !a.HasMember(a.Creator) {
			return fmt.Errorf("account %s creator not a member", a.Address)
		}
		for _, m := range a.Members {
			if _, ok := memberAddrs[m.Address]; ok {
				return fmt.Errorf("duplicate member %s of account %s", m.Address, a.Address)
			}
			memberAddrs[m.Address] = struct{}{}
		}
		accountAddrs[a.Address] = a
	}

	// weed out duplicate proposals
	proposalIds := make(map[uint64]struct{})
	for _, p := range gs.Proposals {
		if _, ok := proposalIds[p.Id]; ok {
			return fmt.Errorf("duplicate proposal id: %d", p.Id)
		}
		proposalIds[p.Id] = struct{}{}
		// check if related account exists
		a, ok := accountAddrs[p.AccountAddress]
		if !ok {
			return fmt.Errorf("proposal id %d: account %s does not exists", p.Id, p.AccountAddress)
		}
		if !a.HasMember(p.Proposer) {
			return fmt.Errorf("proposal id %d proposer %s: not a member of account %s", p.Id, p.Proposer, p.AccountAddress)
		}
	}

	// weed out duplicate votes
	type voteKey struct {
		proposalId uint64
		voter      string
	}
	voteIds := make(map[voteKey]struct{})
	for _, v := range gs.Votes {
		id := voteKey{v.ProposalId, v.VoterAddress}
		if _, ok := voteIds[id]; ok {
			return fmt.Errorf("duplicate vote: %+v", id)
		}
		voteIds[id] = struct{}{}
		// check if related account exists
		a, ok := accountAddrs[v.AccountAddress]
		if !ok {
			return fmt.Errorf("vote %+v: account %s does not exists", id, v.AccountAddress)
		}
		if !a.HasMember(v.VoterAddress) {
			return fmt.Errorf("vote %+v: voter not a member of account %s", id, v.AccountAddress)
		}
		// check if related proposal exists
		if _, ok := proposalIds[v.ProposalId]; !ok {
			return fmt.Errorf("vote %+v: proposal does not exists", id)
		}
	}

	return gs.Params.ValidateBasic()
}

var _ codectypes.UnpackInterfacesMessage = GenesisState{}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, p := range data.Proposals {
		err := p.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
