package set_test

import (
	"github.com/loilo-inc/logos/set"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestIntSet(t *testing.T) {
	values := []int{1, 2}
	val1 := values[0]
	val2 := values[1]
	t.Run("basic", func(t *testing.T) {
		s := set.NewIntSet()
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
		s := set.NewIntSet()
		s.Add(val1)
		s.Add(val2)
		s.ForEach(func(v int) (ok bool) {
			assert.True(t, s.Has(v))
			return true
		})
	})
	t.Run("atomic", func(t *testing.T) {
		wg := sync.WaitGroup{}
		s := set.NewIntSet()
		add := func(v int) {
			wg.Add(1)
			go func() {
				s.Add(v)
				wg.Done()
			}()
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
