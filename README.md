# CallX

HTTP Client easy call API

[![Build Status](http://img.shields.io/travis/prongbang/callx.svg)](https://travis-ci.org/prongbang/callx)
[![Codecov](https://img.shields.io/codecov/c/github/prongbang/callx.svg)](https://codecov.io/gh/prongbang/callx)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/callx)](https://goreportcard.com/report/github.com/prongbang/callx)
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/prongbang)

### Install

```
go get github.com/prongbang/callx
```

### Benchmark

```shell
Benchmark_CallXRequests/GET-10         	 2182681	       553.8 ns/op
Benchmark_CallXRequests/POST
Benchmark_CallXRequests/POST-10        	 1993790	       603.6 ns/op
Benchmark_CallXRequests/PUT
Benchmark_CallXRequests/PUT-10         	 1932697	       625.2 ns/op
Benchmark_CallXRequests/PATCH
Benchmark_CallXRequests/PATCH-10       	 1905487	       627.5 ns/op
Benchmark_CallXRequests/DELETE
Benchmark_CallXRequests/DELETE-10      	 2149676	       556.9 ns/op
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
