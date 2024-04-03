package callx

import (
	"bufio"
	"crypto/tls"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/url"
	"time"
)

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
type Body interface{}

// Form custom type
type Form io.Reader

// Custom callx request model
type Custom struct {
	URL    string
	Method string
	Header Header
	Body   Body
	Form   Form
}

// Config callx model
type Config struct {
	BaseURL string

	// Client name. Used in User-Agent request header.
	//
	// Default client name is used if not set.
	Name string

	// Maximum duration for full request writing and response reading (including body).
	//
	Timeout time.Duration

	// Maximum duration for full response reading (including body).
	//
	// By default response read timeout is unlimited.
	ReadTimeout time.Duration

	// Maximum duration for full request writing (including body).
	//
	// By default request write timeout is unlimited.
	WriteTimeout time.Duration

	// Idle keep-alive connections are closed after this duration.
	//
	// By default idle connections are closed after DefaultMaxIdleConnDuration.
	MaxIdleConnDuration time.Duration

	Interceptor []Interceptor

	// TLS config for https connections.
	//
	// Default TLS config is used if not set.
	TLSConfig *tls.Config

	// InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
	InsecureSkipVerify bool

	// TCPDialer contains options to control a group of Dial calls.
	TCPDialer *fasthttp.TCPDialer

	// Maximum number of connections per each host which may be established.
	//
	// DefaultMaxConnsPerHost is used if not set.
	MaxConnsPerHost int

	// Per-connection buffer size for responses' reading.
	// This also limits the maximum header size.
	//
	// Default buffer size is used if 0.
	ReadBufferSize int

	// Per-connection buffer size for requests' writing.
	//
	// Default buffer size is used if 0.
	WriteBufferSize int

	// RetryIf controls whether a retry should be attempted after an error.
	//
	// By default will use isIdempotent function.
	RetryIf fasthttp.RetryIfFunc

	// StreamResponseBody enables response body streaming.
	StreamResponseBody bool
}

// Interceptor the interface
type Interceptor interface {
	Request(req *fasthttp.Request)
	Response(res *fasthttp.Response)
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
	request(urlStr string, method string, header Header, payload interface{}) Response
}

type callxMethod struct {
	Config *Config
	Client *fasthttp.Client
}

func (n *callxMethod) Get(url string) Response {
	return n.request(url, http.MethodGet, nil, nil)
}

func (n *callxMethod) Post(url string, body Body) Response {
	return n.request(url, http.MethodPost, nil, body)
}

func (n *callxMethod) Patch(url string, body Body) Response {
	return n.request(url, http.MethodPatch, nil, body)
}

func (n *callxMethod) Put(url string, body Body) Response {
	return n.request(url, http.MethodPut, nil, body)
}

func (n *callxMethod) Delete(url string) Response {
	return n.request(url, http.MethodDelete, nil, nil)
}

func (n *callxMethod) Req(custom Custom) Response {
	if custom.Form != nil {
		return n.request(custom.URL, custom.Method, custom.Header, custom.Form)
	}
	return n.request(custom.URL, custom.Method, custom.Header, custom.Body)
}

func (n *callxMethod) AddInterceptor(intercept ...Interceptor) {
	for _, ins := range intercept {
		n.Config.Interceptor = append(n.Config.Interceptor, ins)
	}
}

func (n *callxMethod) request(urlStr string, method string, header Header, payload interface{}) Response {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	resNotFound := Response{Code: http.StatusNotFound}

	if n.Config.BaseURL != "" {
		urlStr = n.Config.BaseURL + urlStr
	}
	endpointURL, err := url.Parse(urlStr)
	if err != nil {
		return resNotFound
	}
	req.SetRequestURI(endpointURL.String())
	req.Header.SetMethod(method)
	setRequestInterceptor(req, n.Config.Interceptor)
	setHeaders(req, header)

	if payload != nil {
		if form, fok := payload.(Form); fok {
			req.SetBodyStreamWriter(func(w *bufio.Writer) {
				defer func(w *bufio.Writer) {
					_ = w.Flush()
				}(w)
				_, _ = w.ReadFrom(form)
			})
		} else {
			if data, e := json.Marshal(payload); e == nil {
				req.SetBodyRaw(data)
			}
		}
	}

	err = n.Client.Do(req, resp)
	if err != nil {
		return resNotFound
	}
	setResponseInterceptor(resp, n.Config.Interceptor)

	return Response{
		Code: resp.StatusCode(),
		Data: resp.Body(),
	}
}

func setRequestInterceptor(req *fasthttp.Request, interceptors []Interceptor) {
	for _, interceptor := range interceptors {
		interceptor.Request(req)
	}
}

func setResponseInterceptor(res *fasthttp.Response, interceptors []Interceptor) {
	for _, interceptor := range interceptors {
		interceptor.Response(res)
	}
}

func setHeaders(req *fasthttp.Request, header Header) {
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
}

// New callx
func New(config Config) CallX {

	// increase DNS cache time to an hour instead of default minute
	tcpDialer := &fasthttp.TCPDialer{
		Concurrency:      4096,
		DNSCacheDuration: time.Hour,
	}
	if config.TCPDialer != nil {
		tcpDialer = config.TCPDialer
	}

	client := &fasthttp.Client{
		Name:                          config.Name,
		ReadTimeout:                   time.Second * 30,
		WriteTimeout:                  time.Second * 30,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		MaxIdleConnDuration:           time.Hour * 1,
		Dial:                          tcpDialer.Dial,
		ReadBufferSize:                config.ReadBufferSize,
		WriteBufferSize:               config.WriteBufferSize,
		RetryIf:                       config.RetryIf,
		StreamResponseBody:            config.StreamResponseBody,
	}

	if config.MaxConnsPerHost > 0 {
		client.MaxConnsPerHost = config.MaxConnsPerHost
	}

	if config.MaxIdleConnDuration > 0 {
		client.MaxIdleConnDuration = config.MaxIdleConnDuration
	}

	if config.Timeout > 0 {
		client.ReadTimeout = time.Second * config.Timeout
		client.WriteTimeout = time.Second * config.Timeout
	} else {
		if config.ReadTimeout > 0 {
			client.ReadTimeout = time.Second * config.ReadTimeout
			config.Timeout = time.Second * config.ReadTimeout
		}

		if config.WriteTimeout > 0 {
			client.WriteTimeout = time.Second * config.WriteTimeout
			config.Timeout = time.Second * config.WriteTimeout
		}

		if config.TLSConfig != nil {
			client.TLSConfig = config.TLSConfig
		} else if config.InsecureSkipVerify {
			client.TLSConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}
		}
	}

	return &callxMethod{
		Config: &config,
		Client: client,
	}
}
