package api_v1

import (
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu/errstr"
)

// DefaultFee for transfer in the Spacemesh network
const DefaultFee = 10

// DefaultGasLimit for tranfer in the Spacemesh network
const DefaultGasLimit = 100

/*
Transfer creates transaction and submits it to network
*/
func (c *ClientAgent) Transfer(
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
		err = errstr.Wrap(1, err, "failed to encode inner transaction")
		return
	}
	copy(tx.Signature[:], ed25519.Sign2(key, b))
	if b, err = types.InterfaceToBytes(&tx); err != nil {
		err = errstr.Wrap(1, err, "failed to encode whole transaction")
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

/*
Transfer creates transaction and submits it to network. It panics on error
*/
func (c *ClientAgent) LuckyTransfer(
	amount uint64,
	from types.Address, nonce uint64, key ed25519.PrivateKey,
	to types.Address,
	fee, gasLimit uint64) types.TransactionID {

	txid, err := c.Transfer(amount, from, nonce, key, to, fee, gasLimit)
	if err != nil {
		panic(errstr.Frame(2, err))
	}
	return txid
}
