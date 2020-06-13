package api_v1

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

func (st MiningStatus) String() string {
	switch st {
	case MiningIdel:
		return "Idel"
	case MiningInProgress:
		return "In Progress"
	case MiningDone:
		return "Done"
	}
	return "Unknown"
}

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
}

/*
MiningStats describes mining
*/
type MiningStats struct {
	Status                 MiningStatus `json:"status"`
	Coinbase               string       `json:"coinbase"`
	SmeshingRemainingBytes int          `json:"remainingBytes,string"`
	DataDir                string       `json:"dataDir"`
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
func (c *ClientAgent) GetNodeInfo() (info NodeInfo, err error) {
	if info.NodeStatus, err = c.GetNodeStatus(); err != nil {
		return
	}
	info.MiningStats, err = c.GetMiningStats()
	return
}

/*
LuckyNodeInfo calls to node for the integral node information. It panics on error
*/
func (c *ClientAgent) LuckyNodeInfo() (info NodeInfo) {
	fu.LuckyCall(c.GetNodeInfo, &info)
	return
}

/*
GetNodeStatus calls to node for the node status
*/
func (c *ClientAgent) GetNodeStatus() (st NodeStatus, err error) {
	err = c.post("/nodestatus", nil, &st)
	return
}

/*
LuckyNodeStatus calls to node for the node status. It panics on error
*/
func (c *ClientAgent) LuckyNodeStatus() (st NodeStatus) {
	fu.LuckyCall(c.GetNodeStatus, &st)
	return
}

/*
GetMiningStats calls to node for the mining stats
*/
func (c *ClientAgent) GetMiningStats() (st MiningStats, err error) {
	err = c.post("/stats", nil, &st)
	return
}

/*
LuckyMiningStats calls to node for the mining stats. It panics on error
*/
func (c *ClientAgent) LuckyMiningStatus() (st MiningStats) {
	fu.LuckyCall(c.GetMiningStats, &st)
	return
}
