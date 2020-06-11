package wallet

import (
	"encoding/hex"
	"encoding/json"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"os"
	"strings"
)

/*
It's an implementation of the legacy (CLI_Wallet) non-encrypted wallet.
*/

/*
DefaultAccountsJson specifies use local accounts.json file as an wallet
*/
const DefaultAccountsJson = "accounts.json"

type account struct {
	Account
	// there can be additional information related to wallet logic
}

/*
LegacyWallet implements wallet with WalletImpl interface
*/
type LegacyWallet struct{
	accounts []account
}

func (wal *LegacyWallet) load(path string) (err error) {

	type keys struct {
		PubKey  string `json:"pubkey"`
		PrivKey string `json:"privkey"`
	}
	m := map[string]keys{}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(path)
	if err != nil { return }
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&m); err != nil {
		return
	}

	wal.accounts = make([]account,0,len(m))
	for k,v := range m {
		a := account{Account{Name:k}}
		if a.Address, err = types.StringToAddress(v.PubKey); err != nil {
			return fu.Wrap(err, "failed to decode public key")
		}
		if a.Private, err = hex.DecodeString(v.PrivKey); err != nil {
			return fu.Wrap(err, "failed to decode private key")
		}
		wal.accounts = append(wal.accounts,a)
	}
	return
}

/*
Lookup implements WalletImpl interface
*/
func (w *LegacyWallet) Lookup(alias string) (acc Account, exists bool) {
	alias = strings.ToLower(alias)
	for _, a := range w.accounts {
		if a.Name == alias || strings.HasPrefix(alias,"0x") && strings.HasPrefix(a.Address.Hex(),alias) {
			return a.Account, true
		}
	}
	return
}

/*
Unlock implements WalletImpl interface
*/
func (*LegacyWallet) Unlock(string) error {
	// unencrypted
	return nil
}

/*
Legacy defines legacy wallet options
*/
type Legacy struct {
	Path string
}

/*
Load loads wallet content from the file
*/
func (w Legacy) Load() (wal Wallet, err error) {
	legacy := &LegacyWallet{}
	if err = legacy.load(fu.Fne(w.Path,DefaultAccountsJson)); err != nil {
		return
	}
	wal.WalletImpl = legacy
	return
}

/*
Load loads wallet content from the file. It panics on error
*/
func (w Legacy) LuckyLoad() (wal Wallet) {
	fu.LuckyCall(w.Load,&wal)
	return
}
