package yaconf

import (
	"fmt"
	"reflect"
	"strings"
)

type validator interface {
	Validate() error
}

func addPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", prefix, name)
}

func validate(config interface{}, prefix string) []string {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	errors := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		name := f.Tag.Get("yaml")
		if name == "" {
			name = f.Name
		}

		if strings.Contains(f.Tag.Get("yaconf"), "required") && v.Field(i).IsZero() {
			errors = append(errors, fmt.Sprintf("%s is required", addPrefix(prefix, name)))
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			errors = append(errors, validate(v.Field(i).Interface(), addPrefix(prefix, name))...)
			continue
		}

		if f.Type.Kind() == reflect.Ptr {
			if f.Type.Elem().Kind() != reflect.Struct {
				continue
			}
			if v.Field(i).IsNil() {
				errors = append(errors, validate(reflect.New(v.Field(i).Type().Elem()).Interface(), addPrefix(prefix, name))...)
			} else {
				errors = append(errors, validate(v.Field(i).Elem().Interface(), addPrefix(prefix, name))...)
			}
		}

	}
	return errors
}
