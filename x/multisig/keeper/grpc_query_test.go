package keeper_test

import (
	"testing"

	"github.com/atomone-hub/atomone/x/multisig/keeper/testutil"
	"github.com/atomone-hub/atomone/x/multisig/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, _, ctx := testutil.SetupMultisigKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})

	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
