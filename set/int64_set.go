package set

import "sync"

type Int64Set interface {
	Add(v int64)
	Values() []int64
	Size() int
	Has(v int64) bool
}

type int64Set struct {
	m   map[int64]bool
	mux sync.Mutex
}

func (i *int64Set) Has(v int64) bool {
	_, ok := i.m[v]
	return ok
}

func (i *int64Set) Size() int {
	return len(i.m)
}

func (i *int64Set) Add(v int64) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if _, ok := i.m[v]; !ok {
		i.m[v] = true
	}
}

func (i *int64Set) Values() []int64 {
	i.mux.Lock()
	defer i.mux.Unlock()
	ret := make([]int64, i.Size())
	j := 0
	for k := range i.m {
		ret[j] = k
		j++
	}
	return ret
}

func NewInt64Set() Int64Set {
	return &int64Set{m: make(map[int64]bool)}
}
