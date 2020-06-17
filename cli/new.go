package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
)

var cmdNew = &cobra.Command{
	Use:   "new <account>",
	Short: "Create a new account (key pair)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadOrCreateWallet()
		fmt.Printf("Selected wallet%v:\n", fu.Ifs(len(w) > 1, "s", ""))
		for _, x := range w {
			fmt.Println("\t" + x.DisplayName())
		}
		if len(w) > 1 {
			panic(fu.Panic(errors.New("wallet is ambiguous, you must select only one")))
		}
		a, exists := w[0].Lookup(args[0])
		if exists {
			panic(fu.Panic(fmt.Errorf("Account '%v' already exists", a.Name)))
		}
		w[0].LuckyNewPair(args[0])
		w[0].LuckySave()
	},
}
