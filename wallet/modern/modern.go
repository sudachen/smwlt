package modern

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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
	"os"
	"path/filepath"
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

func NewMnemonic() (mnemonic string, err error) {
	bs := make([]byte, 16)
	if _, err = cryptorand.Read(bs); err != nil {
		return
	}
	return bip39.NewMnemonic(bs)
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

func (w *modernWallet) Save() (err error) {

	if _, e := os.Stat(w.path); e == nil {
		_ = os.Remove(w.path + "~")
		if err = os.Rename(w.path, w.path+"~"); err != nil {
			return fu.Wrapf(err, "failed to backup wallet: %v", err.Error())
		}
	}

	defer func() {
		if err != nil {
			if _, e := os.Stat(w.path); e != os.ErrNotExist {
				_ = os.Rename(w.path+"~", w.path)
			}
		}
	}()

	_ = os.MkdirAll(filepath.Dir(w.path), 0755)
	f, err := os.Create(w.path)
	if err != nil {
		return
	}
	defer f.Close()
	if err = w.content.Encode(f); err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	_ = os.Remove(w.path + "~")
	return
}

func (w *modernWallet) NewPair(alias string) (err error) {
	no := len(w.accounts)
	seed := bip39.NewSeed(w.secret.Value("mnemonic").String(), "")
	key := ed25519.NewDerivedKeyFromSeed(seed[:32], uint64(no), []byte(defaultSalt))
	pub := key.Public().(ed25519.PublicKey)[:]

	return w.AddPair(alias, types.BytesToAddress(pub), key, no)
}

func (w *modernWallet) AddPair(alias string, address types.Address, key ed25519.PrivateKey, index int) (err error) {
	a := wallet.Account{alias, address, key, wallet.Wallet{w}}

	/*
		It does not write accounts list because wallet records can contains additional fields not parsed on load.
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
	block, err := aes.NewCipher(w.key)
	if err != nil {
		return
	}
	iv := [16]byte{}
	iv[15] = 5
	crypted := make([]byte, bf.Len())
	cipher.NewCTR(block, iv[:]).XORKeyStream(crypted, bf.Bytes())
	w.content.Map("crypto").Val["cipherText"] = hex.EncodeToString(crypted)

	return
}

func (w *modernWallet) ImportKey(alias string, address types.Address, key ed25519.PrivateKey) (err error) {
	return w.AddPair(alias, address, key, len(w.accounts))
}
