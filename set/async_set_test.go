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
		add := func(v int) {
			wg.Go(func() {
				s.Add(v)
			})
		}
		add(1)
		add(2)
		go s.Size()
		go s.Has(1)
		wg.Wait()
		assert.True(t, s.Has(1))
		assert.True(t, s.Has(2))
	})
}
