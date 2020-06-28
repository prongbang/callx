package callx

import (
	"log"
	"net/http"
)

type loggerInterceptor struct {
}

type jsonContentTypeInterceptor struct {
}

func (l *loggerInterceptor) Interceptor(req *http.Request) {
	log.Println("-->", req.Method, req.URL.Path)
	log.Println(req.URL.RawQuery)
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
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func JSONContentType() Interceptor {
	return &jsonContentTypeInterceptor{}
}

func Logger() Interceptor {
	return &loggerInterceptor{}
}
