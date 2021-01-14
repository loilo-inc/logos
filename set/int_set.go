package set

import "sync"

type IntSet interface {
	Add(v int)
	Delete(v int)
	Values() []int
	ForEach(f func(v int) (ok bool))
	Size() int
	Has(v int) bool
}

type intSet struct {
	m sync.Map
}

func (i *intSet) Has(v int) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *intSet) Size() int {
	cnt := 0
	i.m.Range(func(key, value interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

func (i *intSet) Add(v int) {
	i.m.Store(v, struct{}{})
}

func (i *intSet) Delete(v int) {
	i.m.Delete(v)
}

func (i *intSet) ForEach(f func(i int) bool) {
	i.m.Range(func(key, value interface{}) bool {
		return f(key.(int))
	})
}

func (i *intSet) Values() []int {
	var ret []int
	i.m.Range(func(key, value interface{}) bool {
		v := key.(int)
		ret = append(ret, v)
		return true
	})
	return ret
}

func NewIntSet() IntSet {
	return &intSet{}
}
