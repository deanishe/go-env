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
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Bind populates the fields of a struct from environment variables.
//
// Variables are mapped to fields using `env:"..."` tags, and the
// struct is populated by passing it to Bind(). Unset or empty
// environment variables are ignored.
//
// Untagged fields have a default environment variable assigned to
// them. See VarName() for details of how names are generated.
//
// Bind accepts an optional Env argument. If provided, values will
// be looked up via that Env instead of the program's environment.
func Bind(v interface{}, env ...Env) error {

	e := sysEnv // default env
	if len(env) > 0 {
		e = &envReader{env[0]}
	}

	binds, err := extract(v, e)
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

// binding links an environment variable to the field of a struct.
type binding struct {
	Name   string
	EnvVar string
	Target interface{}
	Kind   reflect.Kind
	env    *envReader
}

// Load populates the target struct from the environment.
func (bind *binding) Load() error {

	rv := reflect.Indirect(reflect.ValueOf(bind.Target))
	typ := rv.Type()

	for i := 0; i < rv.NumField(); i++ {

		field := typ.Field(i)
		value := rv.Field(i)

		if field.Name == bind.Name {
			// Ignore empty/unset fields
			if bind.env.Get(bind.EnvVar) == "" {
				return nil
			}

			return bind.setField(&field, &value)
		}

	}
	return fmt.Errorf("unknown field: %s", bind.Name)
}

func (bind *binding) setField(field *reflect.StructField, rv *reflect.Value) error {

	switch bind.Kind {

	case reflect.Bool:
		b := bind.env.GetBool(bind.EnvVar)
		reflect.Indirect(*rv).SetBool(b)
		// log.Printf("[%s] value=%v", bind.Name, b)

	case reflect.String:

		s := bind.env.GetString(bind.EnvVar)
		reflect.Indirect(*rv).SetString(s)
		// log.Printf("[%s] value=%s", bind.Name, s)

	// Special-case int64, as it may also be a duration.
	case reflect.Int64:

		// Try to parse value as an int, and if that fails, try
		// to parse it as a duration.
		s := bind.env.GetString(bind.EnvVar)

		if _, err := strconv.ParseInt(s, 10, 64); err == nil {

			i := bind.env.GetInt(bind.EnvVar)
			reflect.Indirect(*rv).SetInt(int64(i))

		} else {

			if d, err := time.ParseDuration(s); err == nil {
				reflect.Indirect(*rv).SetInt(int64(d))
			}

		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:

		i := bind.env.GetInt(bind.EnvVar)
		reflect.Indirect(*rv).SetInt(int64(i))
		// log.Printf("[%s] value=%d", bind.Name, i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		i := bind.env.GetInt(bind.EnvVar)
		reflect.Indirect(*rv).SetUint(uint64(i))
		// log.Printf("[%s] value=%d", bind.Name, i)

	case reflect.Float32, reflect.Float64:

		n := bind.env.GetFloat(bind.EnvVar)
		reflect.Indirect(*rv).SetFloat(n)
		// log.Printf("[%s] value=%f", bind.Name, n)

	default:
		return fmt.Errorf("unsupported type: %v", bind.Kind.String())
	}

	return nil

}

func extract(v interface{}, env *envReader) ([]*binding, error) {

	var binds []*binding

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
			varname = VarName(name)
		}

		bind := &binding{
			Name:   name,
			EnvVar: varname,
			Target: v,
			Kind:   field.Type.Kind(),
			env:    env,
		}
		binds = append(binds, bind)

	}

	return binds, nil
}

// VarName generates an environment variable name from a field name.
// This is documented to show how the automatic names are generated.
func VarName(name string) string {
	if !isCamelCase(name) {
		return strings.ToUpper(name)
	}
	return splitCamelCase(name)
}

func isCamelCase(s string) bool {
	if ok, _ := regexp.MatchString(".*[a-z]+[0-9_]*[A-Z]+.*", s); ok {
		return true
	}
	if ok, _ := regexp.MatchString("[A-Z][A-Z][A-Z]+[0-9_]*[a-z]+.*", s); ok {
		return true
	}
	return false
}

func splitCamelCase(name string) string {

	var (
		i     int
		re    *regexp.Regexp
		rest  string
		words []string
	)

	rest = name

	// Start with 3 or more capital letters.
	re = regexp.MustCompile("[A-Z]+([A-Z])[0-9]*[a-z]")
	for {
		if idx := re.FindStringSubmatchIndex(rest); idx != nil {
			i = idx[2]
			s := rest[:i]
			rest = rest[i:]
			words = append(words, s)
		} else {
			break
		}
	}

	re = regexp.MustCompile("[a-z][0-9_]*([A-Z])")
	for {

		if idx := re.FindStringSubmatchIndex(rest); idx != nil {
			i = idx[2]
			s := rest[idx[2]:idx[3]]

			s = rest[:i]
			rest = rest[i:]
			words = append(words, s)
		} else {
			break
		}
	}

	if rest != "" {
		words = append(words, rest)
	}

	if len(words) > 0 {
		s := strings.ToUpper(strings.Join(words, "_"))
		return s
	}

	return ""
}
