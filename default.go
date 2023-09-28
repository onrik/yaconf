package yaconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func setDefaultValue(v reflect.Value) error {
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			setDefaultValue(v.Field(i))
			continue
		}

		if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
			if !v.Field(i).IsNil() {
				setDefaultValue(v.Field(i))
			}
			continue
		}

		tag := f.Tag.Get("yaconf")
		if tag == "" {
			continue
		}
		if strings.HasPrefix(tag, "default") {
			parts := strings.SplitN(tag, "=", 2)
			if len(parts) < 2 || parts[1] == "" {
				continue
			}
			defaultValue := parts[1]

			if f.Type.Kind() == reflect.Ptr {
				if f.Type.Elem().Kind() == reflect.String {
					v.Field(i).Set(reflect.ValueOf(&defaultValue))
					continue
				}
			}

			if f.Type.Kind() == reflect.String {
				v.Field(i).Set(reflect.ValueOf(defaultValue))
				continue
			}

			if f.Type.Kind() == reflect.Bool {
				value, err := strconv.ParseBool(defaultValue)
				if err != nil {
					return fmt.Errorf("%s is invalid value for bool", defaultValue)
				}
				v.Field(i).SetBool(value)
				continue
			}
			if isDuration(v.Field(i)) {
				value, err := time.ParseDuration(defaultValue)
				if err != nil {
					return fmt.Errorf("%s is invalid value for time.Duration", defaultValue)
				}
				v.Field(i).Set(reflect.ValueOf(value))
				continue
			}
			if isInt(v.Field(i)) {
				value, err := strconv.ParseInt(defaultValue, 10, 64)
				if err != nil {
					return fmt.Errorf("%s is invalid value for %s", defaultValue, v.Field(i).Type().Kind())
				}
				v.Field(i).SetInt(value)
				continue
			}
			if isUint(v.Field(i)) {
				value, err := strconv.ParseUint(defaultValue, 10, 64)
				if err != nil {
					return fmt.Errorf("%s is invalid value for %s", defaultValue, v.Field(i).Type().Kind())
				}
				v.Field(i).SetUint(value)
				continue
			}
		}
	}

	return nil
}

func isDuration(v reflect.Value) bool {
	if v.Type().Kind() != reflect.Int64 {
		return false
	}

	return v.Type() == reflect.TypeOf(time.Duration(0))
}

func isInt(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

func isUint(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func FillDefaultValues(config interface{}) error {
	v := reflect.ValueOf(config)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return setDefaultValue(v)
}
