// Package walkertype walks types
package walkertype

import (
	reflect "reflect"
)

type OnKind func(*OnKindOptions) *OnKindWalkReturn

// Is called on only key of maps
type OnMapKey func(*OnMapKeyOptions) *OnKindWalkReturn

type OnKindWalkReturn struct {
	Handled bool
}

type WalkOptions struct {
	Object interface{}
	OnKind OnKind
	// Only key of maps
	OnMapKey OnMapKey
}

type OnKindOptions struct {
	CaseKind reflect.Kind
	Copy     reflect.Value
	Original reflect.Value
}

type OnMapKeyOptions struct {
	Key       reflect.Value
	Copy      reflect.Value
	Original  reflect.Value
	ValueCopy reflect.Value
}

// https://gist.github.com/hvoecking/10772475
func Walk(options *WalkOptions) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(options.Object)
	copy := reflect.New(original.Type()).Elem()

	WalkRecursive(&WalkRecursiveOptions{
		Copy:     copy,
		Original: original,
		OnKind:   options.OnKind,
		OnMapKey: options.OnMapKey,
	})

	// Remove the reflection wrapper
	return copy.Interface()
}

type WalkRecursiveOptions struct {
	Copy     reflect.Value
	Original reflect.Value
	OnKind   OnKind
	OnMapKey OnMapKey
}

func WalkRecursive(options *WalkRecursiveOptions) {
	kind := options.Original.Kind()

	if options.OnKind != nil {
		if options.OnKind(&OnKindOptions{
			CaseKind: kind,
			Copy:     options.Copy,
			Original: options.Original,
		}).Handled {
			return
		}
	}

	switch kind {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := options.Original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		options.Copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		WalkRecursive(&WalkRecursiveOptions{
			Copy:     options.Copy.Elem(),
			Original: originalValue,
			OnKind:   options.OnKind,
			OnMapKey: options.OnMapKey,
		})

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := options.Original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		WalkRecursive(&WalkRecursiveOptions{
			Copy:     copyValue,
			Original: originalValue,
			OnKind:   options.OnKind,
			OnMapKey: options.OnMapKey,
		})
		options.Copy.Set(copyValue)

		// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < options.Original.NumField(); i += 1 {
			if options.Original.Field(i).CanSet() {
				WalkRecursive(&WalkRecursiveOptions{
					Copy:     options.Copy.Field(i),
					Original: options.Original.Field(i),
					OnKind:   options.OnKind,
					OnMapKey: options.OnMapKey,
				})
			}
		}
		// If it is a slice we create a new slice and translate each element
	case reflect.Array, reflect.Slice:
		options.Copy.Set(reflect.MakeSlice(options.Original.Type(), options.Original.Len(), options.Original.Cap()))
		for i := 0; i < options.Original.Len(); i++ {
			WalkRecursive(
				&WalkRecursiveOptions{
					Copy:     options.Copy.Index(i),
					Original: options.Original.Index(i),
					OnKind:   options.OnKind,
					OnMapKey: options.OnMapKey,
				},
			)
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		options.Copy.Set(reflect.MakeMap(options.Original.Type()))
		for _, key := range options.Original.MapKeys() {
			originalValue := options.Original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()

			WalkRecursive(&WalkRecursiveOptions{
				Copy:     copyValue,
				Original: originalValue,
				OnKind:   options.OnKind,
				OnMapKey: options.OnMapKey,
			})

			// @TODO fix this
			// if options.OnMapKey != nil {
			// 	if options.OnMapKey(&OnMapKeyOptions{
			// 		Key:       key,
			// 		Copy:      options.Copy,
			// 		Original:  options.Original,
			// 		ValueCopy: copyValue,
			// 	}).Handled {
			// 		return
			// 	}
			// }

			options.Copy.SetMapIndex(key, copyValue)
		}

	// case reflect.String:
	// 	options.Copy.SetString(options.Original.String())

	// And everything else will simply be taken from the original
	default:
		options.Copy.Set(options.Original)
	}
}
