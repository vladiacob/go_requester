# go_requester
Go requester is a package which is helping you to make HTTP requests easy than before.

## Travis CI
![Travis CI](https://travis-ci.org/vladiacob/go_requester.svg)

## Godoc
[https://godoc.org/github.com/vladiacob/go_requester](https://godoc.org/github.com/vladiacob/go_requester)

## How to install
```
go get github.com/vladiacob/go_requester
```

## How to use
### Include go_requester
```
include (
    ..
    requester "github.com/vladiacob/go_requester"
    ..
)
```

### Initialize requester
```
requester := requester.New(http.DefaultClient)
requester.SerUserAgent("test")
```

### Basic auth
```
requester.SetAuthentication("username", "password")
```

### Make request
#### JSON response
```
var clientJSONResponse ClientResponse
response, err := requester.Make("GET", "http://localhost", map[string]string{}, &clientJSONResponse)
```

#### String response
```
var clientStringResponse string
response, err = requester.Make("GET", "http://localhost", map[string]string{}, &clientStringResponse)
```

#### Response
```
type Response struct {
    Status int
    Body   []byte
}
```

