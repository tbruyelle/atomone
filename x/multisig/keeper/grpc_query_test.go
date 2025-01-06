package keeper_test

import (
	"testing"

	"github.com/atomone-hub/atomone/x/multisig/keeper"
	"github.com/atomone-hub/atomone/x/multisig/keeper/testutil"
	"github.com/atomone-hub/atomone/x/multisig/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamsQuery(t *testing.T) {
	k, _, ctx := testutil.SetupMultisigKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	k.Params.Set(ctx, params)
	queryServer := keeper.NewQueryServer(*k)

	response, err := queryServer.Params(wctx, &types.QueryParamsRequest{})

	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestAccountQuery(t *testing.T) {
	k, _, ctx := testutil.SetupMultisigKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	queryServer := keeper.NewQueryServer(*k)
	var account types.Account

	response, err := queryServer.Account(wctx, &types.QueryAccountRequest{
		Address: "",
	})

	require.NoError(t, err)
	require.Equal(t, &types.QueryAccountResponse{Account: &account}, response)
}
