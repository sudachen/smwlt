package fu

import (
	"github.com/sudachen/smwlt/fu/errstr"
	"reflect"
)

func args(a ...interface{}) []reflect.Value {
	r := make([]reflect.Value, len(a))
	for i, x := range a {
		r[i] = reflect.ValueOf(x)
	}
	return r
}

/*
LuckyCall calls specified function and sets result by specified pointer.
It panics if error occurred
*/
func LuckyCall(f, ret interface{}, a ...interface{}) {
	fv := reflect.ValueOf(f)
	v := fv.Call(args(a...))
	if ret != nil {
		if !v[1].IsNil() {
			e := v[1].Interface().(error)
			panic(errstr.Wrap(2, e, e.Error()))
		}
		reflect.ValueOf(ret).Elem().Set(v[0])
	} else {
		if !v[0].IsNil() {
			e := v[0].Interface().(error)
			panic(errstr.Wrap(2, e, e.Error()))
		}
	}
}
