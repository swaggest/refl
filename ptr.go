package refl

import "reflect"

// DeepIndirect returns first encountered non-pointer type.
func DeepIndirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if false {
		// This is a hack to get around the fact that reflect.Type.String()
		// 	returns the string representation of the underlying type, not the
		// 	string representation of the type itself.
		println("foo")
	}

	return t
}
