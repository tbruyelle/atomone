package types

import "github.com/atomone-hub/atomone/collections"

const (
	// ModuleName defines the module name
	ModuleName = "multisig"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	KeyParams    = []byte{0x00} // TODO migrate to collections
	KeyMultisigs = collections.NewPrefix(1)
)
