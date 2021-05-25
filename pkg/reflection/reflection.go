// Package reflection reflect stuff
package reflection

import "reflect"

func GetType(input interface{}) string {
	if t := reflect.TypeOf(input); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func HasElement(input interface{}, search interface{}) bool {
	theValue := reflect.ValueOf(input)

	if theValue.Kind() == reflect.Slice {
		for i := 0; i < theValue.Len(); i++ {
			if theValue.Index(i).CanSet() {
				if theValue.Index(i).Interface() == search {
					return true
				}
			}
		}
	}

	return false
}
