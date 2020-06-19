package cli

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/prompter"
	"github.com/sudachen/smwlt/fu/stdio"
	"github.com/sudachen/smwlt/wallet"
	"strconv"
	"strings"
)

var cmdSend = &cobra.Command{
	Use:   "send <from> <to> <amount> [fee]",
	Short: "Transfer coins from one to another account",
	Args:  cobra.RangeArgs(3, 4),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		from := wallet.LuckyLookup(args[0], w...)
		var to types.Address
		toa, exists := wallet.LookupOne(args[1], w...)
		if !exists {
			x, err := types.StringToAddress(args[1])
			if err != nil {
				panic(errstr.Format(0, "account '%v' does not exist", args[1]))
			}
			to = x
		} else {
			to = toa.Address
		}
		c := newClient()
		nfo := c.LuckyAccountInfo(from.Address)
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			panic(errstr.Format(0, "bad amount '%v'", args[2]))
		}
		fee := int(c.DefaultFee)
		if len(args) > 3 {
			fee, err = strconv.Atoi(args[3])
			if err != nil {
				panic(errstr.Format(0, "bad fee '%v'", args[3]))
			}
		}
		if nfo.Balance < uint64(amount+fee) {
			panic(errstr.Format(0, "not enough balance"))
		}
		stdio.Printf("Transfer coins:\n"+strings.Repeat("\t"+keyValueFormat, 5),
			"From:", from.Address.Hex(),
			"To:", to.Hex(),
			"Balance:", nfo.Balance,
			"Amount:", amount,
			"Fee:", fee)
		ok := *optYes
		if !ok {
			ok = prompter.YN("Confirm transaction", false)
		}
		if !ok {
			stdio.Println("Cancelled")
			return
		}
		txid := c.LuckyTransfer(uint64(amount), from.Address, nfo.Nonce, from.Private, to, uint64(fee), c.DefaultGasLimit)
		stdio.Println("Succeeded with TxID:", txid.String())
	},
}
