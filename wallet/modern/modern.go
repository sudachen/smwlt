package modern

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

type ModernWallet struct {
	// Since file content can be extended with additional fields we have not loose unknown information
	// So do load and save using KV map
	content map[string]interface{}
	name    string
	path    string
	// this part is initialized on unlock
	secret   map[string]interface{}
	accounts []wallet.Account
	contacts []wallet.Contact
	key      []byte // it's used to encrypt data back
}

func mapField(m map[string]interface{}, k ...string) (r map[string]interface{}, err error) {
	r = m
	for i, x := range k {
		if z, ok := r[x]; ok {
			if q, ok := z.(map[string]interface{}); ok {
				r = q
			} else {
				err = fmt.Errorf("invalid wallet file, region '%v' is not a map", strings.Join(k[:i+1], "."))
				return
			}
		} else {
			err = fmt.Errorf("invalid wallet file, no '%v' region", strings.Join(k[:i+1], "."))
			return
		}
	}
	return
}

func listField(m map[string]interface{}, k ...string) (r []interface{}, err error) {
	if m, err = mapField(m, k[:len(k)-1]...); err != nil {
		return
	}
	x := k[len(k)-1]
	if z, ok := m[x]; ok {
		if r, ok = z.([]interface{}); !ok {
			err = fmt.Errorf("invalid wallet file, '%v' is not a list", strings.Join(k, "."))
		}
	} else {
		err = fmt.Errorf("invalid wallet file, no '%v' field", strings.Join(k, "."))
	}
	return
}

func stringField(m map[string]interface{}, k ...string) (r string, err error) {
	if m, err = mapField(m, k[:len(k)-1]...); err != nil {
		return
	}
	x := k[len(k)-1]
	if z, ok := m[x]; ok {
		if r, ok = z.(string); !ok {
			err = fmt.Errorf("invalid wallet file, '%v' is not a string", strings.Join(k, "."))
		}
	} else {
		err = fmt.Errorf("invalid wallet file, no '%v' field", strings.Join(k, "."))
	}
	return
}

func (w *ModernWallet) load(path string) (err error) {

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(path)
	if err != nil {
		return
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&w.content); err != nil {
		return
	}

	if w.name, err = stringField(w.content, "meta", "displayName"); err != nil {
		return
	}

	w.path = path
	return
}

/*
List implements WalletImpl interface
*/
func (w *ModernWallet) List() []wallet.Account {
	accs := make([]wallet.Account, len(w.accounts))
	for i, a := range w.accounts {
		accs[i] = a
	}
	return accs
}

/*
List implements WalletImpl interface
*/
func (w *ModernWallet) Path() string {
	return w.path
}

/*
Lookup implements WalletImpl interface
*/
func (w *ModernWallet) Lookup(alias string) (acc wallet.Account, exists bool) {
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
func (w *ModernWallet) Name() string {
	return w.name
}

const defaultSalt = "Spacemesh blockmesh"

/*
Unlock implements WalletImpl interface
*/
func (w *ModernWallet) Unlock(password string) (err error) {
	if w.secret != nil {
		return
	}
	cfrt, err := stringField(w.content, "crypto", "cipher")
	if err != nil {
		return
	}
	if cfrt != "AES-128-CTR" {
		return fmt.Errorf("unknown cipher %v", cfrt)
	}
	text, err := stringField(w.content, "crypto", "cipherText")
	if err != nil {
		return
	}
	bs, err := hex.DecodeString(text)
	if err != nil {
		return
	}
	salt, err := stringField(w.content, "meta", "meta", "salt")
	if err != nil {
		salt = defaultSalt
	}
	w.key = pbkdf2.Key([]byte(password), []byte(salt), 1000000, 32, sha512.New)
	block, err := aes.NewCipher(w.key)
	if err != nil {
		return
	}
	iv := [16]byte{}
	iv[15] = 5
	plane := make([]byte, len(text))
	cipher.NewCTR(block, iv[:]).XORKeyStream(plane, bs)
	m := map[string]interface{}{}
	if err = json.NewDecoder(bytes.NewBuffer(plane)).Decode(&m); err != nil {
		return
	}
	accs, err := listField(m, "accounts")
	if err != nil {
		return
	}
	for _, a := range accs {
		q := a.(map[string]interface{})
		a := wallet.Account{Wallet: wallet.Wallet{w}}
		if a.Name, err = stringField(q, "displayName"); err != nil {
			return
		}
		pubk, e := stringField(q, "publicKey")
		if a.Address, err = types.StringToAddress(pubk); err != nil {
			return fu.Wrap(err, "failed to decode public key")
		}
		prvk, e := stringField(q, "secretKey")
		if e != nil {
			return e
		}
		if a.Private, err = hex.DecodeString(prvk); err != nil {
			return fu.Wrap(err, "failed to decode private key")
		}
		w.accounts = append(w.accounts, a)
	}
	w.secret = m
	return
}
