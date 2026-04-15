package set

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsyncSet(t *testing.T) {
	t.Run("atomic", func(t *testing.T) {
		wg := sync.WaitGroup{}
		s := NewAsyncSet[int]()
		wg.Go(func() {
			s.Add(1)
		})
		wg.Go(func() {
			s.Add(2)
		})
		wg.Go(func() {
			s.Size()
		})
		wg.Go(func() {
			s.Has(1)
		})
		wg.Wait()
		assert.True(t, s.Has(1))
		assert.True(t, s.Has(2))
	})
}
