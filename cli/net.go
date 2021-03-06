package cli

import (
	"github.com/spf13/cobra"
	"github.com/sudachen/smwlt/fu/stdio"
	"strings"
)

var cmdNet = &cobra.Command{
	Use:   "net",
	Short: "Display the node status",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		nfo := c.LuckyNodeInfo()
		const format = "%-16s %v\n"
		stdio.Printf("Node status:\n"+strings.Repeat("\t"+format, 9),
			"Synced:", nfo.Synced,
			"Synced layer:", nfo.SyncedLayer,
			"Current layer:", nfo.CurrentLayer,
			//"Verified layer:", nfo.VerifiedLayer,
			"Peers:", nfo.Peers,
			"Min peers:", nfo.MinPeers,
			"Max peers:", nfo.MaxPeers,
			"Data directory:", nfo.DataDir,
			"Mining status:", nfo.Status,
			"Coinbase:", nfo.Coinbase.Hex(),
			//"Remaining bytes:", nfo.SmeshingRemainingBytes,
		)
	},
}
