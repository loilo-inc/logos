package set

import "sync"

type BoolSet interface {
	Add(v bool)
	Delete(v bool)
	Values() []bool
	ForEach(f func(v bool) (ok bool))
	Size() int
	Has(v bool) bool
}

type boolSet struct {
	m sync.Map
}

func (i *boolSet) Has(v bool) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *boolSet) Size() int {
	cnt := 0
	i.m.Range(func(key, value interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

func (i *boolSet) Add(v bool) {
	i.m.Store(v, struct{}{})
}

func (i *boolSet) Delete(v bool) {
	i.m.Delete(v)
}

func (i *boolSet) ForEach(f func(i bool) bool) {
	i.m.Range(func(key, value interface{}) bool {
		return f(key.(bool))
	})
}

func (i *boolSet) Values() []bool {
	var ret []bool
	i.m.Range(func(key, value interface{}) bool {
		v := key.(bool)
		ret = append(ret, v)
		return true
	})
	return ret
}

func NewBoolSet() BoolSet {
	return &boolSet{}
}
