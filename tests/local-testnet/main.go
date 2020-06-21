package main

import (
	"fmt"
	testnet "github.com/sudachen/smwlt/tests/local-testnet/go-testnet"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func usage() {
	fmt.Println("Usage: local-testnet [nodes [miners]]")
}

func main() {

	_ = os.Chdir("tests")
	args := os.Args[1:]
	nodes := 4
	miners := 3

	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil || n <= 0 {
			usage()
			os.Exit(1)
		}
		nodes = n
	}

	if len(args) > 1 {
		n, err := strconv.Atoi(args[1])
		if err != nil || n <= 0 || n > nodes {
			usage()
			os.Exit(1)
		}
		miners = n
	}

	term := testnet.Bootstrap(nodes,miners)
	defer term.Sigterm()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<- c
}


