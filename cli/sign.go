package cli

import (
	"encoding/hex"
	"github.com/spacemeshos/ed25519"
	"github.com/spacemeshos/go-spacemesh/common/util"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/stdio"
	"github.com/sudachen/smwlt/wallet"
	"strings"
)

var cmdHexSign = &cobra.Command{
	Use:   "signhex <account> <hex-string>",
	Short: "Sign a hex message with the account private key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		acc := wallet.LuckyLookup(args[0], w...)
		a := strings.TrimPrefix(args[1], "0x")
		msg, err := hex.DecodeString(a)
		if err != nil {
			panic(errstr.Format(0, "failed to decode msg hex string: %v", err.Error()))
		}
		signature := ed25519.Sign2(acc.Private, msg)
		stdio.Println(util.Bytes2Hex(signature))
	},
}

var cmdTextSign = &cobra.Command{
	Use:   "signtext <account> <utf8-string>",
	Short: "Sign a text message with the account private key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		w := loadWallet()
		acc := wallet.LuckyLookup(args[0], w...)
		signature := ed25519.Sign2(acc.Private, []byte(args[1]))
		stdio.Println(util.Bytes2Hex(signature))
	},
}
