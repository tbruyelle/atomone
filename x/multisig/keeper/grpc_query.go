package keeper

import (
	"context"

	"github.com/atomone-hub/atomone/collections"
	"github.com/atomone-hub/atomone/x/multisig/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (k queryServer) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := k.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k queryServer) Account(ctx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addrBz, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address %s: %v", req.Address, err)
	}
	acc, err := k.GetAccount(ctx, addrBz)
	if err != nil {
		return nil, err
	}
	return &types.QueryAccountResponse{Account: &acc}, nil
}

func (k queryServer) Proposals(ctx context.Context, req *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	accountAddr, err := sdk.AccAddressFromBech32(req.AccountAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address %s: %v", req.AccountAddress, err)
	}
	if _, err := k.GetAccount(ctx, accountAddr); err != nil {
		return nil, err
	}
	rng := collections.NewPrefixedPairRange[[]byte, uint64](accountAddr.Bytes())
	it, err := k.Keeper.Proposals.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}
	props, err := it.Values()
	if err != nil {
		return nil, err
	}
	return &types.QueryProposalsResponse{Proposals: props}, nil
}

func (k queryServer) Proposal(ctx context.Context, req *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	accountAddr, err := sdk.AccAddressFromBech32(req.AccountAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address %s: %v", req.AccountAddress, err)
	}
	if _, err := k.GetAccount(ctx, accountAddr); err != nil {
		return nil, err
	}
	prop, err := k.Keeper.GetProposal(ctx, accountAddr, req.ProposalId)
	if err != nil {
		return nil, err
	}
	rng := collections.NewSuperPrefixedTripleRange[[]byte, uint64, []byte](accountAddr.Bytes(), req.ProposalId)
	it, err := k.Keeper.Votes.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}
	votes, err := it.Values()
	if err != nil {
		return nil, err
	}
	return &types.QueryProposalResponse{
		Proposal: prop,
		Votes:    votes,
	}, nil
}
