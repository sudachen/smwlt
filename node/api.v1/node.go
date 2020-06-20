package api_v1

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/sudachen/smwlt/fu"
)

/*
MiningStatus is enumeration of smesher mining statuses
*/
type MiningStatus int

const (
	MiningUnknown MiningStatus = iota
	MiningIdle
	MiningInProgress
	MiningDone
)

func (st MiningStatus) String() string {
	switch st {
	case MiningIdle:
		return "Idle"
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
	//VerifiedLayer uint64 `json:"verifiedLayer,string"`
	Peers         uint64 `json:"peers,string"`
	MinPeers      uint64 `json:"minPeers,string"`
	MaxPeers      uint64 `json:"maxPeers,string"`
}

/*
MiningStats describes mining
*/
type MiningStats struct {
	Status                 MiningStatus  `json:"status"`
	Coinbase               types.Address `json:"-"`
	//SmeshingRemainingBytes int           `json:"remainingBytes,string"`
	DataDir                string        `json:"dataDir"`
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
	out := struct {
		*MiningStats
		StrCoinbase string `json:"coinbase"`
	}{MiningStats: &st}
	err = c.post("/stats", nil, &out)
	if err == nil {
		st.Coinbase, err = types.StringToAddress(out.StrCoinbase)
	}
	return
}

/*
LuckyMiningStats calls to node for the mining stats. It panics on error
*/
func (c *ClientAgent) LuckyMiningStats() (st MiningStats) {
	fu.LuckyCall(c.GetMiningStats, &st)
	return
}

/*
SetCoinbase sets coinbase address
*/
func (c *ClientAgent) SetCoinbase(address types.Address) (err error) {
	in := struct {
		Address string `json:"address"`
	}{address.Hex()}
	err = c.post("/setawardsaddr", &in, &struct{}{})
	return
}

/*
LuckyCoinbase sets coinbase address. It panics on error
*/
func (c *ClientAgent) LuckyCoinbase(address types.Address) {
	fu.LuckyCall(c.SetCoinbase, nil, address)
}
