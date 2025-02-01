package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Router is the interface for a Cosmos SDK module router.
type Router interface {
	Handler(msg sdk.Msg) func(ctx sdk.Context, req sdk.Msg) (*sdk.Result, error)
}
