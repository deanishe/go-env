//
// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2018-01-27
//

package env

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Bind populates the fields of a struct from environment variables.
// Variables are mapped to fields using `env:"..."` tags.
func Bind(v interface{}) error {

	binds, err := extract(v)
	if err != nil {
		return err
	}

	for _, bind := range binds {
		if err := bind.Load(); err != nil {
			return err
		}
	}

	return nil
}

// Binding links an environment variable to the field of a struct.
type Binding struct {
	Name       string
	EnvVar     string
	FieldValue reflect.Value
	Field      reflect.Type
	Target     interface{}
	Kind       reflect.Kind
}

// Load populates the target struct from the environment.
func (bind *Binding) Load() error {

	rv := reflect.Indirect(reflect.ValueOf(bind.Target))
	typ := rv.Type()

	for i := 0; i < rv.NumField(); i++ {

		field := typ.Field(i)
		value := rv.Field(i)

		if field.Name == bind.Name && Get(bind.EnvVar) != "" {
			return bind.setField(&field, &value)
		}

	}
	return nil
}

func (bind *Binding) setField(field *reflect.StructField, rv *reflect.Value) error {

	switch bind.Kind {

	case reflect.Bool:
		b := GetBool(bind.EnvVar)
		reflect.Indirect(*rv).SetBool(b)
		// log.Printf("[%s] value=%v", bind.Name, b)

	case reflect.String:

		s := GetString(bind.EnvVar)
		reflect.Indirect(*rv).SetString(s)
		// log.Printf("[%s] value=%s", bind.Name, s)

	// Special-case int64, as it may also be a duration.
	case reflect.Int64:

		// Try to parse value as an int, and if that fails, try
		// to parse it as a duration.
		s := GetString(bind.EnvVar)

		if _, err := strconv.ParseInt(s, 10, 64); err == nil {

			i := GetInt(bind.EnvVar)
			reflect.Indirect(*rv).SetInt(int64(i))

		} else {

			if d, err := time.ParseDuration(s); err == nil {
				reflect.Indirect(*rv).SetInt(int64(d))
			}

		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:

		i := GetInt(bind.EnvVar)
		reflect.Indirect(*rv).SetInt(int64(i))
		// log.Printf("[%s] value=%d", bind.Name, i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		i := GetInt(bind.EnvVar)
		reflect.Indirect(*rv).SetUint(uint64(i))
		// log.Printf("[%s] value=%d", bind.Name, i)

	case reflect.Float32, reflect.Float64:

		n := GetFloat(bind.EnvVar)
		reflect.Indirect(*rv).SetFloat(n)
		// log.Printf("[%s] value=%f", bind.Name, n)

	default:
		return fmt.Errorf("unsupported type: %v", bind.Kind.String())
	}

	return nil

}

func extract(v interface{}) ([]*Binding, error) {

	var binds []*Binding

	ref := reflect.ValueOf(v)

	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	if ref.Kind() != reflect.Struct {
		return nil, fmt.Errorf("need struct, not %s: %v", ref.Kind(), v)
	}

	typ := ref.Type()

	for i := 0; i < ref.NumField(); i++ {

		var (
			ok      bool
			field   reflect.StructField
			name    string
			tag     string
			varname string
		)

		field = typ.Field(i)
		name = field.Name

		if tag, ok = field.Tag.Lookup("env"); ok {
			if tag == "-" { // Ignore this field
				continue
			}
		}

		if tag != "" {
			varname = tag
		} else {
			varname = envName(name)
		}

		bind := &Binding{
			Name:   name,
			EnvVar: varname,
			Target: v,
			Kind:   field.Type.Kind(),
		}
		binds = append(binds, bind)

	}

	return binds, nil
}

func envName(name string) string {
	words := splitName(name)
	for i, s := range words {
		words[i] = strings.ToUpper(s)
	}
	return strings.Join(words, "_")
}

func splitName(name string) []string {

	var (
		consumeUpper bool
		s            string
		words        []string
	)

	for _, r := range name {

		if unicode.IsUpper(r) {

			if consumeUpper {
				s += string(r)
			} else if s != "" {
				words = append(words, strings.ToLower(s))
				s = string(r)
			} else {
				s += string(r)
			}

			consumeUpper = true
			continue
		}

		consumeUpper = false
		s += string(r)
	}

	if s != "" {
		words = append(words, strings.ToLower(s))
	}

	return words
}
