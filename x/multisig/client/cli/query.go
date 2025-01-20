package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/atomone-hub/atomone/x/multisig/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group multisig queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdQueryAccount(),
		CmdQueryProposals(),
		CmdQueryProposal(),
		CmdQueryParams(),
	)
	return cmd
}

func CmdQueryAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account <address>",
		Short: "shows multisig account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.Account(cmd.Context(),
				&types.QueryAccountRequest{Address: accountAddr.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals <account_address>",
		Short: "shows multisig account's proposals",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.Proposals(cmd.Context(),
				&types.QueryProposalsRequest{AccountAddress: accountAddr.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal <account_address> <proposal_id>",
		Short: "shows a single multisig account's proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			proposalID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.Proposal(cmd.Context(),
				&types.QueryProposalRequest{
					AccountAddress: accountAddr.String(),
					ProposalId:     proposalID,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
