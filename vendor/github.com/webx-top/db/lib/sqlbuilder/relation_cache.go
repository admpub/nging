package sqlbuilder

import (
	"reflect"
	"sync"
)

var (
	relationCache      = make(map[reflect.Type]*parsedRelation)
	relationCacheMutex sync.RWMutex
)

func getRelationCache(t reflect.Type) *parsedRelation {
	relationCacheMutex.RLock()
	r := relationCache[t]
	relationCacheMutex.RUnlock()
	return r
}

func setRelationCache(t reflect.Type, r *parsedRelation) {
	relationCacheMutex.Lock()
	relationCache[t] = r
	relationCacheMutex.Unlock()
}

type kv struct {
	k string
	v interface{}
}

type colType struct {
	col    interface{}
	colStr string
	typ    string
}

type selectorArgs struct {
	orderby []interface{}
	offset  int
	limit   int
	groupby []interface{}
	columns []*colType
}

func NewParsedRelation(relations []string, pipes []Pipe) *parsedRelation {
	return &parsedRelation{
		relations: relations,
		pipes:     pipes,
	}
}

type parsedRelation struct {
	relations    []string
	pipes        []Pipe
	where        *[]*kv
	selectorArgs *selectorArgs
	mutex        sync.RWMutex
}

func (r *parsedRelation) setWhere(where *[]*kv) {
	r.mutex.Lock()
	r.where = where
	r.mutex.Unlock()
}

func (r *parsedRelation) Where() (where *[]*kv) {
	r.mutex.RLock()
	where = r.where
	r.mutex.RUnlock()
	return
}

func (r *parsedRelation) SelectorArgs() (selectorArgs *selectorArgs) {
	r.mutex.RLock()
	selectorArgs = r.selectorArgs
	r.mutex.RUnlock()
	return
}

func (r *parsedRelation) setSelectorArgs(selectorArgs *selectorArgs) {
	r.mutex.Lock()
	r.selectorArgs = selectorArgs
	r.mutex.Unlock()
}
