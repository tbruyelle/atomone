package keeper

import (
	"context"
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

	schema    collections.Schema
	multisigs collections.Map[[]byte, types.Multisig]

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
		authority: authority,
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
