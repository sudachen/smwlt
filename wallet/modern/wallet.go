package modern

import (
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
)

const walletApp = "Spacemesh"

type Wallet struct {
	Path string
	Name string // only for create new wallet
}

/*
Load loads wallet content from the file
*/
func (w Wallet) Load() (wal wallet.Wallet, err error) {
	wal, err = load(w.Path)
	return
}

/*
LuckyLoad loads wallet content from the file. It panics on error
*/
func (w Wallet) LuckyLoad() (wal wallet.Wallet) {
	fu.LuckyCall(w.Load, &wal)
	return
}
