package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	v1 "github.com/atomone-hub/atomone/x/gov/types/v1"
)

func TestGovernor(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	govKeeper, _, _, ctx := setupGovKeeper(t)
	addrs := simtestutil.CreateRandomAccounts(3)
	govAddrs := convertAddrsToGovAddrs(addrs)

	// Add 2 governors
	gov1Desc := v1.NewGovernorDescription("moniker1", "id1", "website1", "sec1", "detail1")
	gov1, err := v1.NewGovernor(govAddrs[0].String(), gov1Desc, time.Now().UTC())
	require.NoError(err)
	gov2Desc := v1.NewGovernorDescription("moniker2", "id2", "website2", "sec2", "detail2")
	gov2, err := v1.NewGovernor(govAddrs[1].String(), gov2Desc, time.Now().UTC())
	gov2.Status = v1.Inactive
	require.NoError(err)
	govKeeper.SetGovernor(ctx, gov1)
	govKeeper.SetGovernor(ctx, gov2)

	// Get gov1
	gov, found := govKeeper.GetGovernor(ctx, govAddrs[0])
	if assert.True(found, "cant find gov1") {
		assert.Equal(gov1, gov)
	}

	// Get gov2
	gov, found = govKeeper.GetGovernor(ctx, govAddrs[1])
	if assert.True(found, "cant find gov2") {
		assert.Equal(gov2, gov)
	}

	// Get all govs
	govs := govKeeper.GetAllGovernors(ctx)
	if assert.Len(govs, 2, "expected 2 governors") {
		// Insert order is not preserved, order is related to the address which is
		// generated randomly, so the order of govs is random.
		for i := 0; i < 2; i++ {
			switch govs[i].GetAddress().String() {
			case gov1.GetAddress().String():
				assert.Equal(gov1, *govs[i])
			case gov2.GetAddress().String():
				assert.Equal(gov2, *govs[i])
			}
		}
	}

	// Get all active govs
	govs = govKeeper.GetAllActiveGovernors(ctx)
	if assert.Len(govs, 1, "expected 1 active governor") {
		assert.Equal(gov1, *govs[0])
	}

	// IterateGovernors
	govs = nil
	govKeeper.IterateGovernors(ctx, func(i int64, govI v1.GovernorI) bool {
		gov := govI.(v1.Governor)
		govs = append(govs, &gov)
		return false
	})
	if assert.Len(govs, 2, "expected 2 governors") {
		for i := 0; i < 2; i++ {
			switch govs[i].GetAddress().String() {
			case gov1.GetAddress().String():
				assert.Equal(gov1, *govs[i])
			case gov2.GetAddress().String():
				assert.Equal(gov2, *govs[i])
			}
		}
	}
}
