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
	ErrWrongMemberAddress              = sdkerrors.Register(ModuleName, 4, "wrong member address")                    //nolint:staticcheck
	ErrZeroThreshold                   = sdkerrors.Register(ModuleName, 5, "threshold must be greater than 0")        //nolint:staticcheck
	ErrTotalWeightGreaterThanThreshold = sdkerrors.Register(ModuleName, 6, "threshold must be less than or equal to the total weight")
	ErrWeightsOverflow                 = sdkerrors.Register(ModuleName, 7, "sum of heights overflow")
)
