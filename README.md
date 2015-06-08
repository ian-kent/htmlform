htmlform [![GoDoc](https://godoc.org/github.com/ian-kent/htmlform?status.svg)](https://godoc.org/github.com/ian-kent/htmlform)
========

**Experimental**: Not suitable for production use.

htmlform simplifies the interaction between web handlers, data models, decoding, validation and template logic.

### Example

Example struct:
```go
// using gorilla/schema and bluesuncorp/validator.v5

type registrationForm struct {
    Email string `validate:"required,email" schema:"email"`
    Password string `validate:"required,min=10" schema:"password" htmlform:"type=password"`
}
```

Call htmlform:
```go
var model registrationForm
form := htmlform.Create(&model, nil, []string{}, []string{}).WithCSRF(nosurf.FormFieldName, nosurf.Token(req))
```

Returns this:
```go
map[string]interface{}{
  "Email": FormField{
    "Name":          "email",
    "Type":          "text",
    "Value":         "",
    "Errors":        map[string]interface{}{},
  },
  "_CSRF": CSRF{
    "FieldName": "csrf_token",
    "Token":     "generated-csrf-token",
  },
}
```

Which you can use in a template:
```go
{{ define "text_field" }}
    {{ if .Errors }}
        <strong>Error</strong>
    {{ end }}
    <input type="{{ .Type }}" name="{{ .Name }}" value="{{ .Value }}">
{{ end }}

{{ define "hidden_field" }}
    <input type="hidden" name="{{ .Name }}" value="{{ .Value }}">
{{ end }}

{{ template "text_field" .Form.Email }}
{{ template "hidden_field" .Form._CSRF }}
```

### Errors

The second argument to `htmlform.Create` is a `func(string) map[string]interface{}`.

If set, it's called for each field to retrieve errors, making it easy to map
errors from form validation into a structure useful for the template.

Example output from validator.v5:

```go
errs := map[string]interface{}{
    "email": map[string]interface{
        "min": map[string]interface{
            "Param": "10",
        },
    },
}
errFunc := func(field string) map[string]interface{} {
  return errs[field]
}
```

Pass errFunc to `htmlform.Create` as the second argument.

From a template:

```go
{{ define "form_error" }}
    {{ if .Errors }}
        {{ if .Errors.min }}{{ .Name }} must be at least {{ .Errors.min.Param }} characters{{ end }}
    {{ end }}
{{ end }}

{{ template "form_error" .Form.Email }}
```

### Licence

Copyright ©‎ 2015, Ian Kent (http://iankent.uk).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
