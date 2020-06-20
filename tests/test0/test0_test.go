package test0

import (
	"fmt"
	api "github.com/sudachen/smwlt/node/api.v1"
	"testing"
	"time"
)

func client(t *testing.T) *api.ClientAgent {
	return api.Client{Verbose: t.Logf, Endpoint:"localhost:19090"}.New()
}

func Test_WaitForTestnet(t *testing.T) {
	c := client(t)
	left := 15
	for {
		if _,err := c.GetNodeInfo(); err == nil {
			break
		}
		left--
		if left == 0 {
			t.FailNow()
		}
		fmt.Printf("local-testnet is not available, %v attempt left\n", left)
		time.Sleep(time.Second*2)
	}
}
