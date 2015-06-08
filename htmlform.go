// Package htmlform is a very simple struct to HTML form mapper

package htmlform

import (
	"fmt"
	"reflect"
	"strings"
)

// Form is a HTML form
type Form map[string]interface{}

// FormField is a HTML form field
type FormField map[string]interface{}

// CSRF is a HTML form field containing a CSRF token
type CSRF struct {
	FieldName string
	Token     string
}

// Name returns the form fields name
func (c CSRF) Name() string { return c.FieldName }

// Value returns the field value, formatted with fmt.Sprintf("%s")
func (c CSRF) Value() string { return c.Token }

// Errors is a function which, given a field name, returns a FieldError
type Errors func(field string) map[string]interface{}

// NameStructTag sets the struct tag used to set the field name, e.g. schema
// If set to an empty string, the default name will be used.
var NameStructTag = "schema"

// HtmlformStructTag sets the struct tag used to set other form properties, e.g.
// setting the form field type to "password".
//
// If the name= paramter is set, the value inferred or given by StructTag is replaced.
var HtmlformStructTag = "htmlform"

// Create returns a Form
func Create(model interface{}, errs Errors, namespace []string, htmlNamespace []string) Form {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		if v.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("htmlform: interface must be a pointer to struct, got %s", v.Kind()))
		}
		panic(fmt.Sprintf("htmlform: interface must be a pointer to struct, got pointer to %s", v.Elem().Kind()))
	}

	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	form := Form{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		t := v.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			if v.IsNil() {
				v.Set(reflect.New(t))
			}
			v = v.Elem()
		}

		fieldType := "hidden"

		switch f.Type.Kind() {
		case reflect.Map:
			// TODO: ?
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			// TODO: ?
			fieldType = "array"
		case reflect.String:
			fieldType = "text"
		case reflect.Int:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fieldType = "number"
		case reflect.Struct:
			fieldType = "struct"
		}

		n := f.Name
		if len(NameStructTag) > 0 {
			if v := f.Tag.Get(NameStructTag); len(v) > 0 {
				n = v
			}
		}

		args := make(map[string]string)
		if len(HtmlformStructTag) > 0 {
			if v := f.Tag.Get(HtmlformStructTag); len(v) > 0 {
				p := strings.Split(v, ",")
				for _, p2 := range p {
					p3 := strings.SplitN(p2, "=", 2)
					if len(p3) == 1 {
						args[p3[0]] = "1"
					} else {
						args[p3[0]] = p3[1]
					}
				}
			}
		}

		if fieldType == "struct" {
			form[f.Name] = Create(v.Field(i).Addr().Interface(), errs, append(namespace, f.Name), append(htmlNamespace, n))
			continue
		}

		if ft, ok := args["type"]; ok {
			fieldType = ft
		}
		if fn, ok := args["name"]; ok {
			n = fn
		}

		var er map[string]interface{}
		if errs != nil {
			er = errs(strings.Join(append(namespace, f.Name), "."))
		}
		form[f.Name] = FormField{
			"Name":          strings.Join(append(htmlNamespace, n), "."),
			"Type":          fieldType,
			"Value":         fmt.Sprintf("%s", v.Field(i).Interface()),
			"Errors":        er,
			"namespace":     namespace,
			"htmlNamespace": htmlNamespace,
		}
	}

	return form
}

// WithCSRF returns the form with a CSRF token
func (f Form) WithCSRF(field, token string) Form {
	f["_CSRF"] = &CSRF{
		FieldName: field,
		Token:     token,
	}
	return f
}
