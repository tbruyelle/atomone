package types

import (
	"cosmossdk.io/errors"
)

// x/multisig module sentinel errors
var (
	ErrMissingMembers                  = errors.Register(ModuleName, 1, "members must be specified")
	ErrDuplicateMember                 = errors.Register(ModuleName, 2, "duplicate member address found")
	ErrZeroMemberWeight                = errors.Register(ModuleName, 3, "member weight must be greater than zero")
	ErrZeroThreshold                   = errors.Register(ModuleName, 4, "threshold must be greater than 0")
	ErrTotalWeightGreaterThanThreshold = errors.Register(ModuleName, 5, "threshold must be less than or equal to the total weight")
	ErrWeightsOverflow                 = errors.Register(ModuleName, 6, "sum of heights overflow")
	ErrNotAMember                      = errors.Register(ModuleName, 7, "not a member of account")
	ErrInvalidSigner                   = errors.Register(ModuleName, 8, "expected multisig account as only signer for proposal message")
	ErrUnroutableProposalMsg           = errors.Register(ModuleName, 9, "proposal message not recognized by router")
	ErrInvalidVote                     = errors.Register(ModuleName, 10, "invalid vote option")
	ErrExecuteWoThreshold              = errors.Register(ModuleName, 11, "cannot execute proposal with unmet threshold")
)
