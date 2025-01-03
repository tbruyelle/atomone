package keeper

import (
	"context"
	"errors"

	govtypes "github.com/atomone-hub/atomone/x/gov/types"
	"github.com/atomone-hub/atomone/x/multisig/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateMultisig(goCtx context.Context, msg *types.MsgCreateMultisig) (*types.MsgCreateMultisigResponse, error) {
	if len(msg.Members) == 0 {
		return nil, types.ErrMissingMembers
	}
	if msg.Threshold <= 0 {
		return nil, types.ErrZeroThreshold
	}

	// set members
	totalWeight := uint64(0)
	membersMap := map[string]struct{}{} // to check for duplicates
	for i := range msg.Members {
		if _, ok := membersMap[msg.Members[i].Address]; ok {
			return nil, types.ErrDuplicateMember
		}

		membersMap[msg.Members[i].Address] = struct{}{}

		if msg.Members[i].Weight == 0 {
			return nil, types.ErrZeroMemberWeight
		}
		addrBz, err := sdk.AccAddressFromBech32(msg.Members[i].Address)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrWrongMemberAddress, "address: %s", msg.Members[9].Address)
		}
		_ = addrBz
		// TODO check members in x/auth?

		totalWeight, err = safeAdd(totalWeight, msg.Members[i].Weight)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrWeightsOverflow, "%v", err)
		}
	}

	// threshold must be less than or equal to the total weight
	if totalWeight < uint64(msg.Threshold) {
		return nil, types.ErrTotalWeightGreaterThanThreshold
	}

	// TODO create multsig address
	// TODO add to store

	return nil, nil
}

func safeAdd(nums ...uint64) (uint64, error) {
	var sum uint64
	for _, num := range nums {
		if newSum := sum + num; newSum < sum {
			return 0, errors.New("overflow")
		} else {
			sum = newSum
		}
	}
	return sum, nil
}

// UpdateParams implements the MsgServer.UpdateParams method.
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
