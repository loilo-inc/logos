//go:generate go run github.com/loilo-inc/logos/cmd/digen $GOFILE
package testdata

import (
	logosSet "github.com/loilo-inc/logos/set"
	"time"
)

// di:"manifest"
type Manifest struct {
	String                 string
	SSet                   logosSet.StringSet
	Slice                  []int
	SlicePtr               *[]int
	SliceWithLen           [1]int
	SliceWithLenPtr        *[1]int
	SliceOfSlice           [][]int
	SliceOfSliceWithLen    [1][1]int
	SliceOfSLicePtr        *[][]int
	SliceOfSliceWithLenPtr *[1][1]int
	TimeStr                time.Time
	TimePtr                *time.Time
	Iface                  interface{}
	Struct                 struct{}
}
