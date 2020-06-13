package wallet

import (
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"strings"
	"time"
)

/*
Account structure for the wallet interface
*/
type Account struct {
	Name    string
	Address types.Address
	Created time.Time
	Private ed25519.PrivateKey
	Wallet  Wallet
}

/*
Wallet implementation
*/
type WalletImpl interface {
	Name() string
	Unlock(key string) error
	Lookup(alias string) (Account, bool)
	List() []Account
}

/*
Wallet decorator
*/
type Wallet struct {
	WalletImpl
}

/*
Unlock wallet. It panics if failed to unlock
*/
func (wal Wallet) LuckyUnlock(key string) {
	if err := wal.Unlock(key); err != nil {
		panic(fu.Panic(err, 2))
	}
}

/*
Unlock wallets
*/
func Unlock(key string, w ...Wallet) (ok bool) {
	for _, wal := range w {
		ok = (wal.Unlock(key) == nil) || ok
	}
	return
}

/*
Lookup for accounts in wallets
*/
func Lookup(alias string, w ...Wallet) (acc []Account) {
	for _, wal := range w {
		if a, exists := wal.Lookup(alias); exists {
			acc = append(acc, a)
		}
	}
	return
}

/*
Lookup for an account in wallets. It panics if there are more then one account
*/
func LookupOne(alias string, w ...Wallet) (acc Account, exists bool) {
	accs := Lookup(alias, w...)
	if len(accs) > 1 {
		v := []string{}
		for _, a := range accs {
			v = append(v, fmt.Sprintf("\t%v [%v]\n", a.Name, a.Address.Hex(), a.Wallet.Name()))
		}
		panic(fu.Panic(
			fmt.Errorf("Account '%v' is ambiguous:\n"+strings.Join(v, ""),
				alias,
			)))
	}
	if len(accs) > 0 {
		acc = accs[0]
		exists = true
	}
	return
}

/*
Lookup for an account in walltes. It panics if there is not exactly one account
*/
func LuckyLookup(alias string, w ...Wallet) Account {
	acc, exists := LookupOne(alias, w...)
	if !exists {
		panic(fu.Panic(fmt.Errorf("account '%v' does not exist", alias)))
	}
	return acc
}
