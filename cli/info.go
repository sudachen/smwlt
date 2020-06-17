package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/wallet"
	"sort"
	"strings"
)

var cmdInfo = &cobra.Command{
	Use:   "info [account]...",
	Short: "Display accounts info (address, balance, and nonce)",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		c := newClient()
		if len(args) > 0 {
			for _, a := range args {
				accs := wallet.Lookup(a, w...)
				if len(accs) == 0 {
					fmt.Printf("Account %v:\n\tdoes not exist\n", a)
				}
				for _, a := range accs {
					nfo := c.LuckyAccountInfo(a.Address)
					fmt.Printf("Account %v [%v]:\n"+strings.Repeat("\t"+keyValueFormat, 3),
						a.Name,
						a.Wallet.DisplayName(),
						"Address:", a.Address.Hex(),
						"Balance:", nfo.Balance,
						"Nonce:", nfo.Nonce,
					)
				}
			}
		} else {
			m := map[string][]wallet.Account{}
			for _, wal := range w {
				accs := wal.List()
				for _, a := range accs {
					m[a.Name] = append(m[a.Name], a)
				}
			}
			names := make([]string, 0, len(m))
			for k := range m {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, n := range names {
				for _, a := range m[n] {
					nfo, err := c.GetAccountInfo(a.Address)
					if err != nil {
						fmt.Printf("Account %v [%v]:\n\t%v\n",
							a.Name,
							a.Wallet.DisplayName(),
							err.Error(),
						)
					} else {
						fmt.Printf("Account %v [%v]:\n"+strings.Repeat("\t"+keyValueFormat, 3),
							a.Name,
							a.Wallet.DisplayName(),
							"Address:", a.Address.Hex(),
							"Balance:", nfo.Balance,
							"Nonce:", nfo.Nonce,
						)
					}
				}
			}
		}
	},
}
