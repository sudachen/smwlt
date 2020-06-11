package wallet

import (
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"time"
)

type Account struct {
	Name    string
	Address types.Address
	Created time.Time
	Private ed25519.PrivateKey
}

type WalletImpl interface {
	Unlock(key string) error
	Lookup(alias string) (Account, bool)
}

type Wallet struct {
	WalletImpl
}

func (wal Wallet) LuckyUnlock(key string) {
	if err := wal.Unlock(key); err != nil {
		panic(fu.Panic(err, 2))
	}
}

func Lookup(alias string, w ...Wallet) (acc Account, exists bool) {
	for _, wal := range w {
		if acc, exists = wal.Lookup(alias); exists {
			return
		}
	}
	return
}

func LuckyLookup(alias string, w ...Wallet) Account {
	acc, exists := Lookup(alias, w...)
	if !exists {
		panic(fu.Panic(fmt.Errorf("there is no account '%v'", alias), 2))
	}
	return acc
}
