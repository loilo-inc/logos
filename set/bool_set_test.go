package set_test

import (
	"github.com/loilo-inc/logos/set"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBoolSet(t *testing.T) {
	values := []bool{false, true}
	val1 := values[0]
	val2 := values[1]
	t.Run("basic", func(t *testing.T) {
		s := set.NewBoolSet()
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
		assert.ElementsMatch(t, []bool{val1, val2}, s.Values())
		s.Delete(val2)
		assert.Equal(t, 1, s.Size())
		assert.Equal(t, []bool{val1}, s.Values())
	})
	t.Run("ForEach", func(t *testing.T) {
		s := set.NewBoolSet()
		s.Add(val1)
		s.Add(val2)
		s.ForEach(func(v bool) (ok bool) {
			assert.True(t, s.Has(v))
			return true
		})
	})
	t.Run("atomic", func(t *testing.T) {
		wg := sync.WaitGroup{}
		s := set.NewBoolSet()
		add := func(v bool) {
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
