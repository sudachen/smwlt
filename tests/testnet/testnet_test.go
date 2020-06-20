package testnet

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu/errstr"
	api "github.com/sudachen/smwlt/node/api.v1"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/legacy"
	"testing"
)

func onpanic(t *testing.T) {
	if e := recover(); e != nil {
		t.Error(errstr.MessageOf(e))
		t.FailNow()
	}
}

func client(t *testing.T) *api.ClientAgent {
	return api.Client{Verbose: t.Logf,Endpoint: "localhost:19090"}.New()
}

func Test_NodeInfo(t *testing.T) {
	defer onpanic(t)
	c := client(t)
	info := c.LuckyNodeInfo()
	t.Logf("%#v\n", info)
}

func Test_AccountInfo(t *testing.T) {
	defer onpanic(t)
	c := client(t)
	w := legacy.Wallet{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	info := c.LuckyAccountInfo(anton.Address)
	t.Logf("%#v\n", info)
}

/*func Test_Transfer(t *testing.T) {
	defer onpanic(t)
	c := client(t)
	w := legacy.Wallet{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	almog := wallet.LuckyLookup("almog", w)
	anton_info1 := c.LuckyAccountInfo(anton.Address)
	t.Logf("anton: %#v\n", anton_info1.Balance)
	almong_info1 := c.LuckyAccountInfo(almog.Address)
	t.Logf("almog: %#v\n", almong_info1.Balance)
	txid := c.LuckyTransfer(100, anton.Address, anton_info1.Nonce, anton.Private, almog.Address, api.DefaultFee, api.DefaultGasLimit)
	t.Logf("txid: %v\n", txid)

	for {
		txinfo := c.LuckyTransactionInfo(txid)
		assert.Assert(t, txinfo.Status != api.TxRejected)
		if txinfo.Status == api.TxConfirmed {
			break
		}
		time.Sleep(5 * time.Second)
	}

	anton_info2 := c.LuckyAccountInfo(anton.Address)
	t.Logf("anton: %#v\n", anton_info2.Balance)
	almong_info2 := c.LuckyAccountInfo(almog.Address)
	t.Logf("almog: %#v\n", almong_info2.Balance)
}*/

func Test_TxList1(t *testing.T) {
	defer onpanic(t)
	c := client(t)
	w := legacy.Wallet{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	almog := wallet.LuckyLookup("almog", w)
	info := c.LuckyNodeInfo()
	anton_info1 := c.LuckyAccountInfo(anton.Address)
	t.Logf("anton: %v, %v\n", anton_info1.Balance, anton_info1.Nonce)
	almong_info1 := c.LuckyAccountInfo(almog.Address)
	t.Logf("almog: %#v\n", almong_info1.Balance)
	txs1 := []types.TransactionID{}
	for k, x := range []uint64{100, 200, 300, 400, 500} {
		txid := c.LuckyTransfer(x, anton.Address, anton_info1.Nonce+uint64(k), anton.Private, almog.Address, api.DefaultFee, api.DefaultGasLimit)
		t.Logf("txid: %v\n", txid)
		txs1 = append(txs1, txid)
	}
	txs2 := c.LuckyTxList(anton.Address, info.SyncedLayer)
	for i,x := range txs2 {
		//fmt.Println(i,t)
		t.Log(i, c.LuckyTransactionInfo(x))
	}
}
