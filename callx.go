package callx

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var interceptors []Interceptor

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
type Body map[string]interface{}

// Custom callx request model
type Custom struct {
	URL    string
	Method string
	Header Header
	Body   interface{}
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
	Post(url string, body interface{}) Response
	Patch(url string, body interface{}) Response
	Put(url string, body interface{}) Response
	Delete(url string) Response
	Req(custom Custom) Response
	AddInterceptor(req *http.Request)
	request(urlStr string, method string, header Header, payload io.Reader) Response
}

type callxMethod struct {
	Config Config
	Client *http.Client // Reuse the same client for connection pooling
}

func (n *callxMethod) Get(url string) Response {
	return n.request(url, http.MethodGet, nil, nil)
}

func (n *callxMethod) Post(url string, body interface{}) Response {
	return n.request(url, http.MethodPost, nil, getPayload(body))
}

func (n *callxMethod) Patch(url string, body interface{}) Response {
	return n.request(url, http.MethodPatch, nil, getPayload(body))
}

func (n *callxMethod) Put(url string, body interface{}) Response {
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

func (n *callxMethod) AddInterceptor(req *http.Request) {
	for _, interceptor := range interceptors {
		interceptor.Interceptor(req)
	}
}

func (n *callxMethod) request(urlStr string, method string, header Header, payload io.Reader) Response {
	resNotFound := Response{Code: http.StatusNotFound}

	endpointURL, err := url.Parse(urlStr)
	if err != nil {
		return resNotFound
	}

	if !isURL(urlStr) && n.Config.BaseURL != "" {
		endpointURL.Scheme = "http" // You can set this to "https" if needed
		endpointURL.Host = n.Config.BaseURL
	}

	req, err := http.NewRequest(method, endpointURL.String(), payload)
	if err != nil {
		return resNotFound
	}

	n.AddInterceptor(req)
	setHeaders(req, header)

	res, err := n.Client.Do(req)
	if err != nil {
		return resNotFound
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, err := io.ReadAll(res.Body)
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

func getPayload(body interface{}) *strings.Reader {
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
		Client: &http.Client{
			Timeout: time.Second * config.Timeout,
		},
	}
}
