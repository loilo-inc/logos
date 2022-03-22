package set

import "sync"

type Set[T comparable] interface {
	Add(v T)
	Delete(v T)
	Values() []T
	ForEach(f func(v T) (ok bool))
	Size() int
	Has(v T) bool
}

type set[T comparable] struct {
	m sync.Map
}

func (i *set[T]) Has(v T) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *set[T]) Size() int {
	cnt := 0
	i.m.Range(func(key, value interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

func (i *set[T]) Add(v T) {
	i.m.Store(v, struct{}{})
}

func (i *set[T]) Delete(v T) {
	i.m.Delete(v)
}

func (i *set[T]) ForEach(f func(i T) bool) {
	i.m.Range(func(key, value interface{}) bool {
		return f(key.(T))
	})
}

func (i *set[T]) Values() []T {
	var ret []T
	i.m.Range(func(key, value interface{}) bool {
		v := key.(T)
		ret = append(ret, v)
		return true
	})
	return ret
}

func NewSet[T comparable]() Set[T] {
	return &set[T]{}
}
