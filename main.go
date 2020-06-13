package main

import (
	"fmt"
	"github.com/sudachen/smwlt/cli"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
)

func loadWallet(path string, legacy bool, password string) (w []wallet.Wallet) {
	if legacy {
		w = []wallet.Wallet{wallet.Legacy{Path: path}.LuckyLoad()}
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
				fmt.Println(fu.PanicMessage(e))
			}
		}
	}()

	cli.Main()

}
