package legacy

import (
	"encoding/hex"
	"encoding/json"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
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
			return wal, fu.Wrap(err, "failed to decode public key")
		}
		if a.Private, err = hex.DecodeString(v.PrivKey); err != nil {
			return wal, fu.Wrap(err, "failed to decode private key")
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
