# CallX 🚀

[![Codecov](https://img.shields.io/codecov/c/github/prongbang/callx.svg)](https://codecov.io/gh/prongbang/callx)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/callx)](https://goreportcard.com/report/github.com/prongbang/callx)
[![Go Reference](https://pkg.go.dev/badge/github.com/prongbang/callx.svg)](https://pkg.go.dev/github.com/prongbang/callx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> A lightweight, fast, and easy-to-use HTTP client for Go. Make API calls with just a few lines of code!

## ✨ Features

- 🚀 **Ultra-fast performance** - Optimized for speed
- 🛠 **Simple API** - Easy to learn and use
- 🔧 **Highly customizable** - Full control over requests
- 🧪 **Well tested** - High test coverage

## ⚡️ Performance

CallX is designed for optimal performance. Here are the benchmark results:

| HTTP Method  | Operations | Time per Operation |
|-------------|------------|-------------------|
| GET         | 41,756     | 31,823 ns/op      |
| POST        | 38,692     | 35,787 ns/op      |
| POST-ENCODE | 28,848     | 39,314 ns/op      |
| PUT         | 31,401     | 35,046 ns/op      |
| PATCH       | 38,923     | 30,094 ns/op      |
| DELETE      | 41,100     | 29,195 ns/op      |

## 📦 Installation

```bash
go get github.com/prongbang/callx
```

## 🚀 Quick Start

### Basic Usage

```go
// Create a client with base URL
c := callx.Config{
    BaseURL: "https://jsonplaceholder.typicode.com",
    Timeout: 60,
}
req := callx.New(c)

// Make a GET request
data := req.Get("/todos/1")
fmt.Println(string(data.Data))
```

## 🔥 Advanced Features

### Custom Request with Authentication

```go
c := callx.Config{
    Timeout: 60,
}
req := callx.New(c)

custom := callx.Custom{
    URL:    "https://httpbin.org/post",
    Method: http.MethodPost,
    Header: callx.Header{
        callx.Authorization: fmt.Sprintf("%s %s", callx.Bearer, "your-token"),
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

### Form-encoded Requests

```go
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
