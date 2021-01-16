package digen_test

import (
	"github.com/loilo-inc/logos/tools/digen"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	filesEqual := func(t *testing.T, a string, b string) {
		_a, err := ioutil.ReadFile(a)
		assert.Nil(t, err)
		_b, err := ioutil.ReadFile(b)
		assert.Nil(t, err)
		assert.Equal(t, string(_a), string(_b))
	}
	t.Run("basic", func(t *testing.T) {
		dir := path.Join(t.TempDir(), "testdata")
		assert.Nil(t, os.Mkdir(dir, os.ModePerm))
		body, err := ioutil.ReadFile("../../di/di.go")
		assert.Nil(t, err)
		err = digen.Generate("./testdata/manifest.go", dir, string(body))
		assert.Nil(t, err)
		filesEqual(t, filepath.Join(dir, "di.go"), "./testdata/di.go")
		filesEqual(t, filepath.Join(dir, "di_ext.go"), "./testdata/di_ext.go")
	})
}
