package colltest

import (
	"context"
	"io"

	db "github.com/cosmos/cosmos-db"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type contextStoreKey struct{}

// MockStore returns a mock store.KVStoreService and a mock context.Context.
// They can be used to test collections. The StoreService.NewStoreContext
// can be used to instantiate a new empty KVStore.
func MockStore() (*StoreService, context.Context) {
	kv := db.NewMemDB()
	ctx := context.WithValue(context.Background(), contextStoreKey{}, &testStore{kv})
	return &StoreService{}, ctx
}

type StoreService struct{}

func (s StoreService) OpenKVStore(ctx context.Context) storetypes.KVStore {
	return ctx.Value(contextStoreKey{}).(storetypes.KVStore)
}

func (s StoreService) NewStoreContext() context.Context {
	kv := db.NewMemDB()
	return context.WithValue(context.Background(), contextStoreKey{}, &testStore{kv})
}

type testStore struct {
	db db.DB
}

func (t testStore) CacheWrap() storetypes.CacheWrap {
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
