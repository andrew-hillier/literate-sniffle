## `websocket.Conn`

``` go
type Conn struct {
    PayloadType byte
    MaxPayloadBytes int
}
```

``` go
import (
    "net/http"
    "net/http"
)

func socket(ws *websocket.Conn) {
    // Handle receiving and sending data
}

func main() {
    http.Handle("/websocket", websocket.Handler(socket))
    log.Fatal(http.ListenAndServe(":5000", nil))
}
```

`Conn` has methods for reading and writing...
``` go
func (ws *Conn) Read(msg []byte) (n int, err error)
func (ws *Conn) Write(msg []byte) (n int, err error)
```

...however Go provides a simipler method for us to, in the form of this `Codec` type:
``` go
type Codec struct {
    Marshal func(v interface{}) (data []byte, payloadType byte, err error)
    Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
}
```

There are two "convenience" objects which implement this Codec type in the WebSocket package;
- `message` for when you want to send and receive byte or text data
- `JSON` for when you want to send JSON back and forth

### `code.Receive`
``` go
func (cd Codec) Receive(ws *Conn, v interface{}) (err error)
```

``` go
func socket (we *websocket.Conn) {
    go func(c *websocket.Conn) {
        for {
            var msg message
            if err := websocket.JSON.Receive(c, &msg); err != nil {
                break
            }
            fmt.Printf("received message %s\n", msg.Data)
        }
    }
}
```

### `code.Send`
``` go
func (cd Codec) Send(ws *Conn, v interface{}) (err error)
```

``` go
func socket (we *websocket.Conn) {
    products, _ := product.GetTopTenProducts()
    for {
        time.Sleep(10 * time.Second)
        if err := websocket.JSON.Send(ws, products); err != nil {
            break
        }
    }
}
```
## Testing

You can establish a connection in the browser console, and then use the Developer Tools > Network > `websocket` > Messages to see the data sent from the server:
``` js
let ws = new WebSocket("ws://localhost:5000/websocket")
ws.send(JSON.stringify({data: "test message from browser", type: "test"}))
```