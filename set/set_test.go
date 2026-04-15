package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCases[T comparable]() []struct {
	set  func(args ...T) Set[T]
	name string
} {
	return []struct {
		set  func(args ...T) Set[T]
		name string
	}{
		{NewSet[T], "normal"},
		{func(args ...T) Set[T] { return NewAsyncSet(args...) }, "async"},
	}
}

func TestSet(t *testing.T) {
	for _, c := range makeCases[int]() {
		t.Run(c.name, func(t *testing.T) {
			t.Run("basic", func(t *testing.T) {
				s := c.set()
				assert.Zero(t, s.Size())
				s.Add(1)
				assert.Equal(t, 1, s.Size())
				assert.True(t, s.Has(1))
				assert.False(t, s.Has(2))
				s.Add(2)
				assert.True(t, s.Has(2))
				assert.Equal(t, 2, s.Size())
				s.Add(2)
				assert.Equal(t, 2, s.Size())
				assert.ElementsMatch(t, []int{1, 2}, s.Values())
				s.Delete(2)
				assert.Equal(t, 1, s.Size())
				assert.Equal(t, []int{1}, s.Values())
			})
			t.Run("ForEach", func(t *testing.T) {
				s := c.set()
				s.Add(1)
				s.Add(2)
				s.ForEach(func(v int) (ok bool) {
					assert.True(t, s.Has(v))
					return true
				})
			})
		})
	}
}

func TestNewSet(t *testing.T) {
	for _, c := range makeCases[int]() {
		t.Run(c.name, func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				s := c.set()
				assert.Equal(t, 0, s.Size())
				assert.Empty(t, s.Values())
			})
			t.Run("with single value", func(t *testing.T) {
				s := c.set(1)
				assert.Equal(t, 1, s.Size())
				assert.True(t, s.Has(1))
				assert.Equal(t, []int{1}, s.Values())
			})
			t.Run("with multiple values", func(t *testing.T) {
				s := c.set(1, 2, 3)
				assert.Equal(t, 3, s.Size())
				assert.True(t, s.Has(1))
				assert.True(t, s.Has(2))
				assert.True(t, s.Has(3))
				assert.ElementsMatch(t, []int{1, 2, 3}, s.Values())
			})
			t.Run("with duplicate values", func(t *testing.T) {
				s := c.set(1, 2, 2, 3, 1)
				assert.Equal(t, 3, s.Size())
				assert.ElementsMatch(t, []int{1, 2, 3}, s.Values())
			})
		})
	}
	for _, c := range makeCases[string]() {
		t.Run("with string type", func(t *testing.T) {
			s := c.set("a", "b", "c")
			assert.Equal(t, 3, s.Size())
			assert.True(t, s.Has("a"))
			assert.ElementsMatch(t, []string{"a", "b", "c"}, s.Values())
		})
	}
}
