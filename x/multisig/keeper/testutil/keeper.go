package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	govtypes "github.com/atomone-hub/atomone/x/gov/types"
	"github.com/atomone-hub/atomone/x/multisig/keeper"
	"github.com/atomone-hub/atomone/x/multisig/types"
)

type Mocks struct{}

func SetupMsgServer(t *testing.T) (types.MsgServer, *keeper.Keeper, Mocks, sdk.Context) {
	t.Helper()
	k, m, ctx := SetupMultisigKeeper(t)
	return keeper.NewMsgServerImpl(*k), k, m, ctx
}

func SetupMultisigKeeper(t *testing.T) (
	*keeper.Keeper,
	Mocks,
	sdk.Context,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	_ = ctrl
	m := Mocks{}

	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	// banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	return keeper.NewKeeper(encCfg.Codec, key, authority), m, ctx
}
