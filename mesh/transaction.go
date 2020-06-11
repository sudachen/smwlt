package mesh

import (
	"errors"
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
	"strings"
	"time"
)

type TxStatus int
const (
	TxNotfound TxStatus = iota
	TxRejected
	TxPending
	TxConfirmed
)

func (s TxStatus) String() string {
	switch s {
	case TxNotfound: return "Not Found"
	case TxRejected: return "Rejected"
	case TxPending: return "Pending"
	case TxConfirmed: return "Confirmed"
	}
	return "Unknown"
}

type TransactionInfo struct {
	Id types.TransactionID
	From types.Address
	To types.Address
	Amount uint64 `json:"amount,string"`
	Fee uint64 `json:"fee,string"`
	Status TxStatus
	Timestamp time.Time
	LayerId uint64 `json:"layerId,string"`
}

func (c *ClinetAgent) GetTransaction(txid types.TransactionID) (info TransactionInfo, err error) {
	out := struct{
		*TransactionInfo
		IdStr struct{S []byte `json:"id"`} `json:"txId"`
		FromStr struct{S string `json:"address"`} `json:"sender"`
		ToStr struct{S string `json:"address"`} `json:"receiver"`
		TsStr int64 `json:"timestamp,string"`
		StatusStr string `json:"status"`
	}{TransactionInfo: &info}

	err = c.post("/gettransaction",&struct{Id []byte `json:"id"`}{txid.Bytes()},&out)
	if e := errors.Unwrap(err); e != nil {
		if strings.Contains(e.Error(),"not found") {
			return TransactionInfo{Status: TxNotfound}, nil
		}
		return
	}

	info.Id = types.TransactionID(types.BytesToHash(out.IdStr.S))
	info.From = types.HexToAddress(out.FromStr.S)
	info.To = types.HexToAddress(out.ToStr.S)
	info.Timestamp = time.Unix(out.TsStr,0)

	switch out.StatusStr {
	case "PENDING": info.Status = TxPending
	case "REJECTED": info.Status = TxRejected
	case "CONFIRMED": info.Status = TxConfirmed
	default:
		err = fmt.Errorf("unexpected transaction status %v", out.StatusStr)
		return
	}

	return
}

func (c *ClinetAgent) LuckyTransaction(txid types.TransactionID) (info TransactionInfo) {
	fu.LuckyCall(c.GetTransaction,&info,txid)
	return
}

