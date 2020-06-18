package main

import (
	"fmt"
	"github.com/sudachen/smwlt/cli"
	"github.com/sudachen/smwlt/fu/errstr"
	"os"
)

func main() {

	defer func() {
		if !*cli.OptTrace {
			if e := recover(); e != nil {
				fmt.Fprintln(os.Stderr, errstr.MessageOf(e))
				os.Exit(1)
			}
		}
	}()

	cli.Main()

}
