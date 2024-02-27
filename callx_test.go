package callx_test

import (
	"fmt"
	"github.com/prongbang/callx"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Get")
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

	data := req.Get("/todos/1?q=callx")
	if data.Code != 200 {
		t.Error("CallX Get Error", data)
	}
}

func Test_PostBodyNil(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		Timeout: 60,
	}
	req := callx.New(c)

	data := req.Post(fmt.Sprintf("%s/todos?q=callx&type=http", ts.URL), nil)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_PostBodyError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := map[string]interface{}{
		"error": make(chan int),
	}
	data := req.Post("/todos", body)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_PostServerNotFound(t *testing.T) {
	c := callx.Config{
		BaseURL: "http://localhost:1234/todos",
		Timeout: 60,
	}
	req := callx.New(c)

	data := req.Post("/todos", nil)
	if data.Code != 404 {
		t.Error("CallX Post Error")
	}
}

func Test_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := map[string]interface{}{}
	data := req.Post("/todos", body)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_Post_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := []interface{}{"1", "2", "3", "4"}
	data := req.Post("/todos", body)
	if data.Code != 200 {
		t.Error("CallX Post Error")
	}
}

func Test_Put(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Put")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := map[string]interface{}{}
	data := req.Put("/todos/1", body)
	if data.Code != 200 {
		t.Error("CallX Put Error")
	}
}

func Test_Patch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Patch")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
	}
	req := callx.New(c)

	body := map[string]interface{}{}
	data := req.Patch("/todos/1", body)
	if data.Code != 200 {
		t.Error("CallX Patch Error")
	}
}

func Test_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Delete")
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
		_, _ = fmt.Fprintln(w, "Hello, Req Post")
	}))
	defer ts.Close()

	c := callx.Config{
		Timeout: 60,
		Interceptor: []callx.Interceptor{
			callx.HeaderInterceptor(callx.Header{
				"X-TOKEN": "Bearer XYZ",
			}),
		},
	}
	req := callx.New(c)

	req.AddInterceptor(callx.LoggerInterceptor())

	custom := callx.Custom{
		URL:    fmt.Sprintf("%s/post", ts.URL),
		Method: http.MethodPost,
		Header: callx.Header{
			callx.Authorization: fmt.Sprintf("%s %s", callx.Bearer, "eyJh9.e30.EtU"),
		},
		Body: callx.Body(map[string]interface{}{
			"username": "root",
			"password": "pass",
		}),
	}
	data := req.Req(custom)
	if data.Code != 200 {
		t.Error("CallX Req Post Error")
	}
}

func Test_ReqEncoded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Req Post")
	}))
	defer ts.Close()

	c := callx.Config{
		Timeout: 60,
	}
	req := callx.New(c)

	form := url.Values{}
	form.Set("message", "Test")
	custom := callx.Custom{
		URL:    fmt.Sprintf("%s/post", ts.URL),
		Method: http.MethodPost,
		Header: callx.Header{
			callx.Authorization: "Bearer XTZ",
			callx.ContentType:   "application/x-www-form-urlencoded",
		},
		Form: strings.NewReader(form.Encode()),
	}

	data := req.Req(custom)
	if data.Code != 200 {
		t.Error("CallX Req Post Error")
	}
}

func Test_ReqMethodNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, Post")
	}))
	defer ts.Close()

	c := callx.Config{
		Timeout: 60,
	}
	req := callx.New(c)
	custom := callx.Custom{
		URL:    fmt.Sprintf("%s/post", ts.URL),
		Method: "!#@!@",
	}
	data := req.Req(custom)
	if data.Code != 400 {
		t.Error("CallX Post Error", string(data.Data))
	}
}

func Benchmark_CallXRequests(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, World")
	}))
	defer ts.Close()

	c := callx.Config{
		BaseURL: ts.URL,
		Timeout: 60,
		Interceptor: []callx.Interceptor{
			callx.HeaderInterceptor(callx.Header{
				callx.Authorization: "Bearer XTZ",
			}),
		},
	}
	req := callx.New(c)

	b.Run("GET", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = req.Get("/todos/1?id=1")
		}
	})

	b.Run("POST", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = req.Post("/todos?id=1", callx.Body(map[string]any{"body": "post"}))
		}
	})

	b.Run("POST-ENCODE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			form := url.Values{}
			form.Set("message", "Test")
			custom := callx.Custom{
				URL:    fmt.Sprintf("%s/post", ts.URL),
				Method: http.MethodPost,
				Header: callx.Header{
					callx.Authorization: "Bearer XTZ",
					callx.ContentType:   "application/x-www-form-urlencoded",
				},
				Form: strings.NewReader(form.Encode()),
			}
			_ = req.Req(custom)
		}
	})

	b.Run("PUT", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = req.Put("/todos/1?id=1", callx.Body(map[string]any{"body": "put"}))
		}
	})

	b.Run("PATCH", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = req.Patch("/todos/1?id=1", callx.Body(map[string]any{"body": "patch"}))
		}
	})

	b.Run("DELETE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = req.Delete("/todos/1?id=1")
		}
	})
}
