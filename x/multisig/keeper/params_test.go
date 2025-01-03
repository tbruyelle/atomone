package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/atomone-hub/atomone/x/multisig/keeper/testutil"
	"github.com/atomone-hub/atomone/x/multisig/types"
)

func TestGetParams(t *testing.T) {
	k, _, ctx := testutil.SetupMultisigKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
