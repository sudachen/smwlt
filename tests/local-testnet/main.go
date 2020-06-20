package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	_ = os.Chdir("tests")

	term := bootsrap()
	defer term.Sigterm()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<- c
}


