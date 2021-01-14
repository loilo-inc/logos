package set

import "sync"

type StringSet interface {
	Add(v string)
	Delete(v string)
	Values() []string
	ForEach(f func(v string) (ok bool))
	Size() int
	Has(v string) bool
}

type stringSet struct {
	m sync.Map
}

func (i *stringSet) Has(v string) bool {
	_, ok := i.m.Load(v)
	return ok
}

func (i *stringSet) Size() int {
	cnt := 0
	i.m.Range(func(key, value interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

func (i *stringSet) Add(v string) {
	i.m.Store(v, struct{}{})
}

func (i *stringSet) Delete(v string) {
	i.m.Delete(v)
}

func (i *stringSet) ForEach(f func(i string) bool) {
	i.m.Range(func(key, value interface{}) bool {
		return f(key.(string))
	})
}

func (i *stringSet) Values() []string {
	var ret []string
	i.m.Range(func(key, value interface{}) bool {
		v := key.(string)
		ret = append(ret, v)
		return true
	})
	return ret
}

func NewStringSet() StringSet {
	return &stringSet{}
}
