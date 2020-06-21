package go_testnet

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type pidmapVal struct{Name string;No int}
type pidmap map[string]pidmapVal
var p2pid = pidmap{}
var muP2pid = sync.Mutex{}

func nameOf(no int) string {
	if no != 0 {
		return fmt.Sprintf("NODE[%v]",no)
	}
	return "BOOTSTRAP"
}

func (m pidmap) Set(pid string, no int) {
	muP2pid.Lock()
	m[pid] =  pidmapVal{ nameOf(no), no }
	muP2pid.Unlock()
}

func (m pidmap) GetString(pid string) (r string){
	muP2pid.Lock()
	r = m[pid].Name
	muP2pid.Unlock()
	if r == "" { return "{"+pid+"}" }
	return
}

const localIdentity = "Local node identity >> "

func parse(m map[string]interface{}, no int, prefix string) {
	if M, ok := m["M"].(string); ok {
		if M == "PoST initialization completed" {
			fmt.Printf(prefix + "PoST initialization completed\n")
		} else if M == "message sent" {
			fmt.Printf(prefix+"Message sent %v\n",m["msg_type"])
		/*} else if M == "new_connection" {
			fmt.Printf(prefix+"Connection %v => %v\n",
					p2pid.GetString(m["src"].(string)),
					p2pid.GetString(m["dst"].(string)))*/
		/*} else if M == "now connected" {
			fmt.Printf(prefix+"Connected %v peers\n", m["n_peers"])*/
		} else if strings.HasPrefix(M,localIdentity) {
			p2pid.Set(M[len(localIdentity):],no)
		} else {
			//fmt.Println(prefix,M,m)
		}
	}
}

func collect(c chan string,no int) {
	prefix := nameOf(no)+":"
	for s := range c {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(s),&m); err == nil {
			parse(m, no, prefix)
		}
	}
}

func collect_poet(c chan string) {
	for s := range c {
		_ = s
		//fmt.Println("POET:"+s)
	}
}
