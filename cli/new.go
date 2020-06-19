package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/prompter"
	"github.com/sudachen/smwlt/fu/stdio"
	"strings"
)

var cmdNew = &cobra.Command{
	Use:   "new <account>",
	Short: "Create a new account (key pair)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadOrCreateWallet()
		stdio.Printf("Selected wallet%v:\n", fu.Ifs(len(w) > 1, "s", ""))
		for _, x := range w {
			stdio.Println("\t" + x.DisplayName())
		}
		if len(w) > 1 {
			panic(errstr.New(0, "wallet is ambiguous, you must select only one"))
		}
		a, exists := w[0].Lookup(args[0])
		if exists {
			panic(errstr.Format(0, "Account '%v' already exists", a.Name))
		}
		yes := *optYes
		if !yes {
			yes = prompter.YN(fmt.Sprintf("Create account '%v' in this wallet", args[0]), false)
		}
		if yes {
			w[0].LuckyNewPair(args[0])
			w[0].LuckySave()
			a, _ = w[0].Lookup(args[0])
			stdio.Printf("Account %v [%v]:\n"+strings.Repeat("\t"+keyValueFormat, 1),
				args[0], w[0].DisplayName(),
				"Address:", a.Address.Hex(),
			)
			stdio.Println("Successfully created")
		} else {
			stdio.Println("Cancelled")
		}
	},
}
