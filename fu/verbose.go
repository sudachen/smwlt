package fu

import "fmt"

var VerboseOpt = false
var VerboseOptP *bool = &VerboseOpt

func Verbose(f string, a ...interface{}) {
	if VerboseOptP != nil && *VerboseOptP {
		Printfln("# "+f, a...)
	}
}

func Printfln(f string, a ...interface{}) {
	fmt.Printf(f+"\n", a...)
}
