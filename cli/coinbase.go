package cli

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/wallet"
)

var cmdCoinbase = &cobra.Command{
	Use:   "coinbase <account>",
	Short: "Set the account as coinbase account in the node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		acc := wallet.LuckyLookup(args[0], w...)
		c := newClient()
		nfo := c.LuckyMiningStats()
		if nfo.Coinbase == acc.Address {
			fmt.Printf("Node coinbase already is %v\n", acc.Address.Hex())
			return
		}
		ok := *optYes
		if !ok {
			ok = prompter.YN(fmt.Sprintf("Set node coinbase to %v", acc.Address.Hex()), false)
		}
		if !ok {
			fmt.Println("Cancelled")
			return
		}
		c.LuckyCoinbase(acc.Address)
		nfo = c.LuckyMiningStats()
		if nfo.Coinbase != acc.Address {
			panic(fu.Panic(fmt.Errorf("oops, coinbase is not updated"), 2))
		}
		fmt.Printf("Succeeded, codebase is %v now\n", nfo.Coinbase.Hex())
	},
}
