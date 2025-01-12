package types

import "github.com/atomone-hub/atomone/collections"

const (
	// ModuleName defines the module name
	ModuleName = "multisig"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	KeyParams         = collections.NewPrefix(0)
	KeyAccounts       = collections.NewPrefix(1)
	KeyAccountNumber  = collections.NewPrefix(2)
	KeyProposals      = collections.NewPrefix(3)
	KeyProposalNumber = collections.NewPrefix(4)
	KeyVotes          = collections.NewPrefix(5)
)
