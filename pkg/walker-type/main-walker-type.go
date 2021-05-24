// Package walkertype walks types
package walkertype

import reflect "reflect"

type Walker struct{}

type WalkerReturn struct {
	Handled bool
}

type WalkerOnWalkMapOptions struct {
	Document TypeMap
	Key      string
	Kind     reflect.Kind
	Value    TypeGeneric
}

type (
	TypeMap     map[string]interface{}
	TypeArray   []interface{}
	TypeGeneric interface{}
)

func (self *Walker) WalkMap(input TypeMap, onWalk func(options *WalkerOnWalkMapOptions) *WalkerReturn) TypeMap {
	for k, v := range input {
		vt := reflect.TypeOf(v)
		kind := vt.Kind()

		resultWalk := onWalk(&WalkerOnWalkMapOptions{
			Document: input,
			Key:      k,
			Value:    v,
			Kind:     kind,
		})

		if !resultWalk.Handled {
			switch vt.Kind() {
			case reflect.Map:
				if mv, ok := v.(TypeMap); ok {
					input[k] = self.WalkMap(mv, onWalk)
				} else {
					panic("error.")
				}

			case reflect.Array, reflect.Slice:
				if mv, ok := v.(TypeArray); ok {
					input[k] = self.WalkArray(mv)
				} else {
					panic("error.")
				}

			default:
				input[k] = self.WalkGeneric(v)
			}
		}
	}

	return input
}

func (self *Walker) WalkArray(x TypeArray) TypeArray {
	return x
}

func (self *Walker) WalkGeneric(x TypeGeneric) TypeGeneric {
	return x
}
