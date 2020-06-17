package tests

import (
	"github.com/sudachen/smwlt/cli"
	"github.com/sudachen/smwlt/fu"
	"gotest.tools/assert"
	"testing"
)

func Test_Main1(t *testing.T) {
	c := cli.CLI()
	c.SetArgs([]string{"-h"})
	err := c.Execute()
	assert.NilError(t, err)
}

func Test_Main2(t *testing.T) {
	c := cli.CLI()
	c.SetArgs([]string{"--unknown-key"})
	assert.Assert(t, PanicWith("unknown flag", func() {
		if err := c.Execute(); err != nil {
			panic(fu.Panic(err))
		}
	}))
}
