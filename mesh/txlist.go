package mesh

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
)

func (c *ClinetAgent) GetTxList(address types.Address, startLayer uint64) (txs []types.TransactionID, err error) {
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

func (c *ClinetAgent) LuckyTxList(address types.Address, startLayer uint64) (txs []types.TransactionID) {
	fu.LuckyCall(c.GetTxList, &txs, address, startLayer)
	return
}
