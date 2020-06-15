package main

import (
	"fmt"
	"github.com/sudachen/smwlt/cli"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
	legacy2 "github.com/sudachen/smwlt/wallet/legacy"
	"os"
)

func loadWallet(path string, legacy bool, password string) (w []wallet.Wallet) {
	if legacy {
		w = []wallet.Wallet{legacy2.Wallet{Path: path}.LuckyLoad()}
	} else {
		panic(fu.Panic(fmt.Errorf("unsupported wallet type")))
	}
	if password != "" {
		ok := wallet.Unlock(password, w...)
		if !ok {
			panic(fu.Panic(fmt.Errorf("there is nothing to unlock, wrong password(?)")))
		}
	}
	return
}

func main() {

	defer func() {
		if !*cli.OptTrace {
			if e := recover(); e != nil {
				fmt.Fprintln(os.Stderr, fu.PanicMessage(e))
				os.Exit(1)
			}
		}
	}()

	cli.Main()

}
