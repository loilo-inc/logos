package set

import "sync"

type BoolSet interface {
	Add(v bool)
	Values() []bool
	Size() int
	Has(v bool) bool
}

type boolSet struct {
	m   map[bool]bool
	mux sync.Mutex
}

func (i *boolSet) Has(v bool) bool {
	_, ok := i.m[v]
	return ok
}

func (i *boolSet) Size() int {
	return len(i.m)
}

func (i *boolSet) Add(v bool) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if _, ok := i.m[v]; !ok {
		i.m[v] = true
	}
}

func (i *boolSet) Values() []bool {
	i.mux.Lock()
	defer i.mux.Unlock()
	ret := make([]bool, i.Size())
	j := 0
	for k := range i.m {
		ret[j] = k
		j++
	}
	return ret
}

func NewBoolSet() BoolSet {
	return &boolSet{m: make(map[bool]bool)}
}
