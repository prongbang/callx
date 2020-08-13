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

// Constant Header
const (
	Authorization = "Authorization"
	ContentType   = "Content-Type"
	Accept        = "Accept"
	Basic         = "Basic"
	Bearer        = "Bearer"
)

// Header custom type
type Header map[string]string

// Body custom type
type Body map[string]string

// Custom callx request model
type Custom struct {
	URL    string
	Method string
	Header Header
	Body   Body
	Form   io.Reader
}

// Config callx model
type Config struct {
	BaseURL     string
	Timeout     time.Duration
	Interceptor []Interceptor
}

// Interceptor the interface
type Interceptor interface {
	Interceptor(req *http.Request)
}

// Response callx model
type Response struct {
	Code int
	Data []byte
}

// CallX the interface
type CallX interface {
	Get(url string) Response
	Post(url string, body Body) Response
	Patch(url string, body Body) Response
	Put(url string, body Body) Response
	Delete(url string) Response
	Req(custom Custom) Response
	AddInterceptor(intercept ...Interceptor)
	request(url string, method string, header Header, payload io.Reader) Response
}

type callxMethod struct {
	Config Config
}

func (n *callxMethod) Get(url string) Response {
	return n.request(url, http.MethodGet, nil, nil)
}

func (n *callxMethod) Post(url string, body Body) Response {
	return n.request(url, http.MethodPost, nil, getPayload(body))
}

func (n *callxMethod) Patch(url string, body Body) Response {
	return n.request(url, http.MethodPatch, nil, getPayload(body))
}

func (n *callxMethod) Put(url string, body Body) Response {
	return n.request(url, http.MethodPut, nil, getPayload(body))
}

func (n *callxMethod) Delete(url string) Response {
	return n.request(url, http.MethodDelete, nil, nil)
}

func (n *callxMethod) Req(custom Custom) Response {
	if custom.Form != nil {
		return n.request(custom.URL, custom.Method, custom.Header, custom.Form)
	}
	return n.request(custom.URL, custom.Method, custom.Header, getPayload(custom.Body))
}

func (n *callxMethod) AddInterceptor(intercept ...Interceptor) {
	for _, ins := range intercept {
		interceptors = append(interceptors, ins)
	}
}

func (n *callxMethod) request(url string, method string, header Header, payload io.Reader) Response {
	resNotFound := Response{Code: http.StatusNotFound}
	var ts time.Duration = 60
	if n.Config.Timeout > 0 {
		ts = n.Config.Timeout
	}
	client := &http.Client{
		Timeout: time.Second * ts,
	}
	endpoint := n.Config.BaseURL + url
	if isURL(url) {
		endpoint = url
	}
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return resNotFound
	}
	n.AddInterceptor(n.Config.Interceptor...)
	for _, inp := range interceptors {
		inp.Interceptor(req)
	}
	setHeaders(req, header)

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

func setHeaders(req *http.Request, header Header) {
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
}

func getPayload(body Body) *strings.Reader {
	data, err := json.Marshal(body)
	if err != nil {
		return strings.NewReader("")
	}
	return strings.NewReader(string(data))
}

func isURL(url string) bool {
	return strings.Index(url, "http://") > -1 || strings.Index(url, "https://") > -1
}

// New callx
func New(config Config) CallX {
	return &callxMethod{
		Config: config,
	}
}
