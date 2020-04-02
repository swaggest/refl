package refl_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/refl"
	"github.com/swaggest/refl/internal/sample"
)

func TestIsSliceOrMap(t *testing.T) {
	assert.True(t, refl.IsSliceOrMap(new(***[]sample.TestSampleStruct)))
	assert.True(t, refl.IsSliceOrMap(new(***map[string]sample.TestSampleStruct)))
	assert.True(t, refl.IsSliceOrMap([]int{}))
	assert.True(t, refl.IsSliceOrMap(map[int]int{}))
	assert.False(t, refl.IsSliceOrMap(new(***sample.TestSampleStruct)))
	assert.False(t, refl.IsSliceOrMap(nil))
}

func TestIsStruct(t *testing.T) {
	assert.False(t, refl.IsStruct(new(***[]sample.TestSampleStruct)))
	assert.False(t, refl.IsStruct(new(***map[string]sample.TestSampleStruct)))
	assert.False(t, refl.IsStruct([]int{}))
	assert.False(t, refl.IsStruct(map[int]int{}))
	assert.True(t, refl.IsStruct(new(***sample.TestSampleStruct)))
	assert.True(t, refl.IsStruct(sample.TestSampleStruct{}))
	assert.False(t, refl.IsStruct(nil))
}

type Map map[int]int

type S struct {
	Map
}

func TestFindEmbeddedSliceOrMap(t *testing.T) {
	assert.NotNil(t, refl.FindEmbeddedSliceOrMap(S{}))
}

func TestIsZero(t *testing.T) {
	type MyStruct struct {
		S  string
		SS struct {
			I int
		}
	}

	assert.True(t, refl.IsZero(reflect.ValueOf(0)))
	assert.True(t, refl.IsZero(reflect.ValueOf(uint(0))))
	assert.True(t, refl.IsZero(reflect.ValueOf(0.0)))
	assert.True(t, refl.IsZero(reflect.ValueOf(complex128(0.0))))
	assert.True(t, refl.IsZero(reflect.ValueOf("")))
	assert.True(t, refl.IsZero(reflect.ValueOf(false)))
	assert.True(t, refl.IsZero(reflect.ValueOf(([]int)(nil))))

	assert.True(t, refl.IsZero(reflect.ValueOf([2]int{})))
	assert.False(t, refl.IsZero(reflect.ValueOf([2]int{1})))

	assert.True(t, refl.IsZero(reflect.ValueOf(MyStruct{})))
	assert.False(t, refl.IsZero(reflect.ValueOf(MyStruct{S: "s"})))
}
