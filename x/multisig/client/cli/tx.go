package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	msgv1 "cosmossdk.io/api/cosmos/msg/v1"

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
		NewCreateAccountCmd(),
		NewDraftProposalCmd(),
		NewCreateProposalCmd(),
		NewVoteCmd(),
		NewExecuteProposalCmd(),
	)
	return cmd
}

const FlagThreshold = "threshold"

// NewCreateAccountCmd implements creating a new multisig account command.
func NewCreateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account --threshold X <addr1,weight1> [addr2,weight2]...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create a new multisig account",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()
			// Parse multisig members
			var members []types.Member
			for _, arg := range args {
				parts := strings.Split(arg, ",")
				weight, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return fmt.Errorf("parse weight %s: %v", parts[1], err)
				}
				members = append(members, types.Member{
					Address: parts[0],
					Weight:  weight,
				})
			}
			threshold, err := cmd.Flags().GetUint64(FlagThreshold)
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

	cmd.Flags().Uint64(FlagThreshold, 0, "Specify the threshold required to pass proposal within the multisig account.")
	cmd.MarkFlagRequired(FlagThreshold)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCreateProposalCmd implements creating a new multisig account proposal command.
func NewCreateProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-proposal <account_address> <path/to/proposal.json>",
		Args:  cobra.ExactArgs(2),
		Short: "Create a new multisig account proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			bz, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			var p proposal
			if err := json.Unmarshal(bz, &p); err != nil {
				return err
			}
			msgs := make([]sdk.Msg, len(p.Messages))
			for i, anyJSON := range p.Messages {
				var msg sdk.Msg
				err := clientCtx.Codec.UnmarshalInterfaceJSON(anyJSON, &msg)
				if err != nil {
					return err
				}
				msgs[i] = msg
			}

			// Build message and broadcast
			msg := &types.MsgCreateProposal{
				Sender:         from.String(),
				AccountAddress: accountAddr.String(),
				Title:          p.Title,
				Summary:        p.Summary,
			}
			anys, err := sdktx.SetMsgs(msgs)
			if err != nil {
				return err
			}
			msg.Messages = anys
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewVoteCmd implements creating a proposal vote command.
func NewVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote <account_address> <proposal_id> <vote>",
		Args:  cobra.ExactArgs(3),
		Short: "Vote on an account's proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			proposalID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			voteOpt := types.VoteOptionFromString(args[2])
			msg := &types.MsgVote{
				Voter:          from.String(),
				AccountAddress: accountAddr.String(),
				ProposalId:     proposalID,
				Vote:           voteOpt,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewExecuteProposalCmd implements executing a proposal command.
func NewExecuteProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-proposal <account_address> <proposal_id>",
		Args:  cobra.ExactArgs(2),
		Short: "Execute a proposal.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			proposalID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			msg := &types.MsgExecuteProposal{
				Executor:       from.String(),
				AccountAddress: accountAddr.String(),
				ProposalId:     proposalID,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

const (
	draftProposalFileName = "draft_proposal.json"
)

type proposal struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	// Msgs defines an array of sdk.Msgs proto-JSON-encoded as Anys.
	Messages []json.RawMessage `json:"messages,omitempty"`
}

// NewDraftProposalCmd let a user generate a draft proposal.
func NewDraftProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft-proposal <account_address>",
		Args:  cobra.ExactArgs(1),
		Short: "Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			accountAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgPrompt := promptui.Select{
				Label: "Select proposal message type:",
				Items: func() []string {
					msgs := clientCtx.InterfaceRegistry.ListImplementations(sdk.MsgInterfaceProtoName)
					sort.Strings(msgs)
					return msgs
				}(),
			}
			_, msgType, err := msgPrompt.Run()
			if err != nil {
				return fmt.Errorf("failed to prompt proposal message type: %w", err)
			}
			msg, err := sdk.GetMsgFromTypeURL(clientCtx.Codec, msgType)
			if err != nil {
				// should never happen
				panic(err)
			}

			// prompt for title and summary
			titlePrompt := promptui.Prompt{
				Label:    "Enter proposal title",
				Validate: client.ValidatePromptNotEmpty,
			}
			title, err := titlePrompt.Run()
			if err != nil {
				return fmt.Errorf("failed to set proposal title: %w", err)
			}
			summaryPrompt := promptui.Prompt{
				Label:    "Enter proposal summary",
				Validate: client.ValidatePromptNotEmpty,
			}
			summary, err := summaryPrompt.Run()
			if err != nil {
				return fmt.Errorf("failed to set proposal summary: %w", err)
			}

			// set messages field
			signerFieldName, err := getSignerFieldName(msg)
			if err != nil {
				fmt.Printf("cannot determine msg %s signer field name: %v", msgType, err)
			}
			defaultValues := make(map[string]string)
			if signerFieldName != "" {
				defaultValues[signerFieldName] = accountAddr.String()
			}
			msgFilled, err := promptMsgFields(msg, "msg", defaultValues)
			if err != nil {
				return fmt.Errorf("failed to set proposal message fields: %w", err)
			}
			message, err := clientCtx.Codec.MarshalInterfaceJSON(msgFilled)
			if err != nil {
				return fmt.Errorf("failed to marshal proposal message: %w", err)
			}

			p := proposal{
				Title:    title,
				Summary:  summary,
				Messages: []json.RawMessage{message},
			}
			if err := writeFile(draftProposalFileName, p); err != nil {
				return err
			}
			cmd.Printf("%s file has successfully been generated.\n", draftProposalFileName)
			return nil
		},
	}
	return cmd
}

// writeFile writes the input to the file
func writeFile(fileName string, input any) error {
	bz, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal proposal: %w", err)
	}
	return os.WriteFile(fileName, bz, 0o600)
}

func getSignerFieldName(msg sdk.Msg) (string, error) {
	// find signer field using "cosmos.msg.v1.signer" proto extension
	protoDesc := protodesc.ToDescriptorProto(proto.MessageReflect(msg).Descriptor())
	protoExts, err := proto.GetExtension(protoDesc.Options, msgv1.E_Signer)
	if err != nil {
		return "", err
	}
	fieldName := protoExts.([]string)[0]
	return strcase.ToCamel(fieldName), nil
}
