package xls

import "reflect"

// isKindOf checks if the given value is of the given kind.
// It returns true if the value is of the given kind.
// It returns false if the value is not of the given kind.
func isKindOf(v interface{}, kind reflect.Kind) bool {
	rv := reflect.ValueOf(v)
	actualValue := reflect.Indirect(rv)
	return actualValue.Kind() == kind
}
