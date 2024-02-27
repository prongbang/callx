package callx_test

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"

	"github.com/prongbang/callx"
)

func Test_LoggerInterceptor(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodGet)
	req.Header.SetRequestURI("/todo")
	header := callx.HeaderInterceptor(callx.Header{
		callx.Authorization: "Bearer Token",
	})
	header.Request(req)
	logger := callx.LoggerInterceptor()
	logger.Request(req)
	if string(req.Header.RequestURI()) != "/todo" {
		t.Error("Request not found")
	}
}

func Test_JSONContentTypeInterceptor(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodGet)
	req.Header.SetRequestURI("/todo")
	jsonIntercept := callx.JSONContentTypeInterceptor()
	jsonIntercept.Request(req)
	if string(req.Header.RequestURI()) != "/todo" {
		t.Error("Request not found")
	}
}
