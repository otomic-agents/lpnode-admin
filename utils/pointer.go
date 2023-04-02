package utils

import (
	"reflect"
)

func PointerGetString(v interface{}) string {
	val := ""
	if v == nil {
		return val
	}
	valType := reflect.TypeOf(v)
	if valType.String() == "string" {
		val = v.(string)
	}
	if valType.String() == "*string" {
		ptr := v.(*string)
		if ptr == nil {
			return val
		}
		val = *v.(*string)
	}
	return val
}
func PointerGetInt64(v interface{}) int64 {
	var val int64 = 0
	if v == nil {
		return val
	}
	valType := reflect.TypeOf(v)
	if valType.String() == "int64" {
		val = v.(int64)
	}
	if valType.String() == "*int64" {
		ptr := v.(*int64)
		if ptr == nil {
			return val
		}
		val = *v.(*int64)
	}
	return val
}
