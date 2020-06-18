package cli

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/stdio"
	"github.com/sudachen/smwlt/wallet"
	"strconv"
	"strings"
)

var cmdTxs = &cobra.Command{
	Use:   "txs <account> [startLayer]",
	Short: "List transactions (outgoing and incoming) for the account",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		acc := wallet.LuckyLookup(args[0], w...)
		c := newClient()
		layer := uint64(0)
		if len(args) > 1 {
			x, err := strconv.Atoi(args[1])
			if err != nil {
				panic(errstr.Format(0, "bad layer '%v'", args[1]))
			}
			layer = uint64(x)
		}

		mark := func(address types.Address) string {
			if address == acc.Address {
				return "*"
			}
			return ""
		}

		txs := c.LuckyTxList(acc.Address, layer)
		stdio.Println("TX list for " + acc.Address.Hex() + ":")
		for i, x := range txs {
			tx := c.LuckyTransactionInfo(x)
			stdio.Printf("%3d:"+strings.Repeat("\t"+keyValueFormat, 7),
				i,
				"Id:", tx.Id.String(),
				"From"+mark(tx.From)+":", tx.From.String(),
				"To"+mark(tx.To)+":", tx.To.String(),
				"Amount:", tx.Amount,
				"Fee:", tx.Fee,
				"Layer:", tx.LayerId,
				"Status:", tx.Status)
		}
	},
}
