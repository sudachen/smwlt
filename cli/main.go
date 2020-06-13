package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	api "github.com/sudachen/smwlt/mesh/api.v1"
	"github.com/sudachen/smwlt/wallet"
)

const keyValueFormat = "%-8s %v\n"

var mainCmd = &cobra.Command{
	Use:              "smwlt",
	TraverseChildren: true,
}

var optWalletFile = mainCmd.Flags().StringP("wallet", "w", "", "wallet filename")
var optLegacy = mainCmd.Flags().BoolP("legacy", "l", false, "use legacy unencrypted file format")
var optPassword = mainCmd.Flags().StringP("password", "p", "", "wallet unlock password")
var optEndpoint = mainCmd.Flags().StringP("endpoint", "e", api.DefaultEndpoint, "host:port to connect mesh node")
var optYes = mainCmd.Flags().BoolP("yes", "y", false, "auto confirm")
var optVerbose = mainCmd.Flags().BoolP("verbose", "v", false, "be verbose")
var OptTrace = mainCmd.Flags().BoolP("trace", "x", false, "backtrace on panic")

func Main() {
	mainCmd.AddCommand(
		cmdInfo,
		cmdSend,
		cmdTxs,
		cmdNet,
		cmdHexSign,
		cmdTextSign,
		cmdCoinbase,
	)

	if err := mainCmd.Execute(); err != nil {
		panic(fu.Panic(err, 1))
	}
}

func loadWallet() (w []wallet.Wallet) {
	if *optLegacy {
		w = []wallet.Wallet{wallet.Legacy{Path: *optWalletFile}.LuckyLoad()}
		// unencrypted
	} else {
		w = []wallet.Wallet{}
		for _, x := range w {
			if *optPassword != "" {
				err := x.Unlock(*optPassword)
				if err == nil {

				}
			}
		}
		if len(w) == 0 && *optPassword != "" {
			panic(fu.Panic(fmt.Errorf("there is nothing to unlock, wrong password(?)")))
		}
		panic(fu.Panic(fmt.Errorf("unsupported wallet type")))
	}
	return
}

type Client struct {
	*api.ClientAgent
	DefaultGasLimit uint64
	DefaultFee      uint64
}

func newClient() Client {
	return Client{
		ClientAgent:     api.Client{Endpoint: *optEndpoint, Verbose: *optVerbose}.New(),
		DefaultGasLimit: api.DefaultGasLimit,
		DefaultFee:      api.DefaultFee,
	}
}
