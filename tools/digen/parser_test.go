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
		assert.ElementsMatch(t, []string{"String", "string"}, m.Fields[0])
		assert.Equal(t, []string{"SSet", "logosSet.StringSet"}, m.Fields[1])
		assert.Equal(t, []string{"Slice", "[]int"}, m.Fields[2])
		assert.Equal(t, []string{"SlicePtr", "*[]int"}, m.Fields[3])
		assert.Equal(t, []string{"SliceWithLen", "[1]int"}, m.Fields[4])
		assert.Equal(t, []string{"SliceWithLenPtr", "*[1]int"}, m.Fields[5])
		assert.Equal(t, []string{"SliceOfSlice", "[][]int"}, m.Fields[6])
		assert.Equal(t, []string{"SliceOfSliceWithLen", "[1][1]int"}, m.Fields[7])
		assert.Equal(t, []string{"SliceOfSLicePtr", "*[][]int"}, m.Fields[8])
		assert.Equal(t, []string{"SliceOfSliceWithLenPtr", "*[1][1]int"}, m.Fields[9])
		assert.Equal(t, []string{"TimeStr", "time.Time"}, m.Fields[10])
		assert.Equal(t, []string{"TimePtr", "*time.Time"}, m.Fields[11])
		assert.Equal(t, []string{"Iface", "interface{}"}, m.Fields[12])
		assert.Equal(t, []string{"Struct", "struct{}"}, m.Fields[13])
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
