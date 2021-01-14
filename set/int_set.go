package set

import "sync"

type IntSet interface {
	Add(v int)
	Values() []int
	Size() int
	Has(v int) bool
}

type intSet struct {
	m   map[int]bool
	mux sync.Mutex
}

func (i *intSet) Has(v int) bool {
	_, ok := i.m[v]
	return ok
}

func (i *intSet) Size() int {
	return len(i.m)
}

func (i *intSet) Add(v int) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if _, ok := i.m[v]; !ok {
		i.m[v] = true
	}
}

func (i *intSet) Values() []int {
	i.mux.Lock()
	defer i.mux.Unlock()
	ret := make([]int, i.Size())
	j := 0
	for k := range i.m {
		ret[j] = k
		j++
	}
	return ret
}

func NewIntSet() IntSet {
	return &intSet{m: make(map[int]bool)}
}
