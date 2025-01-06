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

func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	totalWeight := uint64(0)
	for i := range msg.Members {
		var err error
		totalWeight, err = safeAdd(totalWeight, msg.Members[i].Weight)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrWeightsOverflow, "%v", err)
		}
	}
	// threshold must be less than or equal to the total weight
	if totalWeight < uint64(msg.Threshold) {
		return nil, types.ErrTotalWeightGreaterThanThreshold
	}
	// get the next account number
	num, err := k.accountNumber.Next(goCtx)
	if err != nil {
		return nil, err
	}
	// create account address
	creator, _ := sdk.AccAddressFromBech32(msg.Sender) // error checked in msg.ValidateBasic
	accountAddr, err := k.makeAddress(creator, num, nil)
	if err != nil {
		return nil, err
	}
	if err := k.accounts.Set(goCtx, accountAddr, types.Account{
		Creator:   msg.Sender,
		Members:   msg.Members,
		Threshold: msg.Threshold,
	}); err != nil {
		return nil, err
	}

	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	accountAddrStr := sdk.MustBech32ifyAddressBytes(prefix, accountAddr)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeAccountCreation,
			sdk.NewAttribute(types.AttributeKeyAddress, accountAddrStr),
		),
	)
	return &types.MsgCreateAccountResponse{
		Address: accountAddrStr,
	}, nil
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
	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
