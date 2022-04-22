# http.Handle

There are two main ways for handling incoming HTTP requests:
## `http.Handle`
This registers a handler to handle requests matching a pattern.
``` go
func Handle(pattern string, handler Handler)

type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

``` go
import "net/http"

// 1. create a type, which implements the `Handler` interface.
type fooHandler struct {
    Message string
}
func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(f.Message))
}

// 2. register the handler
func main() {
    http.Handle("/foo", &fooHandler{Message: "hello world"})
}
```

## `http.HandleFunc`
This registers a function to handle requests matching a pattern.
``` go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

``` go
import "net/http"

func main() {
    foo := func(w http.ResponseWriter, _ *http.Response) {
        w.Write([]byte(f.Message))
    }

    http.HandleFunc("/foo", foo)
}
```

# ServeMux
HTTP Request Multiplexor

It's responsible for taking the incoming URL of an incoming request, and matching it to the list of registered patterns, then calling hte handler which was registered with that pattern.

When trying to match the incoming URL, it will always try to find the best match from the list of registered patterns.

For example, if we had the following registered handlers; `/hello/world`, `/hello/world/`, and then received a request `/hello/world/123`, this would be handled by the `/hello/world/`, as this is the closest match to the incoming URL.

## `http.ListenAndServe`
``` go
func ListenAndServe(addr string, handler Handler) error
```

``` go
import "net/http"

func main() {
    foo := func(w http.ResponseWriter, _ *http.Response) {
        w.Write([]byte("hello world"))
    }
    http.HandleFunc("/foo", foo)
    err := http.ListenAndServe(":5000", nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

Calling `http.ListenAndServe` with a `nil` will use the default ServeMux.

Because `http.ListenAndServe` blocks, and only returns non-nil error objects, it's possible to condense down to this:
``` go
    log.Fatal(http.ListenAndServe(":5000", nil))
```

## `http.ListenAndServeTLS`
``` go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```

# `encoding/json` package
Allows us to easily encode and decode our Go data types into JSON using the `Marshal` and `Unmarshal` functions.

## `json.Marshal`
``` go
json.Marshal(v interface{}) ([]byte, error)
```
Because the interface has zero methods, every Go type implements the empty interface, so we can effectively pass in any `struct` to this function, and it will marshal our struct into JSON.

``` go
type foo struct {
    Message string
    Age int
    Name string
    surname string
}

func main() {
    data, _ := json.Marshal(&foo{"4Score", 56, "Abe", "Lincoln"})
    fmt.Print(string(data))
}
```

``` json
{"Message":"4Score","Age":56,"Name":"Abe"}
```

The `surname` field isn't exported, therefore isn't included in the resulting JSON.

## `json.Unmarshal`
``` go
json.Unmarshal(data []byte, v interface{}) error
```

``` go
func main() {
    f := foo{}
    err := json.Unmarshal([]byte(`{"Message":"4Score","Age":56,"Name":"Abe"}`), &f)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(f.Message)
}
```

```
4Score
```

You can also customise the struct fields using tags so that the JSON field names don't have to match the Go struct names:
``` go
type foo struct {
    Message string  `json: "message,omitempty"`
    Age int         `json: "age,omitempty"`
    Name string     `json: "firstName,omitempty"`
    surname string
}
```

# Request

`Request.Method`
- string

`Request.Header`
- Header (map[string][]string)

`Request.Body`
- io.ReadCloser
- Returns EOF when not present

## `request.URL`

The `request.URL` is a struct that looks like this:
``` go
type URL struct {
    Scheme string
    Opaque string
    User *Userinfo
    Host string
    Path string
    RawPath string
    ForceQuery bool
    RawQuery string
    Fragment string
}
```

In a url `https://globomantics.com/api/products/123` the `Path` would be `/api/products/123`.

# Middleware

The request pipeline looks something like this:

```
Client -- Server -- HTTP Mux -- Handler
```

Middleware is basically functionality that is executed either before or after our intended handlers are called.

```
                                Authentication
Client -- Server -- HTTP Mux -- Logging            -- Handler
                                Session Management
```

One option we can implement this is by wrapping our handler in a special kind of adapter function.

The `http.HandlerFunc` type accepts a handler a parameter, and then returns another handler back to the caller. This allows us to execute code before or after the handler, and also gives us access to the request and response objects.

``` go
func middlewareHandler(handler http.Handler) http.Handler {
    return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
        // do stuff before intended handler here
        handler.ServeHttp(w, r)
        // do stuff after intended handler here
    })
}

func intendedFunction(w http.ResponseWriter, r *http.Request) {
    // business logic here
}

func main() {
    intendedHandler := http.HandlerFunc(intendedFunction)
    http.Handle("/foo", middlewareHandler(intendedFunction))
    http.ListenAndServe(":5000", nil)
}
```

# CORS (Cross Origin Resource Sharing)

CORS is designed to prevent cross-origin attacks.

It prevents websites which are served up from a specific origin from accessing resources from another origin.

```
Y  http://globomantics.com/products
Y  http://globomantics.com/api/products/123 
N  http://globomantics.com:8080/products        (due to different port)
N  https://globomantics.com/products            (due to different protocol)
N  http://dev.globomantics.com/dashboard        (due to subdomain)
```

To get around this, our web server needs to add special headers to response that get sent back to the client.

``` go
w.Header().Add("Access-Control-Allow-Origin", "*")
w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CSRF-Token, Accept-Encoding")
```
see more https://developer.mozilla.org/en-US/docs/Gloassary/CORS