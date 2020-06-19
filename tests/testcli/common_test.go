package testcli

import (
	"github.com/sudachen/smwlt/cli"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/stdio"
	"github.com/sudachen/smwlt/tests/expect"
	"gotest.tools/assert"
	"testing"
)

func cliTestOnPanic(t *testing.T, pty *expect.Pty, done chan error, passOnFail bool) {
	if e := recover(); e != nil {
		<-done
		if !passOnFail {
			t.Error(errstr.MessageOf(e))
		} else {
			t.Log("(PASS ON FAIL) " + errstr.MessageOf(e))
		}
		pty.Host.SkipRest()
		if passOnFail {
			return
		}
		t.FailNow()
	}
	// if there is no panic
	pty.Host.SkipRest()
	err := <-done
	if err != nil {
		if !passOnFail {
			t.Error(err.Error())
			t.FailNow()
		}
		t.Log("(PASS ON FAIL) " + err.Error())
	} else {
		if passOnFail {
			t.Error("must fail")
			t.FailNow()
		}
	}
}

func execCLI(t *testing.T, target *expect.OsIo, args ...string) (done chan error) {
	done = make(chan error, 1)
	go func() {
		defer target.LocalIo().Reset()
		defer func() {
			t.Log("CLI done")
			target.Close()
			close(done)
		}()
		c := cli.CLI()
		c.SetOutput(target.Output)
		c.SetArgs(args)
		defer func() {
			if e := recover(); e != nil {
				stdio.Println(e)
			}
		}()
		t.Log("execute CLI with ", args)
		if err := c.Execute(); err != nil {
			//panic(err)
			done <- err
		}
	}()
	return
}

func testCLI(t *testing.T, test func(t *testing.T, pty *expect.Pty), args ...string) {
	pty, err := expect.New()
	assert.NilError(t, err)
	defer pty.Close()

	done := execCLI(t, &pty.Target, args...)

	defer cliTestOnPanic(t, pty, done, false)
	test(t, pty)
}

func failCLI(t *testing.T, test func(t *testing.T, pty *expect.Pty), args ...string) {
	pty, err := expect.New()
	assert.NilError(t, err)
	defer pty.Close()

	done := execCLI(t, &pty.Target, args...)

	defer cliTestOnPanic(t, pty, done, true)
	test(t, pty)
}
