# Using Templating

Simply put, templates in Go are specifically formatted text, which are designed to interact with the data structure to produce formatted ouput.

Data Structure:
``` go
type BlogPost struct {
    Header string
    Message string
}
```

Template:
``` html
<html>
    ...
    <h1>{{.Header}}</h1>
    <div>
        <p>{{.Message}}</p>
    </div>
</html>
```

There are two packages that deal with templating in Go;
- `text/template` base functionality for qorking with templates in Go
- `html/template` same interface but with added security for HTML output

`template.New`
``` go
func New(name string) *Template
```

`Template.Parse`
``` go 
func (t *Template) Parse(text string) (*Template, error)
```

`Template.Execute`
``` go
func (t *Template) Execute(wr io.Writer, data interface{}) error
```

``` go
import "html/template"

type BlogPost struct {
    Header string
    Message string
}

func main() {
    post := BlogPost{"First Post!", "Hello World"}
    tmpl, _ := template.New("post").Parse(`<h1>{{.Header}}</h1><p>{{.Message}}</p>`)
    tmpl.Execute(os.Stdout, post)
}
```

## Pipelines

Pipelines are a sequence of commands that are able to be chained together to produce some kind of an output.

A command can be a simple value or an argument, like we see above by referencing the fields in our structs, or commands can be functions or method calls, that allow us to pass in one or more arguments.

Examples:
``` go
{{ "Hello" }}
{{ 1234 }}
{{ .Message }}
{{ println "Hi" }}
{{ .SayHello }}
{{ .SaySomething "Bye" }}
```

Pipelines can also be chained, for example:
``` go
{{ .SaySomething "Hello" }}     // not chained
{{ "Hello" | .SaySomething }}   // equivilant of above using pipeline chaining
{{ "Hello" | .SaySomething | printf "%s %s" "World" }}
```

## Pipeline Looping

``` go
{{ range pipeline }} T1 {{ end }}                   // standard
{{ range pipeline }} T1 {{ else }} T2 {{ end }}     // T2 = if no results
{{ range $index, $element := pipeline }}            // provides index of the current element
```

Example:
``` go
import "html/template"

tmpl := "{{range .}}{{.}}{{end}}"

func main() {
    items := []string{"one", "two", "three"}
    tmpl, _ := template.New("tmplt").Parse(tmpl)
    err := tmpl.Execute(os.Stdout, items)
}
```

Output:
```
onetwothree
```

## Template Functions

| Function | Example |
|-|-|
| and | `{{if and true true true}} {{end}}` |
| or | `{{if or true false true}} {{end}}` |
| index | `{{index .1}}` |
| len | `{{len .}}` |
| not | `{{if not false}}` |
| print, printf, println | `{{println "hey"}}` |

https://golang.org/pkg/text/template#hdr-Functions

## Template Operators

| Operator | Example |
|-|-|
| eq | `arg1 == arg2` |
| ne | `arg1 != arg2` |
| lt | `arg1 < arg2` |
| le | `arg1 <= arg2` |
| gt | `arg1 > arg2` |
| ge | `arg1 >= arg2` |

## Custom Functions

### `Template.Funcs`

We can use our own Go functions in templates, by calling the `Funcs` method on the `Template`. This is a map of strings to functions, where if the map key is encountered when executing a template, the function will be called.
``` go
func (t *Template) Funcs(funcMap FuncMap) *Template

type FuncMap map[string]interface{}
```

Some rules for template functions though;
- Return a single value
- Return a single value, or an error

Example:
```go
import "html/template"

tmpl := "{{range $index, $element := .}}{{if mod $index 2}}{{.}}{{end}}{{end}}"

func main() {
    items := []string{"one", "two", "three"}
    fm := template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}
    tmpl, _ := template.New("tmplt").Funcs(fm).Parse(tmpl)
    err := tmpl.Execute(os.Stdout, items)
}
```

Output:
```
onethree
```