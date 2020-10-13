package refl_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/refl"
	fancypath "github.com/swaggest/refl/internal/Fancy-Path"
	"github.com/swaggest/refl/internal/sample"
)

type (
	NamedSlice              []string
	NamedMap                map[int]string
	NamedSlicePtr           *[]string
	NamedMapOfNamedSlicePtr *map[int]NamedSlicePtr
)

func TestGoType(t *testing.T) {
	assert.Equal(
		t,
		refl.TypeString("github.com/swaggest/refl/internal/sample.TestSampleStruct"),
		refl.GoType(reflect.TypeOf(sample.TestSampleStruct{})),
	)
	assert.Equal(
		t,
		refl.TypeString("*github.com/swaggest/refl/internal/sample.TestSampleStruct"),
		refl.GoType(reflect.TypeOf(new(sample.TestSampleStruct))),
	)
	assert.Equal(
		t,
		refl.TypeString("*github.com/swaggest/refl/internal/Fancy-Path::fancypath.Sample"),
		refl.GoType(reflect.TypeOf(new(fancypath.Sample))),
	)
	assert.Equal(t, refl.TypeString("github.com/swaggest/refl_test.NamedMapOfNamedSlicePtr"), refl.GoType(reflect.TypeOf(NamedMapOfNamedSlicePtr(nil))))
	assert.Equal(t, refl.TypeString("*github.com/swaggest/refl_test.NamedMapOfNamedSlicePtr"), refl.GoType(reflect.TypeOf(new(NamedMapOfNamedSlicePtr))))

	var nsp NamedSlicePtr

	assert.Equal(t, refl.TypeString("github.com/swaggest/refl_test.NamedSlicePtr"), refl.GoType(reflect.TypeOf(nsp)))
	assert.Equal(t, refl.TypeString("github.com/swaggest/refl_test.NamedSlice"), refl.GoType(reflect.TypeOf(NamedSlice{})))
	assert.Equal(t, refl.TypeString("[]string"), refl.GoType(reflect.TypeOf([]string{})))
	assert.Equal(t, refl.TypeString("github.com/swaggest/refl_test.NamedMap"), refl.GoType(reflect.TypeOf(NamedMap{})))
}
