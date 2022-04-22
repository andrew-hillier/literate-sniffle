## Persisting Data

### `sql.Open`

``` go
func Open(driverName, dataSourceName string) (*DB, error)
```

DB Type
- Configurable pool of zero or more connections
- Creates and frees connections automatically
- Thread-safe

``` go
import "database/sql"
...

var DbConn *sql.DB

func SetupDatabase() {
    var err error
    DbConn, err = sql.Open("mysql", "root:password123@tcp(127.0.0.1:3306)/inventorydb")
    if err != nil {
        log.Fatal(err)
    }
}
```

The sql drivers aren't included in the Go standard library, so we need to add a package for that, in the case of MySql, this is the command to add the package:
```
go get -u github.com/go-sql-driver/mysql
```
Drivers: https://github.com/golang/go/wiki/SQLDrivers

We then need to add this import as an "indirect" reference in our code
``` go
import (
    _ "github.com/go-sql-driver/mysql"
)
```

### DB.Query
``` go
func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
```

Rows Type
- Result of a query
- Use "Next" to advance
- Needs to be closed

#### Rows.Scan
Scan converts the columns read from the database into our Go types.
``` go
func (rs *Rows) Scan(dest ...interface{}) error
```

``` go
...
results, err := db.Query(`select productId, manufacturer, sku from products`)
if err != nil {
    log.Fatal(err)
}
defer results.Close()
products := make([]Product, 0)
for results.Next() {
    var product Product
    results.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, ...)
    products = append(products, product)
}
```

### DB.QueryRow

This will only return the top1 result, and discard the rest.

``` go
func (db *DB) QueryRow(query string, args ...interface{}) *Row
```

#### Row.Scan

``` go
func (rs *Row) Scan(dest ...interface{}) error
```

### DB.Exec

``` go
func (rs *DB) Exec(query string, args ...interface{}) (Result, error)
```

#### sql.Result

``` go
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
```

``` go
...
result, err := db.Exec(`update products set sku=? where productid=?`,
    product.Sku,
    product.ProductID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("number of affected rows %d\n", result.RowsAffected())
...
```

### Connection Pooling
Connection Max Lifetime
 - Sets the maximum amount of time a connection may be used.
 - If set to 0, the connections will be infinitely reused.

Max Idle Connections
 - Sets the maximum number of connections in the idle connection pool.
 - Default is 2

Max Open Connections
 - Sets the maximum number of open connections to the database.
 - If you set this value to a number that is lower that the max idle connections, it will automatically lower the max idle connections to match the max number of open connections.

If a new connection comes in and and you've already reached your maximum number of connections, the new request will now block.

#### Context
Allows you to set a deadline, cancel a signal, or set other reqest-scoped values across API boundries and between processes.

``` go
...
ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
results, err := db.QueryContext(ctx, `select productId, manufacturer, sku from products`)
if err != nil {
    log.Fatal(err)
}
defer results.Close()
products := make([]Product, 0)
for results.Next() {
    results.Scan(&product.ProductID, &product.Manufacturer, &product.Sku ...)
    products = append(products, product)
}
```

Similar functions are avaliable to supply context: 
- QueryContext
- QueryRowContext
- ExecContext

## File Upload
- base64 encode
  - Convert the file to a string and include the JSON payload
- multipart/form-data (more efficient)
  - Uses an HTTP form to submit the raw data

### `Encoding.DecodeString`
``` go
func (enc *Encoding) DecodeString(s string) ([]byte, error)
```

``` go
str := "SGVsbG8gV29ybGQ="
output, err := base64.StdEncoding.DecodeString(str)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%q\n", output)
```

### `Request.FormFile`
``` go
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
```

#### `multipart.File`
``` go
type File interface {
    io.Reader,
    io.ReaderAt,
    io.Seeker,
    io.Closer
}
```

#### `multipart.FileHeader`
``` go
type FileHeader struct {
    Filename string
    Header textproto.MIMEHeader
    Size int64
}
```

``` go
r.ParseMultiPartForm(5 << 20) // 5 Mb
file, handler, err := r.FormFile("uploadFileName")
if err != nil {
    fmt.Println("error reading file from request")
    return
}
defer file.Close()
f, err := os.OpenFile("./filepath/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
defer f.Close()
io.Copy(f, file)
```

``` go
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
    filename = "gopher.png"
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("error reading file")
        return
    }
    defer file.Close()
    w.Header.Set("Content-Disposition", "attachment; filename=" + filename)
    io.Copy(w, file)
}
```