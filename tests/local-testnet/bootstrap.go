package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Terminator []chan struct{}
func (t Terminator) Sigterm() {
	for _,x := range t {
		close(x)
	}
}

const nodesCount = 2
const bootstrapPort = 17513
const bootstratJsonPort = 19090
const bootstratGrpcPort = 19190
const poetPort = 10080
const poetRpcPort = 50002
const TestNetPath = "./local-testnet/"
const BinPath = TestNetPath+"bin/"

var globalOpts = map[string]interface{} {
	"--randcon":nodesCount,
	"--hare-committee-size":5,
	"--hare-max-adversaries":2,
	"--hare-round-duration-sec":10,
	"--layer-duration-sec":60,
	"--layer-average-size":10,
	"--hare-wakeup-delta": 10,
	"--test-mode":nil,
	"--eligibility-confidence-param":5,
	"--eligibility-epoch-offset":0,
	"--genesis-active-size":5,
	"--genesis-conf": TestNetPath+"genesis_accounts.json",
	//"--executable-path": BinPath+"go-spacemesh",
	"--genesis-time": "2020-06-20T04:27:15+00:00",
	"--coinbase": "097598942e44919cf7d11499887a595e41b097acd0a75d65ed8b8c6fa739d297",
	"--poet-server": "127.0.0.1:10080",
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

func bootsrap() (term Terminator) {

	_ = os.RemoveAll(TestNetPath+".data")
	_ = os.MkdirAll(TestNetPath+".data",0777)

	sigterm := make([]chan struct{}, nodesCount+2)
	for i := range sigterm { sigterm[i] = make(chan struct{}) }
	defer func() { for _,x := range sigterm { close(x) } }()

	opts := mergeOpts(map[string]interface{}{
		"--json-server":nil,
		"--json-port":bootstratJsonPort,
		"--grpc-server":nil,
		"--grpc-port":bootstratGrpcPort,
		"--tcp-port":bootstrapPort,
		"--data-folder": fmt.Sprintf(TestNetPath+".data/data.%v",0),
		"--post-datadir": fmt.Sprintf(TestNetPath+".data/post.%v",0),
	})
	c, err := exec(BinPath+"go-spacemesh",opts,sigterm[0]) ;
	if err != nil {
		panic(err)
	}

	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
			if M, ok := m["M"].(string); ok {
				fmt.Fprintln(os.Stderr,"##", s)
				const prefix = "Local node identity >> "
				if strings.HasPrefix(M,prefix) {
					id := M[len(prefix):]
					boot(fmt.Sprintf("spacemesh://%v@127.0.0.1:%v",id,bootstrapPort),sigterm[2:])
					break
				}
			}
		}
	}

	poet(sigterm[1])

	go collect(c)
	term = sigterm
	sigterm = nil
	return
}

func boot(booturl string, sigterm []chan struct{}) {
	for i, t := range sigterm {
		node(i+1,booturl,t)
	}
}

func node(no int, booturl string, sigterm chan struct{}) {
	opts := mergeOpts(map[string]interface{}{
		"--json-server":nil,
		"--json-port":bootstratJsonPort+no,
		"--grpc-server":nil,
		"--grpc-port":bootstratGrpcPort+no,
		"--tcp-port":bootstrapPort+no,
		"--start-mining":nil,
		"--bootstrap":nil,
		"--bootnodes": booturl+fmt.Sprintf("?disc=%v",bootstrapPort),
		"--data-folder": fmt.Sprintf(TestNetPath+".data/%v",no),
		"--post-datadir": fmt.Sprintf(TestNetPath+".data/post.%v",no),
	})
	c, err := exec(BinPath+"go-spacemesh",opts,sigterm) ;
	if err != nil {
		panic(err)
	}
	go collect(c)
	return
}

func poet(sigterm chan struct{}) {
	opts := map[string]interface{} {
		"--rpclisten":fmt.Sprintf("127.0.0.1:%v",poetRpcPort),
		"--restlisten":fmt.Sprintf("127.0.0.1:%v",poetPort),
		"--initialduration":"10s",
		"--duration":"10s",
		"--gateway":fmt.Sprintf("127.0.0.1:%v",bootstratGrpcPort),
	}
	c, err := exec(BinPath+"poet",opts,sigterm) ;
	if err != nil {
		panic(err)
	}
	go collect(c)
	return
}

func collect(c chan string) {
	for s := range c {
		fmt.Fprintln(os.Stderr,s)
	}
}
