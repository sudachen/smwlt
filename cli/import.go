package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/sudachen/smwlt/fu/prompter"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/stdio"
	"strings"
)

var cmdImport = &cobra.Command{
	Use:   "import <account> <hex-string>",
	Short: "Import account key pair from the hex string",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadOrCreateWallet()
		stdio.Printf("Selected wallet%v:\n", fu.Ifs(len(w) > 1, "s", ""))
		for _, x := range w {
			stdio.Println("\t" + x.DisplayName())
		}
		if len(w) > 1 {
			panic(errstr.New(1, "wallet is ambiguous, you must select only one"))
		}
		a, exists := w[0].Lookup(args[0])
		if exists {
			panic(errstr.Format(1, "Account '%v' already exists", a.Name))
		}
		bs, err := hex.DecodeString(args[1])
		if err != nil {
		}
		key := ed25519.PrivateKey(bs)
		address := types.BytesToAddress(key.Public().(ed25519.PublicKey)[:])
		stdio.Printf("Account %v [for import]:\n"+strings.Repeat("\t"+keyValueFormat, 1),
			args[0],
			"Address:", address.Hex(),
		)
		yes := *optYes
		if !yes {
			yes = prompter.YN(fmt.Sprintf("Import '%v' to the wallet", args[0]), false)
		}
		if yes {
			w[0].LuckyImportKey(args[0], address, key)
			w[0].LuckySave()
			stdio.Println("Successfully imported")
		} else {
			stdio.Println("Cancelled")
		}
	},
}
