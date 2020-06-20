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
const poetRpcPort = 50002
const eventsPort = 50003
const TestNetPath = "./local-testnet/"
const BinPath = TestNetPath +"bin/"

var genesisAccounts = []string{
	/*******/  "0x097598942e44919cf7d11499887a595e41b097acd0a75d65ed8b8c6fa739d297",
	/*Almog*/  "0x4d05cfede9928bcd225c008db8110cfeb1f01011e118bdb93f1bb14d2052c276",
	/*Anton*/  "0xdb58184012f26c405bff2d8866bf7ef2d1da7f0b391d1f1364f1d695929df617",
	/*Tap*/    "0x891da146767aa80e3ce3ef826ef675c1bb32e9021844193a163fac231513149a",
	/*Yosher*/ "0x39a27e846f7e9783cd8fcae0f94abe7ba1428df096e13e903ef5b9df85d520e1",
	/*Gavrad*/ "0x0dc90fe42d96e302ae122aa3437e320d792772aba8f459f80e18a45ae754112d",
}

var nodesCount = len(genesisAccounts)-1

var globalOpts = map[string]interface{} {
	"--randcon":                      4,
	"--hare-committee-size":          5,
	"--hare-max-adversaries":         2,
	"--hare-round-duration-sec":      10,
	"--layer-duration-sec":           60,
	"--layer-average-size":           10,
	"--hare-wakeup-delta":            10,
	"--test-mode":                    nil,
	"--eligibility-confidence-param": 5,
	"--eligibility-epoch-offset":     0,
	"--genesis-active-size":          5,
	"--genesis-conf":                 TestNetPath +"genesis_accounts.json",
	"--poet-server":                  fmt.Sprintf("127.0.0.5:%v", poetPort),
	"--genesis-time": 				  time.Now().Format(time.RFC3339),
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

func Bootstrap() (term Terminator) {

	_ = os.RemoveAll(TestNetPath +".data")
	_ = os.MkdirAll(TestNetPath+".data",0777)

	sigterm := make([]chan struct{}, nodesCount+2)
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
		"--coinbase":     genesisAccounts[0],
		"--events-url": fmt.Sprintf("tcp://127.0.0.1:%v", eventsPort),
	})
	c, err := exec(BinPath+"go-spacemesh",opts,sigterm[0]) ;
	if err != nil {
		panic(err)
	}

	id := ""

	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
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
			if M, ok := m["M"].(string); ok {
				if M == "App started." {
					boot(fmt.Sprintf("spacemesh://%v@127.0.0.10:%v",id, bootstrapPort),sigterm[2:])
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

func boot(booturl string, sigterm []chan struct{}) {

	poet(sigterm[1])

	for i, t := range sigterm {
		node(i+1,booturl,t)
	}

}

func node(no int, booturl string, sigterm chan struct{}) {
	opts := mergeOpts(map[string]interface{}{
		"--json-server":  nil,
		"--json-port":    bootstratJsonPort +no,
		"--grpc-server":  nil,
		"--grpc-port":    bootstratGrpcPort +no,
		"--tcp-port":     bootstrapPort +no,
		"--start-mining": nil,
		"--bootstrap":    nil,
		"--bootnodes":    booturl+fmt.Sprintf("?disc=%v", bootstrapPort),
		"--data-folder":  fmt.Sprintf(TestNetPath+".data/data.%v",no),
		"--post-datadir": fmt.Sprintf(TestNetPath+".data/post.%v",no),
		"--coinbase":     genesisAccounts[no],
	})
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
		"--reset": nil,
	}
	c, err := exec(BinPath+"poet",opts,sigterm) ;
	if err != nil {
		panic(err)
	}
	go collect_poet(c)
	return
}

func collect(c chan string,n int) {
	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
			if M, ok := m["M"].(string); ok {
				if M == "PoST initialization completed" {
					fmt.Printf("PoST initialization completed for node %v\n",n)
					break
				}
			}
		}
	}
}

func collect_poet(c chan string) {
	for _ = range c {
		//
	}
}
