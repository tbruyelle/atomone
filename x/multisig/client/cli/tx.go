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
		NewMsgCreateMultisigCmd(),
	)
	return cmd
}

const FlagThreshold = "threshold"

// NewCreate implements creating a new multisig command.
func NewMsgCreateMultisigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [addr1,weight1] [addr2,weight2]...",
		Args:  cobra.MinimumNArgs(4),
		Short: "Create a new multisig account",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// Get voting address
			from := clientCtx.GetFromAddress()

			// Parse multisig members
			var (
				members   []*types.Member
				threshold int64
			)
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
				threshold += int64(weight)
			}

			// Build message and broadcast
			msg := &types.MsgCreateMultisig{
				Sender:    from.String(),
				Members:   members,
				Threshold: threshold,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Int64(FlagThreshold, 0, "Specify the threshold (default to sum of members weight)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
