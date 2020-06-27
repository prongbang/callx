# CallX

HTTP Client easy call API

```
go get github.com/prongbang/callx
```

### How to use

```golang
c := callx.Config{
    BaseURL: "https://jsonplaceholder.typicode.com",
    Timeout: 60,
}
req := callx.New(c)

data := req.Get("/todos/1")
fmt.Println(string(data.Data))
```