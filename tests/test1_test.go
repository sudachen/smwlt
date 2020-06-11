package tests

import (
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/mesh"
	"github.com/sudachen/smwlt/wallet"
	"gotest.tools/assert"
	"testing"
	"time"
)

func Test_NodeInfo(t *testing.T) {
	c := mesh.Client{Verbose: true}.New()
	info := c.LuckyNodeInfo()
	fmt.Printf("%#v\n", info)
}

func Test_AccountInfo(t *testing.T) {
	c := mesh.Client{Verbose: true}.New()
	w := wallet.Legacy{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	info := c.LuckyAccountInfo(anton.Address)
	fmt.Printf("%#v\n", info)
}

func Test_Transfer(t *testing.T) {
	c := mesh.Client{Verbose: true}.New()
	w := wallet.Legacy{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	almog := wallet.LuckyLookup("almog", w)
	anton_info1 := c.LuckyAccountInfo(anton.Address)
	fmt.Printf("anton: %#v\n", anton_info1.Balance)
	almong_info1 := c.LuckyAccountInfo(almog.Address)
	fmt.Printf("almog: %#v\n", almong_info1.Balance)
	txid := c.LuckyTransfer(100, anton.Address, anton_info1.Nonce, anton.Private, almog.Address, mesh.DefaultFee, mesh.DefaultGasLimit)
	fmt.Printf("txid: %v\n", txid)

	for {
		txinfo := c.LuckyTransaction(txid)
		assert.Assert(t, txinfo.Status != mesh.TxRejected)
		if txinfo.Status == mesh.TxConfirmed {
			break
		}
		time.Sleep(5 * time.Second)
	}

	anton_info2 := c.LuckyAccountInfo(anton.Address)
	fmt.Printf("anton: %#v\n", anton_info2.Balance)
	almong_info2 := c.LuckyAccountInfo(almog.Address)
	fmt.Printf("almog: %#v\n", almong_info2.Balance)
}

func Test_TxList(t *testing.T) {
	c := mesh.Client{Verbose: true}.New()
	w := wallet.Legacy{Path: "../accounts.json"}.LuckyLoad()
	anton := wallet.LuckyLookup("anton", w)
	almog := wallet.LuckyLookup("almog", w)
	info := c.LuckyNodeInfo()
	anton_info1 := c.LuckyAccountInfo(anton.Address)
	fmt.Printf("anton: %v, %v\n", anton_info1.Balance, anton_info1.Nonce)
	almong_info1 := c.LuckyAccountInfo(almog.Address)
	fmt.Printf("almog: %#v\n", almong_info1.Balance)
	txs1 := []types.TransactionID{}
	for k, x := range []uint64{100, 200, 300, 400, 500} {
		txid := c.LuckyTransfer(x, anton.Address, anton_info1.Nonce+uint64(k), anton.Private, almog.Address, mesh.DefaultFee, mesh.DefaultGasLimit)
		fmt.Printf("txid: %v\n", txid)
		txs1 = append(txs1, txid)
	}
	txs2 := c.LuckyTxList(anton.Address, info.VerifiedLayer)
	for i, t := range txs2 {
		//fmt.Println(i,t)
		fmt.Println(i, c.LuckyTransaction(t))
	}
}
