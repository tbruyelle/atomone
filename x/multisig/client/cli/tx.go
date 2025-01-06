package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/atomone-hub/atomone/x/multisig/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		NewMsgCreateAccountCmd(),
	)
	return cmd
}

const FlagThreshold = "threshold"

// NewCreateAccountCmd implements creating a new multisig account command.
func NewMsgCreateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account --threshold X <addr1,weight1> [addr2,weight2]...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create a new multisig account",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// Get voting address
			from := clientCtx.GetFromAddress()

			// Parse multisig members
			var members []*types.Member
			for _, arg := range args {
				parts := strings.Split(arg, ",")
				weight, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return fmt.Errorf("parse weight %s: %v", parts[1], err)
				}
				members = append(members, &types.Member{
					Address: parts[0],
					Weight:  weight,
				})
			}
			threshold, err := cmd.Flags().GetInt64(FlagThreshold)
			if err != nil {
				return err
			}

			// Build message and broadcast
			msg := &types.MsgCreateAccount{
				Sender:    from.String(),
				Members:   members,
				Threshold: threshold,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			// exec atomoned q tx TX_HASH |jq -r '.events[] | select(.type == "multisig_creation") | .attributes[] | select(.key == "address")|.value'
		},
	}

	cmd.Flags().Int64(FlagThreshold, 0, "Specify the threshold required to pass proposal within the multisig account.")
	cmd.MarkFlagRequired(FlagThreshold)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
