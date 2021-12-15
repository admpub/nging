package echo

import "context"

func NewKV(k, v string) *KV {
	return &KV{K: k, V: v}
}

//KV 键值对
type KV struct {
	K  string
	V  string
	H  H           `json:",omitempty" xml:",omitempty"`
	X  interface{} `json:",omitempty" xml:",omitempty"`
	fn func(context.Context) interface{}
}

func (a *KV) SetK(k string) *KV {
	a.K = k
	return a
}

func (a *KV) SetV(v string) *KV {
	a.V = v
	return a
}

func (a *KV) SetKV(k, v string) *KV {
	a.K = k
	a.V = v
	return a
}

func (a *KV) SetH(h H) *KV {
	a.H = h
	return a
}

func (a *KV) SetHKV(k string, v interface{}) *KV {
	if a.H == nil {
		a.H = H{}
	}
	a.H.Set(k, v)
	return a
}

func (a *KV) SetX(x interface{}) *KV {
	a.X = x
	return a
}

func (a *KV) SetFn(fn func(context.Context) interface{}) *KV {
	a.fn = fn
	return a
}

func (a *KV) Fn() func(context.Context) interface{} {
	return a.fn
}

type KVList []*KV

func (list *KVList) Add(k, v string, options ...KVOption) {
	a := &KV{K: k, V: v}
	for _, option := range options {
		option(a)
	}
	*list = append(*list, a)
}

func (list *KVList) AddItem(item *KV) {
	*list = append(*list, item)
}

func (list *KVList) Delete(i int) {
	n := len(*list)
	if i+1 < n {
		*list = append((*list)[0:i], (*list)[i+1:]...)
	} else if i < n {
		*list = (*list)[0:i]
	}
}

func (list *KVList) Reset() {
	*list = (*list)[0:0]
}

//NewKVData 键值对数据
func NewKVData() *KVData {
	return &KVData{
		slice: []*KV{},
		index: map[string][]int{},
	}
}

//KVData 键值对数据（保持顺序）
type KVData struct {
	slice []*KV
	index map[string][]int
}

//Slice 返回切片
func (a *KVData) Slice() []*KV {
	return a.slice
}

//Keys 返回所有K值
func (a *KVData) Keys() []string {
	keys := make([]string, len(a.slice))
	for i, v := range a.slice {
		if v == nil {
			continue
		}
		keys[i] = v.K
	}
	return keys
}

//Index 返回某个key的所有索引值
func (a *KVData) Index(k string) []int {
	v, _ := a.index[k]
	return v
}

//Indexes 返回所有索引值
func (a *KVData) Indexes() map[string][]int {
	return a.index
}

//Reset 重置
func (a *KVData) Reset() *KVData {
	a.index = map[string][]int{}
	a.slice = []*KV{}
	return a
}

//Add 添加键值
func (a *KVData) Add(k, v string, options ...KVOption) *KVData {
	if _, y := a.index[k]; !y {
		a.index[k] = []int{}
	}
	a.index[k] = append(a.index[k], len(a.slice))
	an := &KV{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = append(a.slice, an)
	return a
}

func (a *KVData) AddItem(item *KV) *KVData {
	if _, y := a.index[item.K]; !y {
		a.index[item.K] = []int{}
	}
	a.index[item.K] = append(a.index[item.K], len(a.slice))
	a.slice = append(a.slice, item)
	return a
}

//Set 设置首个键值
func (a *KVData) Set(k, v string, options ...KVOption) *KVData {
	a.index[k] = []int{0}
	an := &KV{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = []*KV{an}
	return a
}

func (a *KVData) SetItem(item *KV) *KVData {
	a.index[item.K] = []int{0}
	a.slice = []*KV{item}
	return a
}

func (a *KVData) Get(k string, defaults ...string) string {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]].V
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

func (a *KVData) GetItem(k string, defaults ...func() *KV) *KV {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]]
		}
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	return nil
}

func (a *KVData) Has(k string) bool {
	if _, ok := a.index[k]; ok {
		return true
	}
	return false
}

//Delete 设置某个键的所有值
func (a *KVData) Delete(ks ...string) *KVData {
	indexes := []int{}
	for _, k := range ks {
		v, y := a.index[k]
		if !y {
			continue
		}
		indexes = append(indexes, v...)
	}
	newSlice := []*KV{}
	a.index = map[string][]int{}
	for i, v := range a.slice {
		var exists bool
		for _, idx := range indexes {
			if i != idx {
				continue
			}
			exists = true
			break
		}
		if exists {
			continue
		}
		if _, y := a.index[v.K]; !y {
			a.index[v.K] = []int{}
		}
		a.index[v.K] = append(a.index[v.K], len(newSlice))
		newSlice = append(newSlice, v)
	}
	a.slice = newSlice
	return a
}
