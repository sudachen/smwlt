package cli

import (
	"github.com/spf13/cobra"
)

var cmdNew = &cobra.Command{
	Use:   "new <account>",
	Short: "Create a new account (key pair)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
