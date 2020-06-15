package verbose

import "fmt"

var VerboseOpt = false
var VerboseOptP *bool = &VerboseOpt

func Printfln(f string, a ...interface{}) {
	if VerboseOptP != nil && *VerboseOptP {
		fmt.Printf("# "+f+"\n", a...)
	}
}
