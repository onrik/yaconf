package yaconf

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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
		name := f.Tag.Get("yaml")
		if name == "" {
			name = f.Name
		}

		if f.Tag.Get("yaconf") == "required" && v.Field(i).IsZero() {
			errors = append(errors, fmt.Sprintf("%s is required", addPrefix(prefix, name)))
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			errors = append(errors, validate(v.Field(i).Interface(), addPrefix(prefix, name))...)
			continue
		}

		if f.Type.Kind() == reflect.Ptr {
			if f.Type.Elem().Kind() == reflect.Struct {
				if v.Field(i).IsNil() {
					v.Field(i).Set(reflect.New(f.Type.Elem()))
				}
				errors = append(errors, validate(v.Field(i).Elem().Interface(), addPrefix(prefix, name))...)
			}
		}

	}
	return errors
}

// Read config from file
func Read(filename string, config interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	errors := validate(config, "")
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	if v, ok := config.(validator); ok {
		return v.Validate()
	}

	return nil
}
