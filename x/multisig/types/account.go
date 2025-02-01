package types

import (
	"slices"
)

// HasMember returns true if the account has a member with address `addr`.
func (a Account) HasMember(addr string) bool {
	return slices.ContainsFunc(a.Members, func(m Member) bool {
		return m.Address == addr
	})
}
