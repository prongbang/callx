package callx_test

import (
	"net/http/httptest"
	"testing"

	"github.com/prongbang/callx"
)

func Test_LoggerInterceptor(t *testing.T) {
	req := httptest.NewRequest("GET", "/todo", nil)
	header := callx.HeaderInterceptor(callx.Header{
		callx.Authorization: "Bearer Token",
	})
	header.Interceptor(req)
	logger := callx.LoggerInterceptor()
	logger.Interceptor(req)
	if req.URL.Path != "/todo" {
		t.Error("Request not found")
	}
}

func Test_JSONContentTypeInterceptor(t *testing.T) {
	req := httptest.NewRequest("GET", "/todo", nil)
	jsonIntercept := callx.JSONContentTypeInterceptor()
	jsonIntercept.Interceptor(req)
	if req.URL.Path != "/todo" {
		t.Error("Request not found")
	}
}
