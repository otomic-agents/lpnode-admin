package main

import (
	"fmt"
	"reflect"
)

func GetString(v interface{}) string {
	val := ""
	valType := reflect.TypeOf(v)
	if valType.String() == "string" {
		val = v.(string)
	}
	if valType.String() == "*string" {
		val = *v.(*string)
	}
	return val
}
func main() {
	var i int64 = 100
	t := reflect.TypeOf(&i)
	fmt.Println(t.Name(), "|", t.Kind(), "|", t.String())
}
