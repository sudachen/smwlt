package cli

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu"
	api "github.com/sudachen/smwlt/node/api.v1"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/legacy"
	"github.com/sudachen/smwlt/wallet/modern"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const MajorVersion = 1
const MinorVersion = 0
const keyValueFormat = "%-8s %v\n"

var mainCmd = &cobra.Command{
	Use:           "smwlt",
	Short:         fmt.Sprintf("Spacemesh CLI Wallet %v.%v (https://github.com/sudachen/smwlt)", MajorVersion, MinorVersion),
	SilenceErrors: true,
}

var optWalletFile = mainCmd.PersistentFlags().StringP("wallet-file", "f", "", "use wallet filename")
var optWalletName = mainCmd.PersistentFlags().StringP("wallet-name", "n", "", "select wallet by name")
var optWalletDir = mainCmd.PersistentFlags().StringP("wallet-dir", "d", modern.DefaultDirectory(), "use wallet dir")
var optLegacy = mainCmd.PersistentFlags().BoolP("legacy", "l", false, "use legacy unencrypted file format")
var optPassword = mainCmd.PersistentFlags().StringP("password", "p", "", "wallet unlock password")
var optEndpoint = mainCmd.PersistentFlags().StringP("endpoint", "e", api.DefaultEndpoint, "host:port to connect mesh node")
var optYes = mainCmd.PersistentFlags().BoolP("yes", "y", false, "auto confirm")
var OptTrace = mainCmd.PersistentFlags().BoolP("trace", "x", false, "backtrace on panic")

func init() {
	mainCmd.PersistentFlags().BoolP("help", "h", false, "help for info")
	fu.VerboseOptP = mainCmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")
	mainCmd.AddCommand(
		cmdInfo,
		cmdSend,
		cmdTxs,
		cmdNet,
		cmdHexSign,
		cmdTextSign,
		cmdCoinbase,
		cmdNew,
		cmdExport,
		cmdImport,
	)
}

func CLI() *cobra.Command {
	return mainCmd
}

func Main() {

	cst, _ := terminal.GetState(int(os.Stdin.Fd()))
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		terminal.Restore(int(os.Stdin.Fd()), cst)
		os.Exit(1)
	}()
	if err := CLI().Execute(); err != nil {
		panic(fu.Panic(err, 1))
	}
}

func unlock(w wallet.Wallet, passw *[]string, interactive bool) bool {
	for _, p := range *passw {
		if e := w.Unlock(p); e == nil {
			return true
		}
	}
	if interactive {
		fmt.Printf("Unlocking wallet %v\n", w.DisplayName())
		p := prompter.Password("Enter password [leave empty to skip]")
		if p != "" {
			if e := w.Unlock(p); e == nil {
				fmt.Println("Wallet unlocked")
				*passw = append(*passw, p)
				return true
			} else {
				fmt.Println("Wrong password!")
			}
		} else {
			fmt.Println("Wallet skipped")
		}
	}
	return false
}

func loadWallet(canBeEmpty ...bool) (w []wallet.Wallet) {
	if *optLegacy {
		w = []wallet.Wallet{legacy.Wallet{Path: *optWalletFile}.LuckyLoad()}
		// unencrypted
	} else {
		w = []wallet.Wallet{}
		wx := []wallet.Wallet{}
		passw := []string{}
		if *optPassword != "" {
			passw = append(passw, *optPassword)
		}
		if *optWalletFile != "" {
			wx = []wallet.Wallet{modern.Wallet{Path: *optWalletFile}.LuckyLoad()}
		} else {
			if err := filepath.Walk(*optWalletDir, func(path string, info os.FileInfo, err error) error {
				base := filepath.Base(path)
				if strings.HasPrefix(base, "my_wallet_") && strings.HasSuffix(base, ".json") {
					fu.Verbose("opening wallet file '%v'", base)
					wal, err := modern.Wallet{Path: path}.Load()
					if err == nil {
						if *optWalletName == "" ||
							strings.HasPrefix(strings.ToLower(wal.DisplayName()), strings.ToLower(*optWalletName)) {
							wx = append(wx, wal)
						}
					} else {
						fu.Verbose("failed to open with error: %v", err.Error())
					}
				}
				return nil
			}); err != nil {
				panic(fu.Panic(err))
			}
		}
		for _, x := range wx {
			if unlock(x, &passw, *optPassword == "") {
				w = append(w, x)
			}
		}
		if len(w) == 0 && *optPassword != "" && !fu.Fnf(canBeEmpty...) {
			panic(fu.Panic(fmt.Errorf("there is nothing to unlock, wrong password(?)")))
		}
	}
	return
}

func exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadOrCreateWallet() (w []wallet.Wallet) {
	if *optWalletFile == "" || exist(*optWalletFile) {
		w = loadWallet(true)
	}
	if len(w) == 0 {
		if *optLegacy {
			w = []wallet.Wallet{legacy.Wallet{*optWalletFile}.New()}
		} else {
			if *optWalletName == "" {
				panic(fu.Panic(fmt.Errorf("you must specify new wallet name")))
			}
			p := *optPassword
			for p == "" {
				p = prompter.Password("Enter new wallet password")
				if p != "" {
					if p == prompter.Password("Verify new wallet password") {
						break
					}
					fmt.Println("Does not match")
				}
			}
			path := *optWalletFile
			if path == "" {
				for {
					path = fmt.Sprintf("my_wallet_%s.json", time.Now().UTC().Format("2006-01-02T15-04-05.000")+"Z")
					if *optWalletDir != "" {
						path = filepath.Join(*optWalletDir, path)
					}
					if !exist(path) { // why not?
						break
					}
				}
			}
			mnemonic, err := wallet.NewMnemonic()
			if err != nil {
				panic(fu.Panic(fu.Wrapf(err, "failed to create new mnemonic: %v", err.Error())))
			}
			w = []wallet.Wallet{modern.Wallet{path, *optWalletName}.New(p, mnemonic)}
			fmt.Print("New wallet mnemonic:")
			for i, x := range strings.Split(mnemonic, " ") {
				if i%4 == 0 {
					fmt.Print("\n\t")
				}
				fmt.Printf("%-20s", x)
			}
			fmt.Println("")
		}
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
		ClientAgent:     api.Client{Endpoint: *optEndpoint, Verbose: fu.Verbose}.New(),
		DefaultGasLimit: api.DefaultGasLimit,
		DefaultFee:      api.DefaultFee,
	}
}
