package go_testnet

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Terminator []chan struct{}
func (t Terminator) Sigterm() {
	for _,x := range t {
		close(x)
	}
}

const bootstrapPort = 17513
const bootstratJsonPort = 19090
const bootstratGrpcPort = 19190
const poetPort = 10080
const poetRpcPort = 10082
const eventsPort = 10083
const TestNetPath = "./local-testnet/"
const BinPath = TestNetPath +"bin/"

const genesisCoinbase = "0x0000000000000000000000000000000000000000000000000000000000000001"

var genesisAccounts = []string{
	/*Almog*/  "0x4d05cfede9928bcd225c008db8110cfeb1f01011e118bdb93f1bb14d2052c276",
	/*Anton*/  "0xdb58184012f26c405bff2d8866bf7ef2d1da7f0b391d1f1364f1d695929df617",
	/*Tap*/    "0x891da146767aa80e3ce3ef826ef675c1bb32e9021844193a163fac231513149a",
	/*Yosher*/ "0x39a27e846f7e9783cd8fcae0f94abe7ba1428df096e13e903ef5b9df85d520e1",
	/*Gavrad*/ "0x0dc90fe42d96e302ae122aa3437e320d792772aba8f459f80e18a45ae754112d",
}

var globalOpts = map[string]interface{} {
	"--acquire-port":				  "false",
	"--randcon":                      3,
	"--hare-committee-size":          1,
	"--hare-max-adversaries":         1,
	"--hare-round-duration-sec":      10,
	"--layer-duration-sec":           60,
	"--layer-average-size":           10,
	"--hare-wakeup-delta":            10,
	"--test-mode":                    nil,
	"--eligibility-confidence-param": 5,
	"--eligibility-epoch-offset":     0,
	"--genesis-active-size":          5,
	"--genesis-conf":                 TestNetPath +"genesis_accounts.json",
	"--genesis-time": 				  time.Now().Round(time.Second).Format("2006-01-02T15:04:05-07:00"),
	"--poet-server": 				  fmt.Sprintf("127.0.0.1:%v", poetPort),
}

func mergeOpts(addopt map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for k,v := range globalOpts {
		r[k] = v
	}
	for k,v := range addopt {
		r[k] = v
	}
	return r
}

func Bootstrap(nodes int, miners int) (term Terminator) {
	if miners > nodes { miners = nodes }

	_ = os.RemoveAll(TestNetPath +".data")
	_ = os.MkdirAll(TestNetPath+".data",0777)

	sigterm := make([]chan struct{}, nodes+2)
	for i := range sigterm { sigterm[i] = make(chan struct{}) }
	defer func() { for _,x := range sigterm { close(x) } }()

	opts := mergeOpts(map[string]interface{}{
		"--json-server":  nil,
		"--json-port":    bootstratJsonPort,
		"--grpc-server":  nil,
		"--grpc-port":    bootstratGrpcPort,
		"--tcp-port":     bootstrapPort,
		"--data-folder":  fmt.Sprintf(TestNetPath+".data/data.%v",0),
		"--post-datadir": fmt.Sprintf(TestNetPath+".data/post.%v",0),
		//"--events-url":   fmt.Sprintf("tcp://127.0.0.1:%v", eventsPort),
		//"--coinbase":     genesisCoinbase,
		//"--start-mining": nil,
	})
	c, err := exec(BinPath+"go-spacemesh",opts,sigterm[0]) ;
	if err != nil {
		panic(err)
	}

	id := ""

	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
			parse(m, 0, "BOOTSTRAP:")
			if M, ok := m["M"].(string); ok {
				const prefix = "Local node identity >> "
				if strings.HasPrefix(M,prefix) {
					id = M[len(prefix):]
					break
				}
			}
		}
	}

	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
			parse(m, 0, "BOOTSTRAP:")
			if M, ok := m["M"].(string); ok {
				if M == "App started." {
					boot(fmt.Sprintf("spacemesh://%v@127.0.0.1:%v",id, bootstrapPort), sigterm[1:], miners)
					break
				}
			}
		}
	}

	go collect(c,0)
	term = sigterm
	sigterm = nil
	return
}

func boot(booturl string, sigterm []chan struct{}, miners int) {
	poet(sigterm[0])
	for i, t := range sigterm[1:] {
		node(i+1,booturl,t,i<miners)
	}
	for ( len(p2pid) < len(sigterm)-1 ) {
		time.Sleep(time.Second)
	}
}

func node(no int, booturl string, sigterm chan struct{}, miner bool) {
	opts := mergeOpts(map[string]interface{}{
		"--json-server":  nil,
		"--json-port":    bootstratJsonPort +no,
		"--grpc-server":  nil,
		"--grpc-port":    bootstratGrpcPort +no,
		"--tcp-port":     bootstrapPort +no,
		"--bootstrap":    nil,
		"--bootnodes":    booturl+fmt.Sprintf("?disc=%v", bootstrapPort),
		"--data-folder":  fmt.Sprintf(TestNetPath+".data/data.%v",no),
		"--post-datadir": fmt.Sprintf(TestNetPath+".data/post.%v",no),
	})
	if miner {
		opts["--start-mining"] = nil
		opts["--coinbase"] = genesisAccounts[no%len(genesisAccounts)]
	}
	c, err := exec(BinPath+"go-spacemesh",opts,sigterm) ;
	if err != nil {
		panic(err)
	}
	go collect(c,no)
	return
}

func poet(sigterm chan struct{}) {
	opts := map[string]interface{} {
		"--rpclisten":fmt.Sprintf("127.0.0.1:%v", poetRpcPort),
		"--restlisten":fmt.Sprintf("127.0.0.1:%v", poetPort),
		"--initialduration":"10s",
		"--duration":"10s",
		"--gateway":fmt.Sprintf("127.0.0.1:%v", bootstratGrpcPort),
		"--poetdir":fmt.Sprintf(TestNetPath+".data/%v","poet"),
		"--empty":nil,
		"--disablebroadcast":nil,
	}
	c, err := exec(BinPath+"poet",opts,sigterm) ;
	if err != nil {
		panic(err)
	}
	go collect_poet(c)
	return
}
