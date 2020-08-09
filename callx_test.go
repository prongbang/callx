package callx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prongbang/callx"
)

func Test_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Get")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
		Interceptor: []callx.Interceptor{
			callx.LoggerInterceptor(),
		},
	}
	req := callx.New(c)

	data := req.Get("/todos/1")
	if data.Code != 200 {
		t.Error("CallX Get Error")
	}
}

func Test_PostBodyNil(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	data := req.Post("/todos", nil)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := callx.Body{}
	data := req.Post("/todos", body)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_Put(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Put")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := callx.Body{}
	data := req.Put("/todos/1", body)
	if data.Code != 200 {
		t.Error("CallX Put Error")
	}
}

func Test_Patch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Patch")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := callx.Body{}
	data := req.Patch("/todos/1", body)
	if data.Code != 200 {
		t.Error("CallX Patch Error")
	}
}

func Test_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Delete")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	data := req.Delete("/todos/1")
	if data.Code != 200 {
		t.Error("CallX Delete Error")
	}
}

func Test_Req(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Req Post")
	}))
	defer ts.Close()

	c := callx.Config{
		Timeout: 60,
	}
	req := callx.New(c)

	custom := callx.Custom{
		URL:    ts.URL + "/post",
		Method: http.MethodPost,
		Header: callx.Header{
			callx.Authorization: fmt.Sprintf("%s %s", callx.Bearer, "eyJh9.e30.EtU"),
		},
		Body: callx.Body{
			"username": "root",
			"password": "pass",
		},
	}
	data := req.Req(custom)
	if data.Code != 200 {
		t.Error("CallX Req Post Error")
	}
}