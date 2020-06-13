package api_v1

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
)

/*
GetTxList lists network transfers from/to specified account
*/
func (c *ClientAgent) GetTxList(address types.Address, startLayer uint64) (txs []types.TransactionID, err error) {
	type account struct {
		Address string `json:"address"`
	}
	txlst := struct {
		Txs []string `json:"txs"`
	}{}
	err = c.post("/accounttxs",
		&struct {
			Account    account `json:"account"`
			StartLayer uint64  `json:"startLayer"`
		}{account{address.String()}, startLayer},
		&txlst)
	m := map[types.TransactionID]bool{}
	txs = make([]types.TransactionID, 0, len(txlst.Txs))
	for _, x := range txlst.Txs {
		txid := types.TransactionID(types.HexToHash32(x))
		if !m[txid] {
			m[txid] = true
			txs = append(txs, types.TransactionID(types.HexToHash32(x)))
		}
	}
	return
}

/*
LuckyTxList lists network transfers from/to specified account. It panics on error
*/
func (c *ClientAgent) LuckyTxList(address types.Address, startLayer uint64) (txs []types.TransactionID) {
	fu.LuckyCall(c.GetTxList, &txs, address, startLayer)
	return
}
