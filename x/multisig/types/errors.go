package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/multisig module sentinel errors
var (
	ErrMissingMembers                  = sdkerrors.Register(ModuleName, 1, "members must be specified")               //nolint:staticcheck
	ErrDuplicateMember                 = sdkerrors.Register(ModuleName, 2, "duplicate member address found")          //nolint:staticcheck
	ErrZeroMemberWeight                = sdkerrors.Register(ModuleName, 3, "member weight must be greater than zero") //nolint:staticcheck
	ErrZeroThreshold                   = sdkerrors.Register(ModuleName, 4, "threshold must be greater than 0")        //nolint:staticcheck
	ErrTotalWeightGreaterThanThreshold = sdkerrors.Register(ModuleName, 5, "threshold must be less than or equal to the total weight")
	ErrWeightsOverflow                 = sdkerrors.Register(ModuleName, 6, "sum of heights overflow")
	ErrNotAMember                      = sdkerrors.Register(ModuleName, 7, "not a member of account")
	ErrInvalidSigner                   = sdkerrors.Register(ModuleName, 8, "expected multisig account as only signer for proposal message") //nolint:staticcheck
	ErrUnroutableProposalMsg           = sdkerrors.Register(ModuleName, 9, "proposal message not recognized by router")                     //nolint:staticcheck
)
