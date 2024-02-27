package callx

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type headerInterceptor struct {
	Header Header
}

func (a *headerInterceptor) Response(res *fasthttp.Response) {
}

func (a *headerInterceptor) Request(req *fasthttp.Request) {
	for k, v := range a.Header {
		req.Header.Set(k, v)
	}
}

type loggerInterceptor struct {
}

func (l *loggerInterceptor) Response(res *fasthttp.Response) {
}

func (l *loggerInterceptor) Request(req *fasthttp.Request) {
	fmt.Println("-->", string(req.Header.Method()), string(req.Header.RequestURI()))
	fmt.Println("-->", req.Header.RawHeaders())
	fmt.Println("-->", "END")
}

type jsonContentTypeInterceptor struct {
}

func (j *jsonContentTypeInterceptor) Response(res *fasthttp.Response) {
}

func (j *jsonContentTypeInterceptor) Request(req *fasthttp.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

// HeaderInterceptor provide a instance
func HeaderInterceptor(header Header) Interceptor {
	return &headerInterceptor{
		Header: header,
	}
}

// JSONContentTypeInterceptor provide a instance
func JSONContentTypeInterceptor() Interceptor {
	return &jsonContentTypeInterceptor{}
}

// LoggerInterceptor provide a instance
func LoggerInterceptor() Interceptor {
	return &loggerInterceptor{}
}
