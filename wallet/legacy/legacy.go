package legacy

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/wallet"
	"io"
	"os"
	"strings"
)

type account struct {
	wallet.Account
	// there can be additional information related to wallet logic
}

type legacyWallet struct {
	accounts []account
	path     string
}

func fill(path string) (wal wallet.Wallet) {
	return wallet.Wallet{
		&legacyWallet{
			path: path,
		},
	}
}

func load(path string) (wal wallet.Wallet, err error) {

	w := &legacyWallet{}

	type keys struct {
		PubKey  string `json:"pubkey"`
		PrivKey string `json:"privkey"`
	}
	m := map[string]keys{}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(path)
	if err != nil {
		return
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&m); err != nil {
		return
	}

	w.accounts = make([]account, 0, len(m))
	for k, v := range m {
		a := account{wallet.Account{Name: k, Wallet: wallet.Wallet{w}}}
		if a.Address, err = types.StringToAddress(v.PubKey); err != nil {
			return wal, errstr.Wrap(1, err, "failed to decode public key")
		}
		if a.Private, err = hex.DecodeString(v.PrivKey); err != nil {
			return wal, errstr.Wrap(1, err, "failed to decode private key")
		}
		w.accounts = append(w.accounts, a)
	}
	w.path = path
	wal.WalletImpl = w
	return
}

/*
List implements WalletImpl interface
*/
func (w *legacyWallet) List() []wallet.Account {
	accs := make([]wallet.Account, len(w.accounts))
	for i, a := range w.accounts {
		accs[i] = a.Account
	}
	return accs
}

/*
Path implements WalletImpl interface
*/
func (w *legacyWallet) Path() string {
	return w.path
}

/*
Lookup implements WalletImpl interface
*/
func (w *legacyWallet) Lookup(alias string) (acc wallet.Account, exists bool) {
	alias = strings.ToLower(alias)
	for _, a := range w.accounts {
		if strings.ToLower(a.Name) == alias ||
			strings.HasPrefix(alias, "0x") && strings.HasPrefix(strings.ToLower(a.Address.Hex()), alias) {
			return a.Account, true
		}
	}
	return
}

/*
Name implements WalletImpl interface
*/
func (w *legacyWallet) Name() string {
	return "legacyWallet"
}

/*
Unlock implements WalletImpl interface
*/
func (*legacyWallet) Unlock(string) error {
	// unencrypted
	return nil
}

/*
Save implements WalletImpl interface
*/
func (w *legacyWallet) Save() (err error) {
	return fu.SaveWithBackup(w.path, func(wr io.Writer) error {
		type keys struct {
			PubKey  string `json:"pubkey"`
			PrivKey string `json:"privkey"`
		}
		m := map[string]keys{}
		for _, a := range w.accounts {
			m[a.Name] = keys{hex.EncodeToString(wallet.PublicKey(a.Private)), hex.EncodeToString(a.Private)}
		}
		return json.NewEncoder(wr).Encode(&m)
	})
}

/*
NewPair implements WalletImpl interface
*/
func (w *legacyWallet) NewPair(alias string) (err error) {
	if _, exists := w.Lookup(alias); exists {
		return fmt.Errorf("account '%v' already exists")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return errstr.Wrapf(1, err, "cannot create account: %s", err.Error())
	}
	w.accounts = append(w.accounts, account{wallet.Account{
		Name:    alias,
		Address: types.BytesToAddress(pub[:]),
		Private: priv,
		Wallet:  wallet.Wallet{w},
	}})
	return
}

/*
ImportKey implements WalletImpl interface
*/
func (w *legacyWallet) ImportKey(alias string, address types.Address, key ed25519.PrivateKey) (err error) {
	w.accounts = append(w.accounts, account{wallet.Account{
		Name:    alias,
		Address: address,
		Private: key,
		Wallet:  wallet.Wallet{w},
	}})
	return
}
