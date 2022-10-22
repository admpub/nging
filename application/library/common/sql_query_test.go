package common

import (
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

func (d *testCacheMap) Put(key string, value interface{}, ttlSeconds int64) error {
	d.Store(key, value)
	return nil
}

func (d *testCacheMap) Del(key string) error {
	d.Delete(key)
	return nil
}

func (d *testCacheMap) Get(key string, value interface{}) error {
	return nil
}

func (d *testCacheMap) Do(key string, recv interface{}, fn func() error, ttlSeconds int64) error {
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
	d.Put(key, recv, ttlSeconds)
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
		q.query(&result, func() error {
			result.Int = i + 100
			return nil
		})
		assert.Equal(t, 100, result.Int)
	}

	for i := 0; i < 3; i++ {
		var results []interface{}
		q.query(&results, func() error {
			results = append(results, &TestData{Number: i + 100})
			return nil
		}, `i`, 0)
		assert.Equal(t, []interface{}{&TestData{Number: 100}}, results)
	}
}
