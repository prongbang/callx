# CallX

HTTP Client easy call API

[![Build Status](http://img.shields.io/travis/prongbang/callx.svg)](https://travis-ci.org/prongbang/callx)
[![Codecov](https://img.shields.io/codecov/c/github/prongbang/callx.svg)](https://codecov.io/gh/prongbang/callx)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/callx)](https://goreportcard.com/report/github.com/prongbang/callx)


```
go get github.com/prongbang/callx
```

### How to use

- Using base URL

```golang
c := callx.Config{
    BaseURL: "https://jsonplaceholder.typicode.com",
    Timeout: 60,
}
req := callx.New(c)

data := req.Get("/todos/1")
fmt.Println(string(data.Data))
```

- Custom request

```golang
c := callx.Config{
    Timeout: 60,
}
req := callx.New(c)

custom := callx.Custom{
    URL:    "https://httpbin.org/post",
    Method: http.MethodPost,
    Header: callx.Header{
        callx.Authorization: fmt.Sprintf("%s %s", callx.Bearer, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Et9HFtf9R3GEMA0IICOfFMVXY7kkTX1wr4qCyhIf58U"),
    },
    Body: callx.Body{
        "username": "root",
        "password": "pass",
        "address": []string{
            "087654321",
            "089786756",
        },
    },
}
data := req.Req(custom)
fmt.Println(string(data.Data))
```

- Custom request form encoded

```golang
c := callx.Config{
    Timeout: 60,
}
req := callx.New(c)

form := url.Values{}
form.Set("message", "Test")

custom := callx.Custom{
    URL:    "https://httpbin.org/post",
    Method: http.MethodPost,
    Header: callx.Header{
        callx.Authorization: "Bearer XTZ",
        callx.ContentType:   "application/x-www-form-urlencoded",
    },
    Form: strings.NewReader(form.Encode()),
}
data := req.Req(custom)
fmt.Println(string(data.Data))
```

### Benchmark

- CallX

Before

```shell
Benchmark_HTTPRequests/GET
Benchmark_HTTPRequests/GET-10         	   35385	     32224 ns/op
Benchmark_HTTPRequests/POST
Benchmark_HTTPRequests/POST-10        	   34557	     33460 ns/op
Benchmark_HTTPRequests/PUT
Benchmark_HTTPRequests/PUT-10         	   36110	     33648 ns/op
Benchmark_HTTPRequests/DELETE
Benchmark_HTTPRequests/DELETE-10      	   37135	     31826 ns/op
```

After

```shell
Benchmark_CallXRequests/GET
Benchmark_CallXRequests/GET-10         	 2236383	       538.9 ns/op
Benchmark_CallXRequests/POST
Benchmark_CallXRequests/POST-10        	 2023060	       588.9 ns/op
Benchmark_CallXRequests/PUT
Benchmark_CallXRequests/PUT-10         	 1943686	       613.7 ns/op
Benchmark_CallXRequests/DELETE
Benchmark_CallXRequests/DELETE-10      	 2199108	       540.7 ns/op
```