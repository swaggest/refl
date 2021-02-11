package refl_test

import (
	"encoding/json"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/refl"
)

type (
	structWithEmbedded struct {
		B int `path:"b" json:"-"`
		embedded
	}

	structWithTaggedEmbedded struct {
		B        int `path:"b" json:"-"`
		embedded `json:"emb"`
	}

	structWithIgnoredEmbedded struct {
		B        int `path:"b" json:"-"`
		embedded `json:"-"`
	}

	embedded struct {
		A int `json:"a"`
	}

	structWithInline struct {
		Data struct {
			Deeper struct {
				B int `path:"b" json:"-"`
				embedded
			} `json:"deeper"`
		} `json:"data"`
	}
)

func TestHasTaggedFields(t *testing.T) {
	type AnonymousField struct {
		AnonProp int `json:"anonProp"`
	}

	type mixedStruct struct {
		AnonymousField
		FieldQuery int `query:"fieldQuery"`
		FieldBody  int `json:"fieldBody"`
	}

	assert.True(t, refl.HasTaggedFields(mixedStruct{}, "json"))

	var i interface{ Do() }

	assert.False(t, refl.HasTaggedFields(i, "json"))
	assert.False(t, refl.HasTaggedFields(nil, "json"))

	assert.True(t, refl.HasTaggedFields(new(structWithEmbedded), "json"))
	assert.True(t, refl.HasTaggedFields(new(structWithTaggedEmbedded), "json"))
	assert.False(t, refl.HasTaggedFields(new(structWithIgnoredEmbedded), "json"))

	assert.True(t, refl.HasTaggedFields(new(structWithEmbedded), "path"))
	assert.False(t, refl.HasTaggedFields(new(structWithEmbedded), "query"))

	b, err := json.Marshal(structWithTaggedEmbedded{B: 10, embedded: embedded{A: 20}})
	assert.NoError(t, err)
	assert.Equal(t, `{"emb":{"a":20}}`, string(b))

	b, err = json.Marshal(structWithEmbedded{B: 10, embedded: embedded{A: 20}})
	assert.NoError(t, err)
	assert.Equal(t, `{"a":20}`, string(b))

	b, err = json.Marshal(structWithIgnoredEmbedded{B: 10, embedded: embedded{A: 20}})
	assert.NoError(t, err)
	assert.Equal(t, `{}`, string(b))
}

type schema struct {
	Title      string
	Desc       *string
	Min        *float64
	Max        float64
	Limit      int64
	Offset     *int64
	Deprecated bool
	Required   *bool
}

type value struct {
	Property string `title:"Value" desc:"..." min:"-1.23" max:"10.1" limit:"5" offset:"2" deprecated:"true" required:"f"`
}

func TestPopulateFieldsFromTags(t *testing.T) {
	s := schema{}
	tag := reflect.TypeOf(value{}).Field(0).Tag
	require.NoError(t, refl.PopulateFieldsFromTags(&s, tag))

	assert.Equal(t, "Value", s.Title)
	assert.Equal(t, "...", *s.Desc)
	assert.Equal(t, -1.23, *s.Min)
	assert.Equal(t, 10.1, s.Max)
	assert.Equal(t, int64(5), s.Limit)
	assert.Equal(t, int64(2), *s.Offset)
	assert.Equal(t, true, s.Deprecated)
	assert.Equal(t, false, *s.Required)
}

func BenchmarkPopulateFieldsFromTags(b *testing.B) {
	s := schema{}
	tag := reflect.TypeOf(value{}).Field(0).Tag

	for i := 0; i < b.N; i++ {
		if err := refl.PopulateFieldsFromTags(&s, tag); err != nil {
			b.Fatal(err)
		}
	}
}

func TestFindTaggedName(t *testing.T) {
	se := structWithEmbedded{}

	assert.Equal(t, "a", refl.Tagged(&se, &se.A, "json"))
	assert.Equal(t, "b", refl.Tagged(&se, &se.B, "path"))
	assert.Panics(t, func() {
		assert.Equal(t, "b", refl.Tagged(&se, &se.B, "json"))
	})

	si := structWithInline{}

	assert.Equal(t, "data", refl.Tagged(&si, &si.Data, "json"))
	assert.Equal(t, "deeper", refl.Tagged(&si.Data, &si.Data.Deeper, "json"))
}

func BenchmarkFindTaggedName(b *testing.B) {
	se := structWithEmbedded{}
	si := structWithInline{}

	b.Run("embedded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if a := refl.Tagged(&se, &se.A, "json"); a != "a" {
				b.Fail()
			}
		}
	})

	b.Run("inline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			if data := refl.Tagged(&si, &si.Data, "json"); data != "data" {
				b.Fail()
			}
		}
	})

	b.Run("deep-inline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			if deeper := refl.Tagged(&si.Data, &si.Data.Deeper, "json"); deeper != "deeper" {
				b.Fail()
			}
		}
	})
}

func TestWalkTaggedFields(t *testing.T) {
	type upload struct {
		A struct {
			B int `json:"b"`
		} `formData:"a"`
		Upload1 *multipart.FileHeader `formData:"upload1" description:"Upload with *multipart.FileHeader."`
	}

	var tags []string

	refl.WalkTaggedFields(reflect.ValueOf(new(upload)), func(v reflect.Value, sf reflect.StructField, tag string) {
		refl.WalkTaggedFields(v, func(v reflect.Value, sf reflect.StructField, tag string) {
			tags = append(tags, tag)
		}, "json")
		tags = append(tags, tag)
	}, "formData")

	assert.Equal(t, []string{"b", "a", "upload1"}, tags)
}

func BenchmarkWalkTaggedFields(b *testing.B) {
	type upload struct {
		A struct {
			B int `json:"b"`
		} `formData:"a"`
		Upload1 *multipart.FileHeader `formData:"upload1" description:"Upload with *multipart.FileHeader."`
	}

	for i := 0; i < b.N; i++ {
		refl.WalkTaggedFields(reflect.ValueOf(new(upload)), func(v reflect.Value, sf reflect.StructField, tag string) {
			if tag != "a" && tag != "upload1" {
				b.Fail()
			}
		}, "formData")
	}
}
