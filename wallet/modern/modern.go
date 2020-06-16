package modern

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

type modernWallet struct {
	// Since file content can be extended with additional fields we have not loose unknown information
	// So do load and save using KV map
	content fu.JsonMap
	name    string
	path    string

	// this part is initialized on unlock
	secret   fu.JsonMap
	accounts []wallet.Account
	contacts []wallet.Contact
	key      []byte // it's used to encrypt data back
}

func onpanic(err *error) {
	if p := recover(); p != nil {
		e := fu.PanicError(p)
		if errors.Is(e, fu.JsonTypeError) {
			*err = fu.Wrapf(e, "wallet is broken: %v", e.Error())
			return
		}
		panic(p)
	}
}

func load(path string) (wal wallet.Wallet, err error) {
	defer onpanic(&err)
	w := &modernWallet{}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(path)
	if err != nil {
		return
	}
	defer r.Close()

	if err = w.content.Decode(r); err != nil {
		return
	}

	w.name = w.content.Map("meta").Value("displayName").String()
	w.path = path
	wal.WalletImpl = w
	return
}

/*
List implements WalletImpl interface
*/
func (w *modernWallet) List() []wallet.Account {
	accs := make([]wallet.Account, len(w.accounts))
	for i, a := range w.accounts {
		accs[i] = a
	}
	return accs
}

/*
Path implements WalletImpl interface
*/
func (w *modernWallet) Path() string {
	return w.path
}

/*
Lookup implements WalletImpl interface
*/
func (w *modernWallet) Lookup(alias string) (acc wallet.Account, exists bool) {
	alias = strings.ToLower(alias)
	for _, a := range w.accounts {
		if strings.ToLower(a.Name) == alias ||
			strings.HasPrefix(alias, "0x") && strings.HasPrefix(strings.ToLower(a.Address.Hex()), alias) {
			return a, true
		}
	}
	return
}

/*
Name implements WalletImpl interface
*/
func (w *modernWallet) Name() string {
	return w.name
}

const defaultSalt = "Spacemesh blockmesh"

/*
Unlock implements WalletImpl interface
*/
func (w *modernWallet) Unlock(password string) (err error) {
	defer onpanic(&err)

	if w.secret.Len() != 0 {
		return
	}
	if cfr := w.content.Map("crypto").Value("cipher").String(); cfr != "AES-128-CTR" {
		return fmt.Errorf("unknown cipher %v", cfr)
	}
	bs := w.content.Map("crypto").Value("cipherText").HexBytes()
	salt := fu.Fne(w.content.Map("meta").Map("meta").Value("salt").String(), defaultSalt)
	w.key = pbkdf2.Key([]byte(password), []byte(salt), 1000000, 32, sha512.New)
	block, err := aes.NewCipher(w.key)
	if err != nil {
		return
	}
	iv := [16]byte{}
	iv[15] = 5
	plane := make([]byte, len(bs))
	cipher.NewCTR(block, iv[:]).XORKeyStream(plane, bs)
	m := fu.JsonMap{Val: map[string]interface{}{}}
	if err = m.Decode(bytes.NewBuffer(plane)); err != nil {
		return
	}
	for _, x := range m.List("accounts").Maps() {
		a := wallet.Account{Wallet: wallet.Wallet{w}}
		a.Name = x.Value("displayName").String()
		pubk := x.Value("publicKey").String()
		if a.Address, err = types.StringToAddress(pubk); err != nil {
			return fu.Wrap(err, "failed to decode public key")
		}
		prvk := x.Value("secretKey").String()
		if a.Private, err = hex.DecodeString(prvk); err != nil {
			return fu.Wrap(err, "failed to decode private key")
		}
		w.accounts = append(w.accounts, a)
	}
	w.secret = m
	return
}
