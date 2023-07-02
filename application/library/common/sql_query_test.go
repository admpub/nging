package common

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/admpub/null"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/defaults"
)

func NewTestCacheMap() factory.Cacher {
	return &testCacheMap{Map: sync.Map{}}
}

type testCacheMap struct {
	sync.Map
}

func init() {
	factory.SetCacher(NewTestCacheMap())
}

func (d *testCacheMap) Put(ctx context.Context, key string, value interface{}, ttlSeconds int64) error {
	d.Store(key, value)
	return nil
}

func (d *testCacheMap) Del(ctx context.Context, key string) error {
	d.Delete(key)
	return nil
}

func (d *testCacheMap) Get(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (d *testCacheMap) Do(ctx context.Context, key string, recv interface{}, fn func() error, ttlSeconds int64) error {
	fmt.Println(`key ===========================>`, key)
	value, ok := d.Load(key)
	if ok {
		reflect.Indirect(reflect.ValueOf(recv)).Set(reflect.Indirect(reflect.ValueOf(value)))
		// switch r := recv.(type) {
		// case *null.Int:
		// 	*r = *value.(*null.Int)
		// case *[]interface{}:
		// 	*r = *value.(*[]interface{})
		// }
		return nil
	}
	fn()
	d.Put(ctx, key, recv, ttlSeconds)
	return nil
}

type TestData struct {
	Number int
}

func TestSQLQueryWithCache(t *testing.T) {
	ctx := defaults.NewMockContext()
	q := NewSQLQuery(ctx).CacheKey(`test`)

	for i := 0; i < 3; i++ {
		result := null.Int{}
		q.query(fmt.Sprintf(`%T`, result), &result, func() error {
			result.Int = i + 100
			return nil
		})
		assert.Equal(t, 100, result.Int)
	}

	for i := 0; i < 3; i++ {
		var results []interface{}
		q.query(fmt.Sprintf(`%T`, results), &results, func() error {
			results = append(results, &TestData{Number: i + 100})
			return nil
		}, `i`, 0)
		assert.Equal(t, []interface{}{&TestData{Number: 100}}, results)
	}
}
