package sqlbuilder

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	relationCache      = make(map[string]*parsedRelation)
	relationCacheMutex sync.RWMutex
)

func getRelationCache(t reflect.StructField) *parsedRelation {
	relationCacheMutex.RLock()
	//println(`key:`, t.Type.String()+`#`+string(t.Tag))
	r := relationCache[t.Type.String()+`#`+string(t.Tag)]
	relationCacheMutex.RUnlock()
	return r
}

func setRelationCache(t reflect.StructField, r *parsedRelation) {
	relationCacheMutex.Lock()
	relationCache[t.Type.String()+`#`+string(t.Tag)] = r
	relationCacheMutex.Unlock()
}

type kv struct {
	k string
	v interface{}
}

func (s *kv) String() string {
	return fmt.Sprintf(`{k: %+v, v: %v}`, s.k, s.v)
}

type colType struct {
	col    interface{}
	colStr string
	typ    string
}

func (s *colType) String() string {
	return fmt.Sprintf(`{col: %+v, colStr: %v, typ: %v}`,
		s.col, s.colStr, s.typ)
}

type selectorArgs struct {
	orderby []interface{}
	offset  int
	limit   int
	groupby []interface{}
	columns []*colType
}

func (s *selectorArgs) String() string {
	return fmt.Sprintf(`{orderby: %+v, offset: %v, limit: %v, groupby: %+v, columns: %+v}`,
		s.orderby, s.offset, s.limit, s.groupby, s.columns)
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

func (r *parsedRelation) String() string {
	var where interface{}
	if r.where != nil {
		where = *r.where
	}
	return fmt.Sprintf(`{relations: %+v, pipes: %v, where: %+v, selectorArgs: %v}`,
		r.relations, r.pipes, where, r.selectorArgs)
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
