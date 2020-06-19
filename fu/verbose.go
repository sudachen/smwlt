package fu

import "github.com/sudachen/smwlt/fu/stdio"

var VerboseOpt = false
var VerboseOptP *bool = &VerboseOpt

func Verbose(f string, a ...interface{}) {
	if VerboseOptP != nil && *VerboseOptP {
		stdio.Printfln("# "+f, a...)
	}
}
