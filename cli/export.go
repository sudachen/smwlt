package cli

import (
	"fmt"
	"github.com/spacemeshos/go-spacemesh/common/util"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/wallet"
)

var cmdExport = &cobra.Command{
	Use:   "export <account>",
	Short: "Export account key pair as a hex string",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		acc := wallet.LuckyLookup(args[0], w...)
		fmt.Println(util.Bytes2Hex(acc.Private))
	},
}

