package zinc

import (
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/mesh"
	"github.com/sudachen/smwlt/wallet"
	"math/rand"
	"testing"
	"time"
)

func Test_Stress1(t *testing.T) {
	ca := make([]*mesh.ClinetAgent, 5)
	ca[0] = mesh.Client{}.New()
	for i := 1; i < 5; i++ {
		ca[i] = mesh.Client{Endpoint: fmt.Sprintf("localhost:919%d", i)}.New()
	}
	c := ca[0]
	w := wallet.Legacy{Path: "../accounts.json"}.LuckyLoad()
	a := []wallet.Account{}
	a = append(a, wallet.LuckyLookup("anton", w))
	a = append(a, wallet.LuckyLookup("almog", w))
	a = append(a, wallet.LuckyLookup("gavrad", w))
	a = append(a, wallet.LuckyLookup("tap", w))
	a = append(a, wallet.LuckyLookup("yosher", w))

	nonce := make([]uint64, len(a))

	for i, x := range a {
		c = ca[rand.Int31n(int32(len(ca)))]
		nfo := c.LuckyAccountInfo(x.Address)
		nonce[i] = nfo.Nonce
	}

	var txid types.TransactionID
	for i := 0; i < 1000; i++ {
		for i := 0; i < 1000; i++ {
			j := rand.Int31n(int32(len(a)))
			from := a[j]
			k := rand.Int31n(int32(len(a)))
			for k == j {
				k = rand.Int31n(int32(len(a)))
			}
			to := a[k]
			c = ca[rand.Int31n(int32(len(ca)))]
			var err error
			txid, err = c.Transfer(1, from.Address, nonce[j], from.Private, to.Address, mesh.DefaultFee, mesh.DefaultGasLimit)
			if err == nil {
				nonce[j]++
			} else {
				nonce[j] = c.LuckyAccountInfo(from.Address).Nonce
			}
		}
		time.Sleep(3 * time.Second)
	}
	for {
		txnfo := ca[0].LuckyTransaction(txid)
		if txnfo.Status == mesh.TxConfirmed {
			break
		}
		time.Sleep(10 * time.Second)
	}

}
