package mesh

import (
	"github.com/sudachen/smwlt/fu"
)

/*
MiningStatus is enumeration of smesher mining statuses
*/
type MiningStatus int

const (
	MiningUnknown MiningStatus = iota
	MiningIdel
	MiningInProgress
	MiningDone
)

/*
NodeStatus describes mesh node
*/
type NodeStatus struct {
	Synced        bool   `json:"synced"`
	SyncedLayer   uint64 `json:"syncedLayer,string"`
	CurrentLayer  uint64 `json:"currentLayer,string"`
	VerifiedLayer uint64 `json:"verifiedLayer,string"`
	Peers         uint64 `json:"peers,string"`
	MinPeers      uint64 `json:"minPeers,string"`
	MaxPeers      uint64 `json:"maxPeers,string"`
	DataDir       string `json:"dataDir"`
}

/*
MiningStats describes mining
*/
type MiningStats struct {
	Status                 MiningStatus `json:"status"`
	Coinbase               string       `json:"coinbase"`
	SmeshingRemainingBytes int          `json:"remainingBytes,string"`
}

/*
NodeInfo is an integral node information
*/
type NodeInfo struct {
	NodeStatus
	MiningStats
}

/*
GetNodeInfo calls to node for the integral node information
*/
func (c *ClinetAgent) GetNodeInfo() (info NodeInfo, err error) {
	if info.NodeStatus, err = c.GetNodeStatus(); err != nil {
		return
	}
	info.MiningStats, err = c.GetMiningStats()
	return
}

/*
LuckyNodeInfo calls to node for the integral node information. It panics on error
*/
func (c *ClinetAgent) LuckyNodeInfo() (info NodeInfo) {
	fu.LuckyCall(c.GetNodeInfo, &info)
	return
}

/*
GetNodeStatus calls to node for the node status
*/
func (c *ClinetAgent) GetNodeStatus() (st NodeStatus, err error) {
	err = c.post("/nodestatus", nil, &st)
	return
}

/*
LuckyNodeStatus calls to node for the node status. It panics on error
*/
func (c *ClinetAgent) LuckyNodeStatus() (st NodeStatus) {
	fu.LuckyCall(c.GetNodeStatus, &st)
	return
}

/*
GetMiningStats calls to node for the mining stats
*/
func (c *ClinetAgent) GetMiningStats() (st MiningStats, err error) {
	err = c.post("/stats", nil, &st)
	return
}

/*
LuckyMiningStats calls to node for the mining stats. It panics on error
*/
func (c *ClinetAgent) LuckyMiningStatus() (st MiningStats) {
	fu.LuckyCall(c.GetMiningStats, &st)
	return
}
