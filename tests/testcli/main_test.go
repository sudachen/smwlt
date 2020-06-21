package testcli

import (
	"github.com/sudachen/smwlt/tests/expect"
	"testing"
)

func Test_MainOnly1(t *testing.T) {
	testCLI(t, func(t *testing.T, pty *expect.Pty) {
		pty.Host.Expect(`Spacemesh CLI Wallet`)
	}, "-h")
}

func Test_MainOnly2(t *testing.T) {
	failCLI(t, func(t *testing.T, pty *expect.Pty) {
		pty.Host.Expect(`XXX Spacemesh CLI Wallet`)
	}, "-h")
}

func Test_MainOnly3(t *testing.T) {
	failCLI(t, func(t *testing.T, pty *expect.Pty) {
	}, "--unknown-option")
}
