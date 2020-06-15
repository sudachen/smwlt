package zinc

import (
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/node/api.v1"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/legacy"
	"math/rand"
	"testing"
	"time"
)

func Test_Stress1(t *testing.T) {
	ca := make([]*api_v1.ClientAgent, 5)
	ca[0] = api_v1.Client{}.New()
	for i := 1; i < 5; i++ {
		ca[i] = api_v1.Client{Endpoint: fmt.Sprintf("localhost:919%d", i)}.New()
	}
	c := ca[0]
	w := legacy.Wallet{Path: "../accounts.json"}.LuckyLoad()
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
			txid, err = c.Transfer(1, from.Address, nonce[j], from.Private, to.Address, api_v1.DefaultFee, api_v1.DefaultGasLimit)
			if err == nil {
				nonce[j]++
			} else {
				nonce[j] = c.LuckyAccountInfo(from.Address).Nonce
			}
		}
		time.Sleep(3 * time.Second)
	}
	for {
		txnfo := ca[0].LuckyTransactionInfo(txid)
		if txnfo.Status == api_v1.TxConfirmed {
			break
		}
		time.Sleep(10 * time.Second)
	}

}
