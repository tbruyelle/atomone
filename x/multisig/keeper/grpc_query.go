package keeper

import (
	"context"
	"errors"

	"github.com/atomone-hub/atomone/collections"
	"github.com/atomone-hub/atomone/x/multisig/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) Multisig(goCtx context.Context, req *types.QueryMultisigRequest) (*types.QueryMultisigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addrBz, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address %s: %v", req.Address, err)
	}
	m, err := k.multisigs.Get(goCtx, addrBz)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "multisig %s doesn't exist", req.Address)
		}
		return nil, err
	}
	return &types.QueryMultisigResponse{Multisig: &m}, nil
}
