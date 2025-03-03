package goconf

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadEnv(target any) error {
	// Check if target is a non-nil pointer
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	return loadStructConfig(v.Elem())
}

func loadStructConfig(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct but got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envVar := field.Tag.Get("env")

		if envVar == "" {
			continue
		}

		// Get the environment variable value
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return fmt.Errorf("environment variable %s is not set", envVar)
		}

		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}

		// Handle the field based on its kind
		switch fieldValue.Kind() {

		case reflect.String:
			fieldValue.SetString(envValue)

		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return fmt.Errorf("failed to parse %s as int: %v", envVar, err)
			}

			fieldValue.SetInt(int64(intValue))

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("failed to parse %s as bool: %v", envVar, err)
			}

			fieldValue.SetBool(boolValue)

		case reflect.Struct:
			// If the field is a struct, recurse into it
			if fieldValue.CanAddr() {
				if err := loadStructConfig(fieldValue); err != nil {
					return fmt.Errorf("failed to load nested struct field %s: %v", field.Name, err)
				}
			}

		case reflect.Ptr:
			// If the field is a pointer to a struct, dereference it
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}

			if fieldValue.Elem().Kind() == reflect.Struct {
				if err := loadStructConfig(fieldValue.Elem()); err != nil {
					return fmt.Errorf("failed to load nested struct field %s: %v", field.Name, err)
				}
			}

		default:
			return fmt.Errorf("unsupported type %s for field %s", fieldValue.Kind(), field.Name)
		}
	}

	return nil
}
