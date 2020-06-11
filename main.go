package main

import (
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/mesh"
	"github.com/sudachen/smwlt/wallet"
	"strconv"
	"github.com/Songmu/prompter"
)

func loadWallet(path string, legacy bool) wallet.Wallet {
	if legacy {
		return wallet.Legacy{Path:path}.LuckyLoad()
	} else {
		panic(fu.Panic(fmt.Errorf("unsupported wallet type")))
	}
}

func main() {

	rootCmd := &cobra.Command{Use:"smwlt", TraverseChildren: true}

	walletFile := rootCmd.Flags().StringP("wallet", "w", "", "wallet filename")
	legacy := rootCmd.Flags().BoolP("legacy", "l", false, "use legacy unencrypted file format")
	password := rootCmd.Flags().StringP("password", "p", "", "wallet unlock password")
	endpoint := rootCmd.Flags().StringP("endpoint", "e", mesh.DefaultEndpoint, "host:port to connect mesh node")
	yes := rootCmd.Flags().BoolP("yes", "y", false, "auto confirm")

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "info <account>",
			Short: "get account info",
			Args: cobra.RangeArgs(1,1),
			Run: func(cmd *cobra.Command, args []string) {
				w := loadWallet(*walletFile,*legacy)
				if *password != "" { w.LuckyUnlock(*password) }
				acc, exists := w.Lookup(args[0])
				if !exists { fu.Panic(fmt.Errorf("account '%v' does not exist", args[0])) }
				c := mesh.Client{Endpoint:*endpoint}.New()
				nfo := c.LuckyAccountInfo(acc.Address)
				fmt.Printf("Address: %v\nCreated: %v\nNonce: %v\nBalance: %v\n",
					acc.Address.Hex(),
					acc.Created,
					nfo.Nonce,
					nfo.Balance)
			},
		},
		&cobra.Command{
			Use:   "tx <from> <to> <amount> [fee]",
			Short: "do transfer",
			Args: cobra.RangeArgs(3,4),
			Run: func(cmd *cobra.Command, args []string) {
				w := loadWallet(*walletFile,*legacy)
				if *password != "" { w.LuckyUnlock(*password) }
				from, exists := w.Lookup(args[0])
				if !exists { fu.Panic(fmt.Errorf("account '%v' does not exist", args[0])) }
				to, exists := w.Lookup(args[1])
				if !exists { fu.Panic(fmt.Errorf("account '%v' does not exist", args[1])) }
				c := mesh.Client{Endpoint:*endpoint}.New()
				nfo := c.LuckyAccountInfo(from.Address)
				amount, err := strconv.Atoi(args[2])
				if err != nil {
					panic(fu.Panic(fmt.Errorf("bad amount '%v'",args[2])))
				}
				fee := int(mesh.DefaultFee)
				if len(args) > 3 {
					fee, err = strconv.Atoi(args[3])
					if err != nil {
						panic(fu.Panic(fmt.Errorf("bad fee '%v'",args[3])))
					}
				}
				if nfo.Balance < uint64(amount+fee) {
					panic(fu.Panic(fmt.Errorf("not enough balance")))
				}
				fmt.Printf("From:    %v\nBalance: %d\nTo:      %v\nAmount:  %d\nFee:     %v\n",
					from.Address.Hex(),
					nfo.Balance,
					to.Address.Hex(),
					amount,
					fee)
				ok := *yes
				if !ok {
					ok = prompter.YN("Confirm transaction",false)
				}
				if !ok {
					fmt.Println("cancelled")
					return
				}
				txid := c.LuckyTransfer(uint64(amount),from.Address,nfo.Nonce,from.Private,to.Address,uint64(fee),mesh.DefaultGasLimit)
				fmt.Println("Succeeded with TxID:",txid.String())
			},
		},
		&cobra.Command{
			Use:   "txs <account> [startLayer]",
			Short: "do transfer",
			Args:  cobra.RangeArgs(1, 2),
			Run: func(cmd *cobra.Command, args []string) {
				w := loadWallet(*walletFile,*legacy)
				if *password != "" { w.LuckyUnlock(*password) }
				acc, exists := w.Lookup(args[0])
				if !exists { fu.Panic(fmt.Errorf("account '%v' does not exist", args[0])) }
				c := mesh.Client{Endpoint:*endpoint}.New()
				layer := uint64(0)
				if len(args) > 1 {
					x, err := strconv.Atoi(args[1])
					if err != nil {
						panic(fu.Panic(fmt.Errorf("bad layer '%v'",args[1])))
					}
					layer = uint64(x)
				}

				mark := func(address types.Address) string {
					if address == acc.Address {
						return "*"
					}
					return ""
				}

				txs := c.LuckyTxList(acc.Address,layer)
				for i,x := range txs {
					tx := c.LuckyTransaction(x)
					fmt.Printf("%3d:\t%-8s %v\n\t%-8s %v\n\t%-8s %v\n\t%-8s %d\n\t%-8s %d\n\t%-8s %s\n",
						i,
						"Id:",
						tx.Id.String(),
						"From"+mark(tx.From)+":",
						tx.From.String(),
						"To"+mark(tx.To)+":",
						tx.To.String(),
						"Amount:",
						tx.Amount,
						"Layer:",
						tx.LayerId,
						"Status:",
						tx.Status)
				}
			},
		},
	)

	if err := rootCmd.Execute(); err != nil {
		panic(fu.Panic(err,1))
	}
}
