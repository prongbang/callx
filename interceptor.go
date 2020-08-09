package callx

import (
	"log"
	"net/http"
)

type headerInterceptor struct {
	Header Header
}

type loggerInterceptor struct {
}

type jsonContentTypeInterceptor struct {
}

func (l *loggerInterceptor) Interceptor(req *http.Request) {
	log.Println("-->", req.Method, req.URL.String())
	for k, v := range req.Header {
		head := k
		value := ""
		for _, val := range v {
			value += val + " "
		}
		log.Print(head, ": ", value)
	}
	log.Println("-->", "END")
}

func (j *jsonContentTypeInterceptor) Interceptor(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (a *headerInterceptor) Interceptor(req *http.Request) {
	for k, v := range a.Header {
		req.Header.Set(k, v)
	}
}

func HeaderInterceptor(header Header) Interceptor {
	return &headerInterceptor{
		Header: header,
	}
}

func JSONContentTypeInterceptor() Interceptor {
	return &jsonContentTypeInterceptor{}
}

func LoggerInterceptor() Interceptor {
	return &loggerInterceptor{}
}
