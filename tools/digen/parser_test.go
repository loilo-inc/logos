package digen_test

import (
	"github.com/loilo-inc/logos/tools/digen"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParseManifest(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		src, _ := ioutil.ReadFile("./testdata/manifest.go")
		m, err := digen.ParseManifest(src)
		assert.Nil(t, err)
		assert.ElementsMatch(t, m.Imports, []string{
			`logosSet "github.com/loilo-inc/logos/set"`,
			`"time"`,
		})
		assert.ElementsMatch(t, []string{"string", "String"}, m.Fields[0])
		assert.Equal(t, []string{"logosSet.StringSet", "SSet"}, m.Fields[1])
		assert.Equal(t, []string{"[]int", "Slice"}, m.Fields[2])
		assert.Equal(t, []string{"*[]int", "SlicePtr"}, m.Fields[3])
		assert.Equal(t, []string{"[1]int", "SliceWithLen"}, m.Fields[4])
		assert.Equal(t, []string{"*[1]int", "SliceWithLenPtr"}, m.Fields[5])
		assert.Equal(t, []string{"[][]int", "SliceOfSlice"}, m.Fields[6])
		assert.Equal(t, []string{"[1][1]int", "SliceOfSliceWithLen"}, m.Fields[7])
		assert.Equal(t, []string{"*[][]int", "SliceOfSLicePtr"}, m.Fields[8])
		assert.Equal(t, []string{"*[1][1]int", "SliceOfSliceWithLenPtr"}, m.Fields[9])
		assert.Equal(t, []string{"time.Time", "TimeStr"}, m.Fields[10])
		assert.Equal(t, []string{"*time.Time", "TimePtr"}, m.Fields[11])
		assert.Equal(t, []string{"interface{}", "Iface"}, m.Fields[12])
		assert.Equal(t, []string{"struct{}", "Struct"}, m.Fields[13])
	})
	t.Run("should return empty manifest if empty", func(t *testing.T) {
		m, err := digen.ParseManifest("package main \n")
		assert.Nil(t, err)
		assert.True(t, m.IsEmpty())
	})
	t.Run("should return err if broken syntax", func(t *testing.T) {
		m, err := digen.ParseManifest("package")
		assert.Nil(t, m)
		assert.Error(t, err)
	})
}
