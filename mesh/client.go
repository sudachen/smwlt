package mesh

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sudachen/smwlt/fu"
	"github.com/sudachen/smwlt/log"
	"io/ioutil"
	"net/http"
)

const DefaultEndpoint = "localhost:9090"
const DefaultVerbose = false

type ClinetAgent struct {
	verbose bool
	baseUrl string
	http    *http.Client
}

/*
Client describes mesh client options
*/
type Client struct {
	Endpoint string
	Verbose  bool
}

/*
New creates new mesh ClinetAgent
*/
func (c Client) New() *ClinetAgent {
	return &ClinetAgent{
		baseUrl: "http://" + fu.Fne(c.Endpoint, DefaultEndpoint) + "/v1",
		verbose: fu.Fnf(c.Verbose, DefaultVerbose),
		http:    &http.Client{},
	}
}

func (c *ClinetAgent) post(api string, in, out interface{}) (err error) {
	data, err := json.Marshal(in)
	if err != nil {
		return
	}
	url := c.baseUrl + api
	if c.verbose {
		log.Info("request: %v, body: %s", url, data)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	if c.verbose {
		log.Info("response body: %s", resBody)
	}

	if res.StatusCode != http.StatusOK {
		rb := struct {
			Error string `json:"error"`
		}{}
		json.Unmarshal(resBody, &rb)
		err = fu.Wrapf(errors.New(rb.Error), "`%v` response status code: %d", api, res.StatusCode)
		if c.verbose {
			log.Error("request failed with error " + err.Error())
		}
		return
	}

	return json.Unmarshal(resBody, out)
}

func (c *ClinetAgent) getValue64(api string, in interface{}) (uint64, error) {
	out := struct {
		Value uint64 `json:"value,string"`
	}{}
	err := c.post(api, in, &out)
	return out.Value, err
}

func (c *ClinetAgent) getValueBool(api string, in interface{}) (bool, error) {
	out := struct {
		Value bool `json:"value"`
	}{}
	err := c.post(api, in, &out)
	return out.Value, err
}
