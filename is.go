package refl

import (
	"math"
	"reflect"
)

// IsSliceOrMap checks if variable is a slice/array/map or a pointer to it.
func IsSliceOrMap(i interface{}) bool {
	if i == nil {
		return false
	}

	t := DeepIndirect(reflect.TypeOf(i))

	return t.Kind() == reflect.Slice || t.Kind() == reflect.Map || t.Kind() == reflect.Array
}

// IsStruct checks if variable is a struct or a pointer to a struct.
func IsStruct(i interface{}) bool {
	if i == nil {
		return false
	}

	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}

// FindEmbeddedSliceOrMap checks if variable has a slice/array/map or a pointer to it embedded.
func FindEmbeddedSliceOrMap(i interface{}) reflect.Type {
	if i == nil {
		return nil
	}

	t := DeepIndirect(reflect.TypeOf(i))

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			v := reflect.Zero(f.Type).Interface()
			if IsSliceOrMap(v) {
				return f.Type
			}

			if t := FindEmbeddedSliceOrMap(v); t != nil {
				return t
			}
		}
	}

	return nil
}

// IsZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
//
// Adapted from go1.13 reflect.IsZero.
func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()

		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !IsZero(v.Index(i)) {
				return false
			}
		}

		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !IsZero(v.Field(i)) {
				return false
			}
		}

		return true
	case reflect.Invalid:
		panic("reflect.Value.IsZero: " + v.Kind().String())
	default:
		// This should never happen, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic("reflect.Value.IsZero: " + v.Kind().String())
	}
}
