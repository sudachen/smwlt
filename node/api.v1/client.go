package api_v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sudachen/smwlt/fu"
	"io/ioutil"
	"net/http"
)

const DefaultEndpoint = "localhost:9090"

type ClientAgent struct {
	verbose func(string, ...interface{})
	baseUrl string
	http    *http.Client
}

/*
Client describes node client options
*/
type Client struct {
	Endpoint string
	Verbose  func(string, ...interface{})
}

/*
New creates new node ClientAgent
*/
func (c Client) New() *ClientAgent {
	verbose := c.Verbose
	if verbose == nil {
		verbose = func(string, ...interface{}) {}
	}
	return &ClientAgent{
		baseUrl: "http://" + fu.Fne(c.Endpoint, DefaultEndpoint) + "/v1",
		verbose: verbose,
		http:    &http.Client{},
	}
}

func (c *ClientAgent) post(api string, in, out interface{}) (err error) {
	data, err := json.Marshal(in)
	if err != nil {
		return
	}
	url := c.baseUrl + api
	c.verbose("request: %v, body: %s", url, data)
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
	c.verbose("response body: %s", resBody)

	if res.StatusCode != http.StatusOK {
		rb := struct {
			Error string `json:"error"`
		}{}
		if json.Unmarshal(resBody, &rb) == nil && rb.Error != "" {
			err = errors.New(rb.Error)
		} else {
			err = fmt.Errorf("`%v` response status code: %d", api, res.StatusCode)
		}
		c.verbose("request failed with error " + err.Error())
		return
	}

	return json.Unmarshal(resBody, out)
}

func (c *ClientAgent) getValue64(api string, in interface{}) (uint64, error) {
	out := struct {
		Value uint64 `json:"value,string"`
	}{}
	err := c.post(api, in, &out)
	return out.Value, err
}

func (c *ClientAgent) getValueBool(api string, in interface{}) (bool, error) {
	out := struct {
		Value bool `json:"value"`
	}{}
	err := c.post(api, in, &out)
	return out.Value, err
}
