// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package env

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	env := MapEnv{
		"key":   "value",
		"key2":  "value2",
		"empty": "",
	}

	data := []struct {
		key string
		fb  []string
		out string
	}{
		// valid
		{"key", []string{}, "value"},
		{"key", []string{"value2"}, "value"},
		{"key2", []string{}, "value2"},
		{"key2", []string{"value"}, "value2"},
		// empty
		{"empty", []string{}, ""},
		{"empty", []string{"dave"}, ""},
		// unset
		{"key3", []string{}, ""},
		{"key3", []string{"bob"}, "bob"},
	}

	e := New(env)

	// Verify env is the same
	for k, x := range env {
		v := e.Get(k)
		assert.Equal(t, x, v, "unexpected result")
	}

	// Test Get
	for _, td := range data {
		v := e.Get(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Basic usage of Get. Returns an empty string if variable is unset.
func ExampleGet() {
	// Set some test variables
	_ = os.Setenv("TEST_NAME", "Bob Smith")
	_ = os.Setenv("TEST_ADDRESS", "7, Dreary Lane")

	fmt.Println(Get("TEST_NAME"))
	fmt.Println(Get("TEST_ADDRESS"))
	fmt.Println(Get("TEST_NONEXISTENT")) // unset variable

	// GetString is a synonym
	fmt.Println(GetString("TEST_NAME"))

	// Output:
	// Bob Smith
	// 7, Dreary Lane
	//
	// Bob Smith

	os.Clearenv()
}

// The fallback value is returned if the variable is unset.
func ExampleGet_fallback() {
	// Set some test variables
	_ = os.Setenv("TEST_NAME", "Bob Smith")
	_ = os.Setenv("TEST_ADDRESS", "7, Dreary Lane")
	_ = os.Setenv("TEST_EMAIL", "")

	fmt.Println(Get("TEST_NAME", "default name"))       // fallback ignored
	fmt.Println(Get("TEST_ADDRESS", "default address")) // fallback ignored
	fmt.Println(Get("TEST_EMAIL", "test@example.com"))  // fallback ignored (var is empty, not unset)
	fmt.Println(Get("TEST_NONEXISTENT", "hi there!"))   // unset variable

	// Output:
	// Bob Smith
	// 7, Dreary Lane
	//
	// hi there!

	os.Clearenv()
}

func TestGetInt(t *testing.T) {
	env := MapEnv{
		"one":            "1",
		"two":            "2",
		"zero":           "0",
		"float":          "3.5",
		"word":           "henry",
		"empty":          "",
		"negative":       "-1",
		"negative_float": "-3.5",
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
		{"negative", []int{}, -1},
		{"negative_float", []int{5}, -3},
	}

	e := &Reader{env}
	// Test GetInt
	for _, td := range data {
		v := e.GetInt(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Getting int values with and without fallbacks.
func ExampleGetInt() {
	// Set some test variables
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("PING_INTERVAL", "")

	fmt.Println(GetInt("PORT"))
	fmt.Println(GetInt("PORT", 5000))        // fallback is ignored
	fmt.Println(GetInt("PING_INTERVAL"))     // returns zero value
	fmt.Println(GetInt("PING_INTERVAL", 60)) // returns fallback
	// Output:
	// 3000
	// 3000
	// 0
	// 60

	os.Clearenv()
}

func TestGetUint(t *testing.T) {
	env := MapEnv{
		"one":            "1",
		"two":            "2",
		"zero":           "0",
		"float":          "3.5",
		"word":           "henry",
		"empty":          "",
		"negative":       "-1",
		"negative_float": "-3.5",
	}

	data := []struct {
		key string
		fb  []uint
		out uint
	}{
		// numbers
		{"one", []uint{}, 1},
		{"two", []uint{1}, 2},
		{"zero", []uint{}, 0},
		{"zero", []uint{2}, 0},
		// empty values
		{"empty", []uint{}, 0},
		{"empty", []uint{5}, 5},
		// non-existent values
		{"five", []uint{}, 0},
		{"five", []uint{5}, 5},
		// invalid values
		{"word", []uint{}, 0},
		{"word", []uint{5}, 5},
		{"negative", []uint{}, 0},
		{"negative_float", []uint{5}, 5},
		// floats
		{"float", []uint{}, 3},
		{"float", []uint{5}, 3},
	}

	e := &Reader{env}
	// Test GetInt
	for _, td := range data {
		v := e.GetUint(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Getting int values with and without fallbacks.
func ExampleGetUint() {
	// Set some test variables
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("PING_INTERVAL", "")

	fmt.Println(GetUint("PORT"))
	fmt.Println(GetUint("PORT", 5000))        // fallback is ignored
	fmt.Println(GetUint("PING_INTERVAL"))     // returns zero value
	fmt.Println(GetUint("PING_INTERVAL", 60)) // returns fallback
	// Output:
	// 3000
	// 3000
	// 0
	// 60

	os.Clearenv()
}

func TestGetFloat(t *testing.T) {
	env := MapEnv{
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

	e := &Reader{env}
	// Test GetFloat
	for _, td := range data {
		v := e.GetFloat(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Strings are parsed to floats using strconv.ParseFloat().
func ExampleGetFloat() {
	// Set some test variables
	_ = os.Setenv("TOTAL_SCORE", "172.3")
	_ = os.Setenv("AVERAGE_SCORE", "7.54")

	fmt.Printf("%0.2f\n", GetFloat("TOTAL_SCORE"))
	fmt.Printf("%0.1f\n", GetFloat("AVERAGE_SCORE"))
	fmt.Println(GetFloat("NON_EXISTENT_SCORE", 120.5))
	// Output:
	// 172.30
	// 7.5
	// 120.5

	os.Clearenv()
}

func TestGetDuration(t *testing.T) {
	env := MapEnv{
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

	e := &Reader{env}

	// Test GetDuration
	for _, td := range data {
		v := e.GetDuration(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Durations are parsed using time.ParseDuration.
func ExampleGetDuration() {
	// Set some test variables
	_ = os.Setenv("DURATION_NAP", "20m")
	_ = os.Setenv("DURATION_EGG", "5m")
	_ = os.Setenv("DURATION_BIG_EGG", "")
	_ = os.Setenv("DURATION_MATCH", "1.5h")

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

	os.Clearenv()
}

func TestGetBool(t *testing.T) {
	env := MapEnv{
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

	e := &Reader{env}

	// Test GetBool
	for _, td := range data {
		v := e.GetBool(td.key, td.fb...)
		assert.Equal(t, td.out, v, "unexpected result")
	}
}

// Strings are parsed using strconv.ParseBool().
func ExampleGetBool() {
	// Set some test variables
	_ = os.Setenv("LIKE_PEAS", "t")
	_ = os.Setenv("LIKE_CARROTS", "true")
	_ = os.Setenv("LIKE_BEANS", "1")
	_ = os.Setenv("LIKE_LIVER", "f")
	_ = os.Setenv("LIKE_TOMATOES", "0")
	_ = os.Setenv("LIKE_BVB", "false")
	_ = os.Setenv("LIKE_BAYERN", "FALSE")

	// strconv.ParseBool() supports many formats
	fmt.Println(GetBool("LIKE_PEAS"))
	fmt.Println(GetBool("LIKE_CARROTS"))
	fmt.Println(GetBool("LIKE_BEANS"))
	fmt.Println(GetBool("LIKE_LIVER"))
	fmt.Println(GetBool("LIKE_TOMATOES"))
	fmt.Println(GetBool("LIKE_BVB"))
	fmt.Println(GetBool("LIKE_BAYERN"))

	// Fallback
	fmt.Println(GetBool("LIKE_BEER", true))

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
	// false
	// true

	os.Clearenv()
}
