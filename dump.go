// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package env

import (
	"encoding"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// sentinel error returned by toString to indicate that Dump should try
// further methods to convert a value to a string.
// nolint: stylecheck
var unknownType = errors.New("unknown type")

// DumpOption is a configuration option to Dump.
type DumpOption func(d *dumper)

// IgnoreZeroValues excludes zero values from the returned map of variables.
// Non-nil slices are unaffected by the setting: an empty string is returned
// for empty slices regardless.
func IgnoreZeroValues(d *dumper) { d.noZero = true }

// VarNameFunc specifies a different function to generate the names of the
// variables returned by Dump.
func VarNameFunc(fun func(string) string) DumpOption {
	return func(d *dumper) {
		d.nameFunc = fun
	}
}

// Dump extracts a struct's fields to a map of variables.
// By default, the names (map keys) of the variables are generated using
// VarName. Pass the VarNameFunc option to generate custom keys.
func Dump(v interface{}, opt ...DumpOption) (map[string]string, error) {
	d := &dumper{
		nameFunc: VarName,
	}
	for _, o := range opt {
		o(d)
	}
	return d.dump(v)
}

// Export extracts a struct's fields' values (via Dump) and exports them to the
// environment (via os.Setenv). It accepts the same options as Dump.
func Export(v interface{}, opt ...DumpOption) error {
	vars, err := Dump(v, opt...)
	if err != nil {
		return err
	}
	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// dumper reads a struct's fields and returns them as a map[string]string.
type dumper struct {
	noZero   bool
	nameFunc func(string) string
}

func (d *dumper) dump(v interface{}) (map[string]string, error) {
	vars := map[string]string{}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	rvType := rv.Type()

	for i := 0; i < rvType.NumField(); i++ {
		var (
			val   = rv.Field(i)
			field = rvType.Field(i)
			name  = field.Name
			key   = field.Tag.Get("env")
		)

		if d.noZero && isZero(val) {
			continue
		}

		// skip unexported fields
		if string(name[0]) == strings.ToLower(string(name[0])) || key == "-" {
			continue
		}

		if key == "" {
			key = d.nameFunc(name)
		}

		if val.Kind() == reflect.Ptr && val.IsNil() {
			vars[key] = ""
			continue
		}

		if val.Kind() == reflect.Slice {
			s, err := dumpSlice(val)
			if err != nil {
				return nil, err
			}
			if s == "" && d.noZero {
				continue
			}
			vars[key] = s
			continue
		}

		s, err := toString(val)
		if err != nil && err != unknownType {
			return nil, err
		}
		if err != unknownType {
			vars[key] = s
			continue
		}

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() == reflect.Struct {
			m, err := d.dump(val.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range m {
				vars[k] = v
			}
			continue
		}
	}

	return vars, nil
}

func dumpSlice(rv reflect.Value) (string, error) {
	var values []string
	for i := 0; i < rv.Len(); i++ {
		v := rv.Index(i)
		s, err := toString(v)
		if err != nil && err != unknownType {
			return "", err
		}
		values = append(values, s)
		continue
	}
	return strings.Join(values, ","), nil
}

func toString(rv reflect.Value) (value string, err error) {
	if tm, ok := rv.Interface().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if s, ok := rv.Interface().(fmt.Stringer); ok {
		return s.String(), nil
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
	}
	return "", unknownType
}

// copied from go1.13 source code, so library also runs on 1.12
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(&reflect.ValueError{"isZero", v.Kind()})
	}
}
