package multisig_test

import (
	"testing"

	"github.com/atomone-hub/atomone/x/multisig"
	"github.com/atomone-hub/atomone/x/multisig/testutil"
	"github.com/atomone-hub/atomone/x/multisig/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}
	k, _, ctx := testutil.SetupMultisigKeeper(t)

	multisig.InitGenesis(ctx, *k, genesisState)
	got := multisig.ExportGenesis(ctx, *k)

	require.NotNil(t, got)
	require.Equal(t, genesisState, *got)
}
