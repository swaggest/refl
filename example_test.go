package refl_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/swaggest/refl"
	fancypath "github.com/swaggest/refl/internal/Fancy-Path"
	"github.com/swaggest/refl/internal/sample"
)

func ExampleGoType() {
	fmt.Println(refl.GoType(reflect.TypeOf(new(fancypath.Sample))))
	fmt.Println(refl.GoType(reflect.TypeOf(new(sample.TestSampleStruct))))

	// Output:
	// *github.com/swaggest/refl/internal/Fancy-Path::fancypath.Sample
	// *github.com/swaggest/refl/internal/sample.TestSampleStruct
}

func ExamplePopulateFieldsFromTags() {
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

	s := schema{}
	tag := reflect.TypeOf(value{}).Field(0).Tag

	err := refl.PopulateFieldsFromTags(&s, tag)
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))

	// Output:
	// {"Title":"Value","Desc":"...","Min":-1.23,"Max":10.1,"Limit":5,"Offset":2,"Deprecated":true,"Required":false}
}
