package set

import "sync"

type Int64Set interface {
	Add(v int64)
	Delete(v int64)
	Values() []int64
	ForEach(f func(v int64) (ok bool))
	Size() int
	Has(v int64) bool
}

type int64Set struct {
	m sync.Map
}

func (i *int64Set) Has(v int64) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *int64Set) Size() int {
	cnt := 0
	i.m.Range(func(key, value interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

func (i *int64Set) Add(v int64) {
	i.m.Store(v, struct{}{})
}

func (i *int64Set) Delete(v int64) {
	i.m.Delete(v)
}

func (i *int64Set) ForEach(f func(i int64) bool) {
	i.m.Range(func(key, value interface{}) bool {
		return f(key.(int64))
	})
}

func (i *int64Set) Values() []int64 {
	var ret []int64
	i.m.Range(func(key, value interface{}) bool {
		v := key.(int64)
		ret = append(ret, v)
		return true
	})
	return ret
}

func NewInt64Set() Int64Set {
	return &int64Set{}
}
