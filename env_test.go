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
		"empty": "",
	}

	data := []struct {
		key string
		fb  []int
		out int
	}{
		// numbers
		{"one", []int{}, 1},
		{"two", []int{1}, 2},
		{"zero", []int{}, 0},
		{"zero", []int{2}, 0},
		// empty values
		{"empty", []int{}, 0},
		{"empty", []int{5}, 5},
		// non-existent values
		{"five", []int{}, 0},
		{"five", []int{5}, 5},
		// invalid values
		{"word", []int{}, 0},
		{"word", []int{5}, 5},
		// floats
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

// Getting int values with and without fallbacks.
func ExampleGetInt() {
	// Set some test variables
	os.Setenv("PORT", "3000")
	os.Setenv("PING_INTERVAL", "")

	fmt.Println(GetInt("PORT"))
	fmt.Println(GetInt("PORT", 5000))        // fallback is ignored
	fmt.Println(GetInt("PING_INTERVAL"))     // returns zero value
	fmt.Println(GetInt("PING_INTERVAL", 60)) // returns fallback
	// Output:
	// 3000
	// 3000
	// 0
	// 60

	os.Unsetenv("PORT")
	os.Unsetenv("PING_INTERVAL")
}

func TestGetFloat(t *testing.T) {
	env := Env{
		"one.three": "1.3",
		"two":       "2.0",
		"zero":      "0",
		"empty":     "",
		"word":      "henry",
	}

	data := []struct {
		key string
		fb  []float64
		out float64
	}{
		// numbers
		{"one.three", []float64{}, 1.3},
		{"two", []float64{1}, 2.0},
		{"zero", []float64{}, 0.0},
		{"zero", []float64{3.0}, 0.0},
		// empty
		{"empty", []float64{}, 0.0},
		{"empty", []float64{5.2}, 5.2},
		// non-existent
		{"five", []float64{}, 0.0},
		{"five", []float64{5.0}, 5.0},
		// invalid
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
		"empty": "",
		"word":  "henry",
	}

	data := []struct {
		key string
		fb  []time.Duration
		out time.Duration
	}{
		// valid
		{"5mins", []time.Duration{}, time.Minute * 5},
		{"1hour", []time.Duration{time.Second * 1}, time.Hour * 1},
		// zero
		{"zero", []time.Duration{}, 0},
		{"zero", []time.Duration{time.Second * 2}, 0},
		// empty
		{"empty", []time.Duration{}, 0},
		{"empty", []time.Duration{time.Second * 2}, time.Second * 2},
		// unset
		{"missing", []time.Duration{}, 0},
		{"missing", []time.Duration{time.Second * 2}, time.Second * 2},
		// invalid
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

// Durations are parsed using time.ParseDuration.
func ExampleGetDuration() {
	// Set some test variables
	os.Setenv("DURATION_NAP", "20m")
	os.Setenv("DURATION_EGG", "5m")
	os.Setenv("DURATION_BIG_EGG", "")
	os.Setenv("DURATION_MATCH", "1.5h")

	// returns time.Duration
	fmt.Println(GetDuration("DURATION_NAP"))
	fmt.Println(GetDuration("DURATION_EGG") * 2)
	// fallback with unset variable
	fmt.Println(GetDuration("DURATION_POWERNAP", time.Minute*45))
	// or an empty one
	fmt.Println(GetDuration("DURATION_BIG_EGG", time.Minute*10))
	fmt.Println(GetDuration("DURATION_MATCH").Minutes())

	// Output:
	// 20m0s
	// 10m0s
	// 45m0s
	// 10m0s
	// 90

	os.Unsetenv("DURATION_NAP")
	os.Unsetenv("DURATION_EGG")
	os.Unsetenv("DURATION_BIG_EGG")
	os.Unsetenv("DURATION_MATCH")
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
		// valid
		{"t", []bool{}, true},
		{"f", []bool{true}, false},
		{"1", []bool{}, true},
		{"0", []bool{true}, false},
		{"true", []bool{}, true},
		{"false", []bool{true}, false},
		// empty
		{"empty", []bool{}, false},
		{"empty", []bool{true}, true},
		// missing
		{"missing", []bool{}, false},
		{"missing", []bool{true}, true},
		// invalid
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
