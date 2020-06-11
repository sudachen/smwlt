package mesh

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
)

/*
AccountInfo describes account sate
*/
type AccountInfo struct {
	Nonce uint64
	Balance uint64
}

type addressValue struct { Address string `json:"address"` }

/*
GetAccountInfo returns account state information
*/
func (c *ClinetAgent) GetAccountInfo(address types.Address) (acc AccountInfo, err error) {
	a := addressValue { address.Hex() }

	if acc.Nonce, err = c.accountNonce(a); err != nil { return }
	if acc.Balance, err = c.accountBalance(a); err != nil { return }
	return
}

/*
LuckyAccountInfo returns account state information. It panics if error occurred
*/
func (c *ClinetAgent) LuckyAccountInfo(address types.Address) (acc AccountInfo) {
	fu.LuckyCall(c.GetAccountInfo,&acc,address)
	return
}

/*
GetAccountBalance returns account balance
*/
func (c *ClinetAgent) GetAccountBalance(address types.Address) (balance uint64, err error) {
	return c.accountBalance(addressValue { address.Hex() })
}

func (c *ClinetAgent) accountBalance(a addressValue) (balance uint64, err error) {
	return c.getValue64("/balance", a)
}

/*
GetAccountNonce returns account nonce
*/
func (c *ClinetAgent) GetAccountNonce(address types.Address) (nonce uint64, err error) {
	return c.accountNonce(addressValue { address.Hex() })
}

func (c *ClinetAgent) accountNonce(a addressValue) (nonce uint64, err error) {
	return c.getValue64("/nonce", a)
}

