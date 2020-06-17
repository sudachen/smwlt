package modern

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"strings"
	"time"
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

func now() string {
	return time.Now().UTC().Format("2006-01-02T15-04-05.000") + "Z"
}

func fill(path, name, password, mnemonic string) (wal wallet.Wallet) {
	modern := &modernWallet{name: name, path: path}
	modern.content = fu.JsonMap{Val: map[string]interface{}{
		"meta": map[string]interface{}{
			"displayName": name,
			"created":     now(),
			"netId":       int(0),
			"meta": map[string]interface{}{
				"salt": defaultSalt,
			},
		},
		"crypto": map[string]interface{}{
			"cipher":     "AES-128-CTR",
			"cipherText": "",
		},
		"contacts": []map[string]interface{}{},
	}}
	modern.secret = fu.JsonMap{Val: map[string]interface{}{
		"mnemonic": mnemonic,
		"accounts": []map[string]interface{}{},
	}}
	modern.key = pbkdf2.Key([]byte(password), []byte(defaultSalt), 1000000, 32, sha512.New)
	return wallet.Wallet{modern}
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
	text, err := fu.AesXor(w.key, bs)
	m := fu.JsonMap{Val: map[string]interface{}{}}
	if err = m.Decode(bytes.NewBuffer(text)); err != nil {
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

/*
Save implements WalletImpl interface
*/
func (w *modernWallet) Save() (err error) {
	return fu.SaveWithBackup(w.path, func(wr io.Writer) error {
		return w.content.Encode(wr)
	})
}

/*
NewPair implements WalletImpl interface
*/
func (w *modernWallet) NewPair(alias string) (err error) {
	no := len(w.accounts)
	seed := bip39.NewSeed(w.secret.Value("mnemonic").String(), "")
	key := ed25519.NewDerivedKeyFromSeed(seed[:32], uint64(no), []byte(defaultSalt))
	pub := key.Public().(ed25519.PublicKey)[:]

	return w.AddPair(alias, types.BytesToAddress(pub), key, no)
}

/*
AddPair adds predefined keys pair to the wallet
*/
func (w *modernWallet) AddPair(alias string, address types.Address, key ed25519.PrivateKey, index int) (err error) {
	a := wallet.Account{alias, address, key, wallet.Wallet{w}}

	/*
			It does not write accounts list because wallet can contains additional fields not parsed on load,
		       instead it modifies and writes the JsonMap object of json representation.
			So this code may work even if some parts of wallet format will changed
	*/

	accounts := w.secret.Val["accounts"].([]interface{})
	w.secret.Val["accounts"] = append(accounts, map[string]interface{}{
		"displayName": alias,
		"created":     now(),
		"path":        fmt.Sprintf("0/0/%d", index),
		"publicKey":   a.Address.Hex(),
		"secretKey":   hex.EncodeToString(key[:]),
	})

	bf := bytes.Buffer{}
	if err = w.secret.Encode(&bf); err != nil {
		return
	}
	text, err := fu.AesXor(w.key, bf.Bytes())
	w.content.Map("crypto").Val["cipherText"] = hex.EncodeToString(text)

	return
}

/*
ImportKey implements WalletImpl interface
*/
func (w *modernWallet) ImportKey(alias string, address types.Address, key ed25519.PrivateKey) (err error) {
	return w.AddPair(alias, address, key, len(w.accounts))
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
