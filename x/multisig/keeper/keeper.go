package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/atomone-hub/atomone/collections"
	collcodec "github.com/atomone-hub/atomone/collections/codec"
	"github.com/atomone-hub/atomone/x/multisig/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	router   types.Router

	authority string

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	Accounts      collections.Map[[]byte, types.Account]
	AccountNumber collections.Sequence
	// Proposals key: account_address+proposal_id
	Proposals      collections.Map[collections.Pair[[]byte, uint64], types.Proposal]
	ProposalNumber collections.Sequence
	// Votes key: account_address+proposal_id+voter_address
	Votes collections.Map[collections.Triple[[]byte, uint64, []byte], types.Vote]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	router types.Router,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilderFromAccessor(func(ctx context.Context) storetypes.KVStore {
		return sdk.UnwrapSDKContext(ctx).KVStore(storeKey)
	})
	k := &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		router:   router,
		Params: collections.NewItem(
			sb, types.KeyParams, "params", collcodec.CollValue[types.Params](cdc),
		),
		Accounts: collections.NewMap(
			sb, types.KeyAccounts, "accounts", collections.BytesKey,
			collcodec.CollValue[types.Account](cdc),
		),
		AccountNumber: collections.NewSequence(sb, types.KeyAccountNumber, "accounts_number"),
		Proposals: collections.NewMap(
			sb, types.KeyProposals, "proposals",
			collections.PairKeyCodec(collections.BytesKey, collections.Uint64Key),
			collcodec.CollValue[types.Proposal](cdc),
		),
		ProposalNumber: collections.NewSequence(sb, types.KeyProposalNumber, "proposal_number"),
		Votes: collections.NewMap(
			sb, types.KeyVotes, "votes",
			collections.TripleKeyCodec(collections.BytesKey, collections.Uint64Key, collections.BytesKey),
			collcodec.CollValue[types.Vote](cdc),
		),
		authority: authority,
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAccount(ctx context.Context, addr sdk.AccAddress) (types.Account, error) {
	acc, err := k.Accounts.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return types.Account{}, status.Errorf(codes.NotFound, "multisig account %s doesn't exist", addr.String())
	}
	return acc, err
}

func (k Keeper) GetProposal(ctx context.Context, addr sdk.AccAddress, id uint64) (types.Proposal, error) {
	prop, err := k.Proposals.Get(ctx, collections.Join(addr.Bytes(), id))
	if errors.Is(err, collections.ErrNotFound) {
		return types.Proposal{}, status.Errorf(codes.NotFound, "multisig proposal %d doesn't exist", id)
	}
	return prop, err
}

func (k Keeper) SetProposal(ctx context.Context, addr sdk.AccAddress, id uint64, prop types.Proposal) error {
	return k.Proposals.Set(ctx, collections.Join(addr.Bytes(), id), prop)
}

// NOTE copied from x/accounts
// makeAddress creates an address for the given account.
// It uses the creator address to ensure address squatting cannot happen, for example
// assuming creator sends funds to a new account X nobody can front-run that address instantiation
// unless the creator itself sends the tx.
// AddressSeed can be used to create predictable addresses, security guarantees of the above are retained.
// If address seed is not provided, the address is created using the creator and account number.
func (k Keeper) makeAddress(creator []byte, accNum uint64, addressSeed []byte) ([]byte, error) {
	// in case an address seed is provided, we use it to create the address.
	var seed []byte
	if len(addressSeed) > 0 {
		seed = append(creator, addressSeed...)
	} else {
		// otherwise we use the creator and account number to create the address.
		seed = append(creator, binary.BigEndian.AppendUint64(nil, accNum)...)
	}

	moduleAndSeed := append([]byte(types.ModuleName), seed...)

	addr := sha256.Sum256(moduleAndSeed)

	return addr[:], nil
}
