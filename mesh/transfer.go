package mesh

import (
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
)

const DefaultFee = 1
const DefaultGasLimit = 100

func (c *ClinetAgent) Transfer(
	amount uint64,
	from types.Address, nonce uint64, key ed25519.PrivateKey,
	to types.Address,
	fee, gasLimit uint64) (txid types.TransactionID, err error) {

	tx := types.Transaction{}
	tx.AccountNonce = nonce
	tx.Amount = amount
	tx.Recipient = to
	tx.GasLimit = gasLimit
	tx.Fee = fee

	b, err := types.InterfaceToBytes(&tx.InnerTransaction)
	if err != nil {
		err = fu.Wrap(err, "failed to encode inner transaction")
		return
	}
	copy(tx.Signature[:], ed25519.Sign2(key, b))
	if b, err = types.InterfaceToBytes(&tx); err != nil {
		err = fu.Wrap(err, "failed to encode whole transaction")
		return
	}

	out := struct {
		Id string `json:"id"`
	}{}
	err = c.post("/submittransaction", &struct {
		Tx []byte `json:"tx"`
	}{b}, &out)
	if err != nil {
		return
	}
	txid = types.TransactionID(types.HexToHash32(out.Id))
	return
}

func (c *ClinetAgent) LuckyTransfer(
	amount uint64,
	from types.Address, nonce uint64, key ed25519.PrivateKey,
	to types.Address,
	fee, gasLimit uint64) types.TransactionID {

	txid, err := c.Transfer(amount, from, nonce, key, to, fee, gasLimit)
	if err != nil {
		panic(fu.Panic(err, 2))
	}
	return txid
}
