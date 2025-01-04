package indexes

import (
	"context"
	"io"

	db "github.com/cosmos/cosmos-db"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// TODO remove this when we add testStore to core/store.

type testStore struct {
	db db.DB
}

func (t testStore) OpenKVStore(ctx context.Context) storetypes.KVStore {
	return t
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
	t.db.Delete(key)
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

type company struct {
	City string
	Vat  uint64
}
