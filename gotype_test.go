package refl_test

import (
	"reflect"
	"testing"

	fancypath "github.com/swaggest/refl/internal/Fancy-Path"
	"github.com/swaggest/refl/internal/sample"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/refl"
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
}