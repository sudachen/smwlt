package main

import (
	testnet "github.com/sudachen/smwlt/tests/local-testnet/go-testnet"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	_ = os.Chdir("tests")

	term := testnet.Bootstrap()
	defer term.Sigterm()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<- c
}


