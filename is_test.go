package refl_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
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

func TestIsScalar(t *testing.T) {
	assert.True(t, refl.IsScalar(123))
	assert.True(t, refl.IsScalar(new(int)))
	assert.True(t, refl.IsScalar(123.4))
	assert.True(t, refl.IsScalar(123.4i))
	assert.True(t, refl.IsScalar(true))
	assert.True(t, refl.IsScalar("abc"))
	assert.True(t, refl.IsScalar(new(string)))
	assert.False(t, refl.IsScalar(struct{}{}))
	assert.False(t, refl.IsScalar(nil))
	assert.False(t, refl.IsScalar(interface{ foo() }(nil)))
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

type doer interface {
	Do()
}

type someDoer struct{ val string }

func (someDoer) Do() {}

type somePtrDoer struct{}

func (*somePtrDoer) Do() {}

func TestAs(t *testing.T) {
	var (
		v  interface{}
		vv doer
	)

	sd := someDoer{val: "abc"}
	vv = sd
	v = vv
	target := new(someDoer)

	assert.True(t, refl.As(v, target))
	assert.Equal(t, sd, *target)
	assert.True(t, refl.As(v, new(doer)))

	spd := &somePtrDoer{}
	vv = spd
	v = vv
	targetP := new(somePtrDoer)

	assert.True(t, refl.As(v, targetP))
	assert.Equal(t, spd, targetP)

	assert.False(t, refl.As("abc", new(doer)))
	assert.False(t, refl.As(new(doer), new(json.Marshaler)))

	assert.False(t, refl.As(nil, target))
	assert.False(t, refl.As(v, new(interface{ Unknown() })))
	assert.Panics(t, func() {
		refl.As(v, someDoer{})
	})
}

func ExampleAs() {
	var (
		v  interface{}
		vv json.Marshaler
	)

	vv = json.RawMessage(`{"abc":123}`)
	v = vv

	target := new(json.RawMessage)
	fmt.Println(refl.As(v, target), string(*target))

	// Output:
	// true {"abc":123}
}

func TestNoEmptyFields(t *testing.T) {
	type My struct {
		Foo int
		Bar string
		Baz chan string
	}

	m := &My{
		Foo: 123,
		Bar: "abc",
		Baz: make(chan string),
	}

	var v interface{} = m

	require.NoError(t, refl.NoEmptyFields(v))

	m.Foo = 0
	m.Bar = ""

	assert.EqualError(t, refl.NoEmptyFields(v), "missing: [Foo Bar]")
	assert.EqualError(t, refl.NoEmptyFields(refl.SentinelError("")), "struct expected, refl.SentinelError received")
}
