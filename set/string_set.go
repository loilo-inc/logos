package set

import "sync"

type StringSet interface {
	Add(v string)
	Values() []string
	Size() int
	Has(v string) bool
}

type stringSet struct {
	m   map[string]bool
	mux sync.Mutex
}

func (i *stringSet) Has(v string) bool {
	_, ok := i.m[v]
	return ok
}

func (i *stringSet) Size() int {
	return len(i.m)
}

func (i *stringSet) Add(v string) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if _, ok := i.m[v]; !ok {
		i.m[v] = true
	}
}

func (i *stringSet) Values() []string {
	i.mux.Lock()
	defer i.mux.Unlock()
	ret := make([]string, i.Size())
	j := 0
	for k := range i.m {
		ret[j] = k
		j++
	}
	return ret
}

func NewStringSet() StringSet {
	return &stringSet{m: make(map[string]bool)}
}
