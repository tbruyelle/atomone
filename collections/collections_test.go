package collections

import (
	"context"
	io "io"
	"math"
	"testing"

	db "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type testStore struct {
	db db.DB
}

func (testStore) CacheWrap() storetypes.CacheWrap {
	panic("not implemented")
}

func (testStore) CacheWrapWithTrace(w io.Writer, tc storetypes.TraceContext) storetypes.CacheWrap {
	panic("not implemented")
}

func (testStore) GetStoreType() storetypes.StoreType {
	panic("not implemented")
}

func (t testStore) OpenKVStore(ctx context.Context) storetypes.KVStore {
	return t
}

func (t testStore) Get(key []byte) []byte {
	v, err := t.db.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (t testStore) Has(key []byte) bool {
	has, err := t.db.Has(key)
	if err != nil {
		panic(err)
	}
	return has
}

func (t testStore) Set(key, value []byte) {
	err := t.db.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (t testStore) Delete(key []byte) {
	err := t.db.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (t testStore) Iterator(start, end []byte) storetypes.Iterator {
	it, err := t.db.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	return it
}

func (t testStore) ReverseIterator(start, end []byte) storetypes.Iterator {
	it, err := t.db.ReverseIterator(start, end)
	if err != nil {
		panic(err)
	}
	return it
}

var _ storetypes.KVStore = testStore{}

func deps() (*testStore, context.Context) {
	kv := db.NewMemDB()
	return &testStore{kv}, context.Background()
}

func TestPrefix(t *testing.T) {
	t.Run("panics on invalid int", func(t *testing.T) {
		require.Panics(t, func() {
			NewPrefix(math.MaxUint8 + 1)
		})
	})

	t.Run("string", func(t *testing.T) {
		require.Equal(t, []byte("prefix"), NewPrefix("prefix").Bytes())
	})

	t.Run("int", func(t *testing.T) {
		require.Equal(t, []byte{0x1}, NewPrefix(1).Bytes())
	})

	t.Run("[]byte", func(t *testing.T) {
		bytes := []byte("prefix")
		prefix := NewPrefix(bytes)
		require.Equal(t, bytes, prefix.Bytes())
		// assert if modification happen they do not propagate to prefix
		bytes[0] = 0x0
		require.Equal(t, []byte("prefix"), prefix.Bytes())
	})
}
