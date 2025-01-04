package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

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

	schema         collections.Schema
	multisigs      collections.Map[[]byte, types.Multisig]
	multisigNumber collections.Sequence

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilderFromAccessor(func(ctx context.Context) storetypes.KVStore {
		return sdk.UnwrapSDKContext(ctx).KVStore(storeKey)
	})
	k := &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		multisigs: collections.NewMap(
			sb,
			types.KeyMultisigs,
			"multisigs",
			collections.BytesKey,
			collcodec.CollValue[types.Multisig](cdc),
		),
		multisigNumber: collections.NewSequence(sb, types.KeyMultisigNumbger, "multisig_number"),
		authority:      authority,
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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
