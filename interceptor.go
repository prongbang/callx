package callx

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
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
	maxBodyLog int
}

func (l *loggerInterceptor) Request(req *fasthttp.Request) {
	method := string(req.Header.Method())
	url := string(req.URI().FullURI())
	body := req.Body()

	fmt.Printf("--> %s %s\n", method, url)

	req.Header.VisitAll(func(k, v []byte) {
		fmt.Printf("%s: %s\n", string(k), string(v))
	})

	if len(body) > 0 {
		fmt.Println()
		fmt.Println(l.truncateBody(body))
	}

	fmt.Printf("--> END %s (%d-byte body)\n\n", method, len(body))
}

func (l *loggerInterceptor) Response(res *fasthttp.Response) {
	status := res.StatusCode()
	body := res.Body()

	fmt.Printf("<-- %d %s\n", status, http.StatusText(status))

	res.Header.VisitAll(func(k, v []byte) {
		fmt.Printf("%s: %s\n", string(k), string(v))
	})

	if len(body) > 0 {
		fmt.Println()
		fmt.Println(l.truncateBody(body))
	}

	fmt.Printf("<-- END HTTP (%d-byte body)\n\n", len(body))
}

func (l *loggerInterceptor) truncateBody(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if len(b) > l.maxBodyLog {
		return fmt.Sprintf("%s\n... (truncated, showing %d of %d bytes)", string(b[:l.maxBodyLog]), l.maxBodyLog, len(b))
	}
	return string(b)
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
func LoggerInterceptor(maxBodyLog ...int) Interceptor {
	limit := 1024
	if len(maxBodyLog) > 0 {
		limit = maxBodyLog[0]
	}
	return &loggerInterceptor{
		maxBodyLog: limit,
	}
}
