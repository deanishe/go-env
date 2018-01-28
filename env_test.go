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
	"os"
	"testing"
	"time"
)

// Env is a string: string mapping.
type Env map[string]string

func withEnv(env Env, fn func(env Env)) {
	var (
		prev  = map[string]string{}
		unset = map[string]bool{}
	)

	for k, v := range env {
		if s, ok := os.LookupEnv(k); ok {
			prev[k] = s
		} else {
			unset[k] = true
		}

		// Update env
		os.Setenv(k, v)
	}

	// Ensure env is reset
	defer func() {

		for k, v := range prev {
			os.Setenv(k, v)
		}

		for k := range unset {
			os.Unsetenv(k)
		}

	}()

	// Call function
	fn(env)
}

func TestGet(t *testing.T) {
	env := Env{
		"key":  "value",
		"key2": "value2",
	}

	data := []struct {
		key string
		fb  []string
		out string
	}{
		{"key", []string{}, "value"},
		{"key", []string{"value2"}, "value"},
		{"key2", []string{}, "value2"},
		{"key2", []string{"value"}, "value2"},
		{"key3", []string{}, ""},
		{"key3", []string{"bob"}, "bob"},
	}

	withEnv(env, func(e Env) {

		// Verify env is the same
		for k, x := range env {
			v := Get(k)
			if v != x {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", k, x, v)
			}
		}

		// Test Get
		for _, td := range data {
			v := Get(td.key, td.fb...)
			if v != td.out {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", td.key, td.out, v)
			}

		}
	})
}

// Basic usage of Get. Returns an empty string if variable is unset.
func ExampleGet() {
	// Set some test variables
	os.Setenv("test_name", "Bob Smith")
	os.Setenv("test_address", "7, Dreary Lane")

	fmt.Println(Get("test_name"))
	fmt.Println(Get("test_address"))
	fmt.Println(Get("test_nonexistent")) // unset variable

	// Output:
	// Bob Smith
	// 7, Dreary Lane
	//

	os.Unsetenv("test_name")
	os.Unsetenv("test_address")
}

// The fallback value is returned if the variable is unset.
func ExampleGet_fallback() {
	// Set some test variables
	os.Setenv("test_name", "Bob Smith")
	os.Setenv("test_address", "7, Dreary Lane")

	fmt.Println(Get("test_name", "default name"))       // fallback ignored
	fmt.Println(Get("test_address", "default address")) // fallback ignored
	fmt.Println(Get("test_nonexistent", "hi there!"))   // unset variable

	// Output:
	// Bob Smith
	// 7, Dreary Lane
	// hi there!

	os.Unsetenv("test_name")
	os.Unsetenv("test_address")
}

func TestGetInt(t *testing.T) {
	env := Env{
		"one":   "1",
		"two":   "2",
		"zero":  "0",
		"float": "3.5",
		"word":  "henry",
	}

	data := []struct {
		key string
		fb  []int
		out int
	}{
		{"one", []int{}, 1},
		{"two", []int{1}, 2},
		{"five", []int{5}, 5},
		{"zero", []int{}, 0},
		{"word", []int{}, 0},
		{"word", []int{5}, 5},
		{"float", []int{}, 3},
		{"float", []int{5}, 3},
	}

	withEnv(env, func(e Env) {

		// Test GetInt
		for _, td := range data {
			v := GetInt(td.key, td.fb...)
			if v != td.out {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", td.key, td.out, v)
			}

		}
	})
}

func TestGetFloat(t *testing.T) {
	env := Env{
		"one.three": "1.3",
		"two":       "2.0",
		"zero":      "0",
		"word":      "henry",
	}

	data := []struct {
		key string
		fb  []float64
		out float64
	}{
		{"one.three", []float64{}, 1.3},
		{"two", []float64{1}, 2.0},
		{"five", []float64{5.0}, 5.0},
		{"zero", []float64{}, 0.0},
		{"word", []float64{}, 0.0},
		{"word", []float64{5.0}, 5.0},
	}

	withEnv(env, func(e Env) {

		// Test GetFloat
		for _, td := range data {
			v := GetFloat(td.key, td.fb...)
			if v != td.out {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", td.key, td.out, v)
			}

		}
	})
}

func TestGetDuration(t *testing.T) {
	env := Env{
		"5mins": "5m",
		"1hour": "1h",
		"zero":  "0",
		"word":  "henry",
	}

	data := []struct {
		key string
		fb  []time.Duration
		out time.Duration
	}{
		{"5mins", []time.Duration{}, time.Minute * 5},
		{"1hour", []time.Duration{time.Second * 1}, time.Hour * 1},
		{"zero", []time.Duration{time.Second * 2}, 0},
		{"zero", []time.Duration{}, 0},
		{"word", []time.Duration{}, 0},
		{"word", []time.Duration{time.Second * 5}, time.Second * 5},
	}

	withEnv(env, func(e Env) {

		// Test GetDuration
		for _, td := range data {
			v := GetDuration(td.key, td.fb...)
			if v != td.out {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", td.key, td.out, v)
			}

		}
	})
}

func TestGetBool(t *testing.T) {
	env := Env{
		"empty": "",
		"t":     "t",
		"f":     "f",
		"1":     "1",
		"0":     "0",
		"true":  "true",
		"false": "false",
		"word":  "nonsense",
	}

	data := []struct {
		key string
		fb  []bool
		out bool
	}{
		{"empty", []bool{}, false},
		{"empty", []bool{true}, true},
		{"t", []bool{}, true},
		{"f", []bool{}, false},
		{"1", []bool{}, true},
		{"0", []bool{}, false},
		{"true", []bool{}, true},
		{"false", []bool{}, false},
		{"word", []bool{}, false},
		{"word", []bool{true}, true},
	}

	withEnv(env, func(e Env) {

		// Test GetBool
		for _, td := range data {
			v := GetBool(td.key, td.fb...)
			if v != td.out {
				t.Errorf("Bad '%s'. Expected=%v, Got=%v", td.key, td.out, v)
			}

		}
	})
}
