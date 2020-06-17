package legacy

/*
It's an implementation of the legacy (CLI_Wallet) non-encrypted wallet.
*/

import (
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
)

/*
DefaultAccountsJson specifies use local accounts.json file as an wallet
*/
const DefaultAccountsJson = "accounts.json"

/*
Wallet defines legacy wallet options
*/
type Wallet struct {
	Path string
}

/*
Load loads wallet content from the file
*/
func (w Wallet) Load() (wal wallet.Wallet, err error) {
	wal, err = load(fu.Fne(w.Path, DefaultAccountsJson))
	return
}

/*
New creates new wallet
*/
func (w Wallet) New() (wal wallet.Wallet) {
	return fill(fu.Fne(w.Path, DefaultAccountsJson))
}

/*
LuckyLoad loads wallet content from the file. It panics on error
*/
func (w Wallet) LuckyLoad() (wal wallet.Wallet) {
	fu.LuckyCall(w.Load, &wal)
	return
}
