/*
Package fu implements internal functional utilities
*/
package fu

import "reflect"

/*
Fnf returns first non-false value
*/
func Fnf(a ...bool) bool {
	for _, x := range a {
		if x {
			return true
		}
	}
	return false
}

/*
Fne returns first non-empty string
*/
func Fne(a ...string) string {
	for _, x := range a {
		if len(x) > 0 {
			return x
		}
	}
	return ""
}

/*
Opti returns optional or default int value
*/
func Opti(v int, a ...int) int {
	for _, x := range a {
		return x
	}
	return v
}

/*
Ifs returns string selected by logical expression
*/
func Ifs(expr bool, ontrue, onfalse string) string {
	if expr {
		return ontrue
	}
	return onfalse
}

/*
Option returns first option of the required type
*/
func Option(t interface{}, o interface{}) reflect.Value {
	xs := reflect.ValueOf(o)
	tv := reflect.ValueOf(t)
	for i := 0; i < xs.Len(); i++ {
		x := xs.Index(i)
		if x.Kind() == reflect.Interface {
			x = x.Elem()
		}
		if x.Type() == tv.Type() {
			return x
		}
	}
	return tv
}
