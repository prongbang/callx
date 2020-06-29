package callx

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var interceptors []Interceptor = []Interceptor{}

type Config struct {
	BaseURL     string
	Timeout     time.Duration
	Interceptor []Interceptor
}

type Interceptor interface {
	Interceptor(req *http.Request)
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
	AddInterceptor(intercept ...Interceptor)
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

func (n *callxMethod) AddInterceptor(intercept ...Interceptor) {
	for _, ins := range intercept {
		interceptors = append(interceptors, ins)
	}
}

func (n *callxMethod) request(url string, method string, payload io.Reader) Response {
	resNotFound := Response{Code: http.StatusNotFound}
	var ts time.Duration = 60
	if n.Config.Timeout > 0 {
		ts = n.Config.Timeout
	}
	client := &http.Client{
		Timeout: time.Second * ts,
	}
	endpoint := n.Config.BaseURL + url
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return resNotFound
	}
	req.URL.RawQuery = req.URL.Query().Encode()
	n.AddInterceptor(n.Config.Interceptor...)
	for _, inp := range interceptors {
		inp.Interceptor(req)
	}

	res, err := client.Do(req)
	if err != nil {
		return resNotFound
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resNotFound
	}
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
