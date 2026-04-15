package set_test

import (
	"sync"
	"testing"

	"github.com/loilo-inc/logos/v2/set"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	values := []int{1, 2}
	val1 := values[0]
	val2 := values[1]
	t.Run("basic", func(t *testing.T) {
		s := set.NewSet[int]()
		assert.Zero(t, s.Size())
		s.Add(val1)
		assert.Equal(t, 1, s.Size())
		assert.True(t, s.Has(val1))
		assert.False(t, s.Has(val2))
		s.Add(val2)
		assert.True(t, s.Has(val2))
		assert.Equal(t, 2, s.Size())
		s.Add(val2)
		assert.Equal(t, 2, s.Size())
		assert.ElementsMatch(t, []int{val1, val2}, s.Values())
		s.Delete(val2)
		assert.Equal(t, 1, s.Size())
		assert.Equal(t, []int{val1}, s.Values())
	})
	t.Run("ForEach", func(t *testing.T) {
		s := set.NewSet[int]()
		s.Add(val1)
		s.Add(val2)
		s.ForEach(func(v int) (ok bool) {
			assert.True(t, s.Has(v))
			return true
		})
	})
	t.Run("atomic", func(t *testing.T) {
		wg := sync.WaitGroup{}
		s := set.NewSet[int]()
		add := func(v int) {
			wg.Go(func() {
				s.Add(v)
			})
		}
		add(val1)
		add(val2)
		go s.Size()
		go s.Has(val1)
		wg.Wait()
		assert.True(t, s.Has(val1))
		assert.True(t, s.Has(val2))
	})
}
func TestNewSet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := set.NewSet[int]()
		assert.Equal(t, 0, s.Size())
		assert.Empty(t, s.Values())
	})
	t.Run("with single value", func(t *testing.T) {
		s := set.NewSet(1)
		assert.Equal(t, 1, s.Size())
		assert.True(t, s.Has(1))
		assert.Equal(t, []int{1}, s.Values())
	})
	t.Run("with multiple values", func(t *testing.T) {
		s := set.NewSet(1, 2, 3)
		assert.Equal(t, 3, s.Size())
		assert.True(t, s.Has(1))
		assert.True(t, s.Has(2))
		assert.True(t, s.Has(3))
		assert.ElementsMatch(t, []int{1, 2, 3}, s.Values())
	})
	t.Run("with duplicate values", func(t *testing.T) {
		s := set.NewSet(1, 2, 2, 3, 1)
		assert.Equal(t, 3, s.Size())
		assert.ElementsMatch(t, []int{1, 2, 3}, s.Values())
	})
	t.Run("with string type", func(t *testing.T) {
		s := set.NewSet("a", "b", "c")
		assert.Equal(t, 3, s.Size())
		assert.True(t, s.Has("a"))
		assert.ElementsMatch(t, []string{"a", "b", "c"}, s.Values())
	})
}
