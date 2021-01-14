package set_test

import (
	"github.com/loilo-inc/logos/set"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSet(t *testing.T) {
	s := set.NewStringSet()
	assert.Zero(t, s.Size())
	values := []string{"1", "2"}
	val1 := values[0]
	val2 := values[1]
	s.Add(val1)
	assert.Equal(t, 1, s.Size())
	assert.True(t, s.Has(val1))
	assert.False(t, s.Has(val2))
	s.Add(val2)
	assert.True(t, s.Has(val2))
	assert.Equal(t, 2, s.Size())
	s.Add(val2)
	assert.Equal(t, 2, s.Size())
	assert.ElementsMatch(t, []string{val1, val2}, s.Values())
}
