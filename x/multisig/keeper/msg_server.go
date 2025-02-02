package keeper

import (
	"context"
	"errors"
	"fmt"

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

// CreateAccount implements the MsgServer.CreateAccount method.
func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	// TODO: require a Deposit to avoid spam accounts?
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
	num, err := k.AccountNumber.Next(goCtx)
	if err != nil {
		return nil, err
	}
	// create account address
	creator, _ := sdk.AccAddressFromBech32(msg.Sender) // error checked in msg.ValidateBasic
	accountAddr, err := k.makeAddress(creator, num, nil)
	if err != nil {
		return nil, err
	}
	// store account
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	accountAddrStr := sdk.MustBech32ifyAddressBytes(prefix, accountAddr)
	err = k.Accounts.Set(goCtx, accountAddr, types.Account{
		Address:   accountAddrStr,
		Creator:   msg.Sender,
		Members:   msg.Members,
		Threshold: msg.Threshold,
	})
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeAccountCreation,
			sdk.NewAttribute(types.AttributeKeyAccountAddress, accountAddrStr),
		),
	)
	return &types.MsgCreateAccountResponse{
		Address: accountAddrStr,
	}, nil
}

// CreateProposal implements the MsgServer.CreateProposal method.
func (k msgServer) CreateProposal(goCtx context.Context, msg *types.MsgCreateProposal) (*types.MsgCreateProposalResponse, error) {
	// TODO: require a Deposit to avoid spam proposal?
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Fetch account
	accountAddr, err := sdk.AccAddressFromBech32(msg.AccountAddress)
	if err != nil {
		return nil, err
	}
	acc, err := k.GetAccount(ctx, accountAddr)
	if err != nil {
		return nil, err
	}
	// Ensure msg.Sender is a member's account
	if !acc.HasMember(msg.Sender) {
		return nil, types.ErrNotAMember
	}
	// Check proposal messages
	msgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}
	for _, msg := range msgs {
		// assert that the multisig account is the only signer of the message
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, types.ErrInvalidSigner
		}
		if !signers[0].Equals(accountAddr) {
			return nil, types.ErrInvalidSigner
		}
		// use the msg service router to see that there is a valid route for that
		// message.
		if k.router.Handler(msg) == nil {
			return nil, sdkerrors.Wrap(types.ErrUnroutableProposalMsg, sdk.MsgTypeURL(msg))
		}
	}

	// Store proposal
	proposalID, err := k.AccountNumber.Next(ctx)
	if err != nil {
		return nil, err
	}
	submitTime := ctx.BlockTime()
	prop, err := types.NewProposal(proposalID, msg.AccountAddress, submitTime, msg.Sender, msg.Title, msg.Summary, msgs)
	if err != nil {
		return nil, err
	}
	err = k.SetProposal(ctx, accountAddr, proposalID, prop)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeProposalCreation,
			sdk.NewAttribute(types.AttributeKeyAccountAddress, msg.AccountAddress),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprint(proposalID)),
		),
	)
	// Return proposal id
	return &types.MsgCreateProposalResponse{ProposalId: proposalID}, nil
}

// Vote implements the MsgServer.Vote method.
func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	// find proposal
	accountAddr := sdk.MustAccAddressFromBech32(msg.AccountAddress)
	_, err := k.GetProposal(goCtx, accountAddr, msg.ProposalId)
	if err != nil {
		return nil, err
	}
	// find account
	acc, err := k.GetAccount(goCtx, accountAddr)
	if err != nil {
		return nil, err
	}
	// check voter is part of account's members
	if !acc.HasMember(msg.Voter) {
		return nil, types.ErrNotAMember
	}
	// Store (or replace) vote
	voterAddr := sdk.MustAccAddressFromBech32(msg.Voter)
	err = k.SetProposalVote(goCtx, accountAddr, msg.ProposalId, voterAddr,
		types.Vote{
			AccountAddress: msg.AccountAddress,
			ProposalId:     msg.ProposalId,
			VoterAddress:   msg.Voter,
			Vote:           msg.Vote,
		})
	if err != nil {
		return nil, err
	}
	return &types.MsgVoteResponse{}, nil
}

// ExecuteProposal implements the MsgServer.ExecuteProposal method.
func (k msgServer) ExecuteProposal(goCtx context.Context, msg *types.MsgExecuteProposal) (*types.MsgExecuteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// find proposal
	accountAddr := sdk.MustAccAddressFromBech32(msg.AccountAddress)
	proposal, err := k.GetProposal(ctx, accountAddr, msg.ProposalId)
	if err != nil {
		return nil, err
	}
	// find account
	acc, err := k.GetAccount(ctx, accountAddr)
	if err != nil {
		return nil, err
	}
	// check voter is part of account's members
	if !acc.HasMember(msg.Executor) {
		return nil, types.ErrNotAMember
	}

	// tally votes and check threshold
	membersWeights := make(map[string]uint64)
	for _, m := range acc.Members {
		membersWeights[m.Address] = m.Weight
	}
	votes, err := k.Keeper.GetProposalVotes(ctx, accountAddr, msg.ProposalId)
	if err != nil {
		return nil, err
	}
	var yesWeight uint64
	for _, v := range votes {
		if v.Vote == types.VoteOption_VOTE_OPTION_YES {
			weight, ok := membersWeights[v.VoterAddress]
			if !ok {
				panic("should not happen, voter is not part of account members")
			}
			yesWeight += weight
		}
	}
	if yesWeight < acc.Threshold {
		return nil, sdkerrors.Wrapf(types.ErrExecuteWoThreshold, "account threshold: %d, current proposal yes vote weight: %d", acc.Threshold, yesWeight)
	}

	// Execute proposal
	execTime := ctx.BlockTime()
	proposal.ExecTime = &execTime
	msgs, err := proposal.GetMsgs()
	if err != nil {
		return nil, err
	}
	var (
		logMsg string
		resp   types.MsgExecuteProposalResponse
	)
	resp.Responses, err = k.executeMsgs(ctx, msgs)
	if err != nil {
		// Error during messages execution, update proposal and report
		logMsg = fmt.Sprintf("execution failed: %v", err)
		proposal.Status = types.StatusFailed
		resp.Error = err.Error()
	} else {
		// Messages execution success, update proposal status
		logMsg = "execution success"
		proposal.Status = types.StatusPassed
		// Delete votes
		// NOTE(tb): we might want to keep the votes for DAO traceability.
		// Maybe the best scenario is to keep the votes and the proposal only if
		// the account is identified as a DAO.
		if err := k.RemoveProposalVotes(ctx, accountAddr, msg.ProposalId); err != nil {
			return nil, err
		}
	}
	if err := k.SetProposal(ctx, accountAddr, proposal.Id, proposal); err != nil {
		return nil, err
	}
	ctx.Logger().Info("proposal executed", "proposalId", proposal.Id, "results", logMsg)
	return &resp, nil
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
	if err := k.Params.Set(goCtx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
