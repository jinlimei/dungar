package utils

import (
	"reflect"
	"runtime"
)

// GetFunctionName will do some fancy reflection stuff and grab the function name
func GetFunctionName(a interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name()
}
