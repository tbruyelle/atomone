package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/atomone-hub/atomone/x/multisig/types"
)

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {
	gen := types.NewGenesisState(
		types.NewParams(),
	)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gen)
}
