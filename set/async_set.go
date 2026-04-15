package set

import "sync"

type AsyncSet[T comparable] Set[T]

type asyncSet[T comparable] struct {
	m sync.Map
}

func (i *asyncSet[T]) Has(v T) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *asyncSet[T]) Size() int {
	cnt := 0
	i.m.Range(func(key, value any) bool {
		cnt++
		return true
	})
	return cnt
}

func (i *asyncSet[T]) Add(v T) {
	i.m.Store(v, struct{}{})
}

func (i *asyncSet[T]) Delete(v T) {
	i.m.Delete(v)
}

func (i *asyncSet[T]) ForEach(f func(i T) bool) {
	i.m.Range(func(key, value any) bool {
		return f(key.(T))
	})
}

func (i *asyncSet[T]) Values() []T {
	var ret []T
	i.m.Range(func(key, value any) bool {
		ret = append(ret, key.(T))
		return true
	})
	return ret
}

func NewAsyncSet[T comparable](args ...T) AsyncSet[T] {
	s := &asyncSet[T]{}
	for _, v := range args {
		s.Add(v)
	}
	return s
}
