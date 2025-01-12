package types

import (
	"slices"
)

func (a Account) HasMember(member string) bool {
	return slices.ContainsFunc(a.Members, func(m *Member) bool {
		return m.Address == member
	})
}
