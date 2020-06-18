package wallet

import (
	cryptorand "crypto/rand"
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/tyler-smith/go-bip39"
	"path/filepath"
	"strings"
)

/*
Account structure for the wallet interface
*/
type Account struct {
	Name    string
	Address types.Address
	Private ed25519.PrivateKey
	Wallet  Wallet
}

type Contact struct {
	Name    string
	Address types.Address
}

/*
Wallet implementation
*/
type WalletImpl interface {
	Name() string
	Path() string
	Unlock(key string) error
	Lookup(alias string) (Account, bool)
	List() []Account
	Save() error
	NewPair(alias string) error
	ImportKey(alias string, address types.Address, key ed25519.PrivateKey) error
}

/*
Wallet decorator
*/
type Wallet struct {
	WalletImpl
}

/*
LuckyUnlock unlocks wallet. It panics if failed to unlock
*/
func (wal Wallet) LuckyUnlock(key string) {
	if err := wal.Unlock(key); err != nil {
		panic(fu.Panic(err, 2))
	}
}

/*
DisplayName retruns composition of wallet name and its file
*/
func (wal Wallet) DisplayName() string {
	return fmt.Sprintf("%v(%v)", wal.Name(), filepath.Base(wal.Path()))
}

/*
LuckySave saves wallet. It panics if failed to unlock
*/
func (wal Wallet) LuckySave() {
	if err := wal.Save(); err != nil {
		panic(fu.Panic(err, 2))
	}
}

/*
LuckyNewPair creates keys pair. It panics if failed to unlock
*/
func (wal Wallet) LuckyNewPair(alias string) {
	if err := wal.NewPair(alias); err != nil {
		panic(fu.Panic(err, 2))
	}
}

/*
LuckyImportKey imports key to the wallet. It panics if failed to unlock
*/
func (wal Wallet) LuckyImportKey(alias string, address types.Address, key ed25519.PrivateKey) {
	if err := wal.ImportKey(alias, address, key); err != nil {
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
			v = append(v, fmt.Sprintf("\t%v [%v]\n", a.Name, a.Address.Hex(), a.Wallet.DisplayName()))
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

/*
PublicKey retrives publick key form private
*/
func PublicKey(key ed25519.PrivateKey) ed25519.PublicKey {
	return key.Public().(ed25519.PublicKey)
}

/*
Address retrives Address key form private
*/
func Address(key ed25519.PrivateKey) types.Address {
	return types.BytesToAddress(key.Public().(ed25519.PublicKey))
}

/*
GenPair generates new wkeys pair
*/
func GenPair(no int, mnemonic string, salt string) (address types.Address, key  ed25519.PrivateKey) {
	seed := bip39.NewSeed(mnemonic, "")
	binsalt := []byte(salt)
	key = ed25519.NewDerivedKeyFromSeed(seed[32:], uint64(no), binsalt)
	address = Address(key)
	return
}

/*
NewMnemonic generates new wallet mnemonic
*/
func NewMnemonic() (mnemonic string, err error) {
	bs := make([]byte, 16)
	if _, err = cryptorand.Read(bs); err != nil {
		return
	}
	return bip39.NewMnemonic(bs)
}
