package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/atomone-hub/atomone/x/multisig/keeper"
	"github.com/atomone-hub/atomone/x/multisig/testutil"
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
	var (
		addr    = "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
		account = types.Account{
			Address:   addr,
			Creator:   addr,
			Threshold: 2,
			Members: []*types.Member{
				{
					Address: "1",
					Weight:  1,
				},
				{
					Address: "2",
					Weight:  2,
				},
			},
		}
	)
	tests := []struct {
		name            string
		address         string
		setup           func(context.Context, *keeper.Keeper)
		expectedErr     string
		expectedAccount types.Account
	}{
		{
			name:        "empty address",
			expectedErr: "rpc error: code = InvalidArgument desc = invalid address : empty address string is not allowed",
		},
		{
			name:        "invalid address",
			address:     "xxx",
			expectedErr: "rpc error: code = InvalidArgument desc = invalid address xxx: decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			name:        "invalid address prefix",
			address:     "atone1j8xu70h426xgds2huz9ljpq4jum0dxgar8nq7t",
			expectedErr: "rpc error: code = InvalidArgument desc = invalid address atone1j8xu70h426xgds2huz9ljpq4jum0dxgar8nq7t: invalid Bech32 prefix; expected cosmos, got atone",
		},
		{
			name:        "address not found",
			address:     addr,
			expectedErr: fmt.Sprintf("rpc error: code = NotFound desc = multisig %s doesn't exist", addr),
		},
		{
			name:    "ok",
			address: addr,
			setup: func(ctx context.Context, k *keeper.Keeper) {
				addrBz := sdk.MustAccAddressFromBech32(addr)
				k.Accounts.Set(ctx, addrBz, account)
			},
			expectedAccount: account,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, _, ctx := testutil.SetupMultisigKeeper(t)
			wctx := sdk.WrapSDKContext(ctx)
			queryServer := keeper.NewQueryServer(*k)
			if tt.setup != nil {
				tt.setup(wctx, k)
			}

			response, err := queryServer.Account(wctx, &types.QueryAccountRequest{
				Address: tt.address,
			})

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &types.QueryAccountResponse{Account: &tt.expectedAccount}, response)
		})
	}
}
