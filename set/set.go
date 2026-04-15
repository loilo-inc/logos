package set

type Set[T comparable] interface {
	Add(v T)
	Delete(v T)
	Values() []T
	ForEach(f func(v T) (ok bool))
	Size() int
	Has(v T) bool
}

type set[T comparable] struct {
	m map[T]struct{}
}

func (i *set[T]) Has(v T) bool {
	_, ok := i.m[v]
	return ok
}

func (i *set[T]) Size() int {
	return len(i.m)
}

func (i *set[T]) Add(v T) {
	i.m[v] = struct{}{}
}

func (i *set[T]) Delete(v T) {
	delete(i.m, v)
}

func (i *set[T]) ForEach(f func(i T) bool) {
	for k := range i.m {
		if !f(k) {
			break
		}
	}
}

func (i *set[T]) Values() []T {
	var ret []T
	for k := range i.m {
		ret = append(ret, k)
	}
	return ret
}

func NewSet[T comparable](args ...T) Set[T] {
	s := &set[T]{m: make(map[T]struct{})}
	for _, v := range args {
		s.Add(v)
	}
	return s
}
