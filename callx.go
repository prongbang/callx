package callx

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseURL     string
	Timeout     time.Duration
	Interceptor []Interceptor
}

type Interceptor interface {
	Interceptor(req *http.Request, res *http.Response)
}

type Response struct {
	Code int
	Data []byte
}

type CallX interface {
	Get(url string) Response
	Post(url string, body map[string]interface{}) Response
	Patch(url string, body map[string]interface{}) Response
	Put(url string, body map[string]interface{}) Response
	Delete(url string) Response
	request(url string, method string, payload io.Reader) Response
}

type callxMethod struct {
	Config Config
}

func (n *callxMethod) Get(url string) Response {
	return n.request(url, http.MethodGet, nil)
}

func (n *callxMethod) Post(url string, body map[string]interface{}) Response {
	return n.request(url, http.MethodPost, getPayload(body))
}

func (n *callxMethod) Patch(url string, body map[string]interface{}) Response {
	return n.request(url, http.MethodPatch, getPayload(body))
}

func (n *callxMethod) Put(url string, body map[string]interface{}) Response {
	return n.request(url, http.MethodPut, getPayload(body))
}

func (n *callxMethod) Delete(url string) Response {
	return n.request(url, http.MethodDelete, nil)
}

func (n *callxMethod) request(url string, method string, payload io.Reader) Response {
	var ts time.Duration = 60
	if n.Config.Timeout > 0 {
		ts = n.Config.Timeout
	}
	client := &http.Client{
		Timeout: time.Second * ts,
	}
	endpoint := n.Config.BaseURL + url
	req, _ := http.NewRequest(method, endpoint, payload)

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	for _, inp := range n.Config.Interceptor {
		inp.Interceptor(req, res)
	}
	defer res.Body.Close()
	return Response{
		Code: res.StatusCode,
		Data: body,
	}
}

func getPayload(body map[string]interface{}) *strings.Reader {
	if body == nil {
		return nil
	}
	data, _ := json.Marshal(body)
	return strings.NewReader(string(data))
}

func New(config Config) CallX {
	return &callxMethod{
		Config: config,
	}
}
