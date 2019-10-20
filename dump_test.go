// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package env

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// implements encoding.TextMarshaler, always fails.
type badMarshaller string

func (bm badMarshaller) MarshalText() ([]byte, error) {
	return nil, errors.New("oops")
}

func (bm badMarshaller) UnmarshalText(data []byte) error {
	return errors.New("oops")
}

type BadTarget struct {
	Oops badMarshaller
}

type DumpTarget struct {
	Ignored      string `env:"-"`
	Empty        string
	EmptySlice   []string
	Nil          *string
	Zero         int
	ZeroTime     time.Time
	ZeroDuration time.Duration
	unexported   string

	String  string
	Strings []string

	Bool  bool
	Bools []bool

	Int  int
	Ints []int

	Int8P  *int8   `env:"INT8"`
	Ints8P []*int8 `env:"INTS8"`

	Int16  int16
	Ints16 []int16

	Int32  int32
	Ints32 []int32

	Int64P  *int64   `env:"INT64"`
	Ints64P []*int64 `env:"INTS64"`

	Uint  uint
	Uints []uint

	Uint8P  *uint8   `env:"UINT8"`
	Uints8P []*uint8 `env:"UINTS8"`

	Uint16  uint16
	Uints16 []uint16

	Uint32  uint32
	Uints32 []uint32

	Uint64P  *uint64   `env:"UINT64"`
	Uints64P []*uint64 `env:"UINTS64"`

	Float32  float32
	Floats32 []float32

	Float64  float64
	Floats64 []float64

	Duration  time.Duration
	Durations []time.Duration

	URLP  *url.URL   `env:"URL"`
	URLSP []*url.URL `env:"URLS"`

	// *time.Time implements encoding.TextMarshaler
	TimeP  *time.Time   `env:"TIME"`
	TimesP []*time.Time `env:"TIMES"`

	Nested *Nested

	Unsupported map[string]string
}

func dumpTestValues() (map[string]string, DumpTarget) {
	var (
		i8_1    int8   = 3
		i8_2    int8   = 4
		i64_1   int64  = 9
		i64_2   int64  = 10
		u8_1    uint8  = 3
		u8_2    uint8  = 4
		u64_1   uint64 = 9
		u64_2   uint64 = 10
		url1, _        = url.Parse("http://www.example.com")
		url2, _        = url.Parse("http://www.example.org")
		t1             = time.Now().Truncate(time.Second)
		t2             = time.Now().Add(time.Hour).Truncate(time.Second)
	)

	x := map[string]string{
		"EMPTY":         "",
		"EMPTY_SLICE":   "",
		"NIL":           "",
		"ZERO":          "0",
		"ZERO_TIME":     "0001-01-01T00:00:00Z",
		"ZERO_DURATION": "0s",

		"STRING":  "hello",
		"STRINGS": "hello,dolly",

		"BOOL":  "true",
		"BOOLS": "true,true",

		"INT":    "1",
		"INTS":   "1,2",
		"INT8":   "3",
		"INTS8":  "3,4",
		"INT16":  "5",
		"INTS16": "5,6",
		"INT32":  "7",
		"INTS32": "7,8",
		"INT64":  "9",
		"INTS64": "9,10",

		"UINT":    "1",
		"UINTS":   "1,2",
		"UINT8":   "3",
		"UINTS8":  "3,4",
		"UINT16":  "5",
		"UINTS16": "5,6",
		"UINT32":  "7",
		"UINTS32": "7,8",
		"UINT64":  "9",
		"UINTS64": "9,10",

		"FLOAT32":  "1.1",
		"FLOATS32": "1.1,2.2",

		"FLOAT64":  "3.3",
		"FLOATS64": "3.3,4.4",

		"NESTED_STRING": "nested",
		"NESTED_NUM":    "42",

		"DURATION":  "1m0s",
		"DURATIONS": "1m0s,2m0s",

		"URL":  "http://www.example.com",
		"URLS": "http://www.example.com,http://www.example.org",

		"TIME": t1.Format(time.RFC3339),
		"TIMES": strings.Join([]string{
			t1.Format(time.RFC3339),
			t2.Format(time.RFC3339),
		}, ","),
	}

	v := DumpTarget{
		Ignored:    "ignored",
		Empty:      "",
		unexported: "unexported",

		String:  "hello",
		Strings: []string{"hello", "dolly"},

		Bool:  true,
		Bools: []bool{true, true},

		Int:     1,
		Ints:    []int{1, 2},
		Int8P:   &i8_1,
		Ints8P:  []*int8{&i8_1, &i8_2},
		Int16:   5,
		Ints16:  []int16{5, 6},
		Int32:   7,
		Ints32:  []int32{7, 8},
		Int64P:  &i64_1,
		Ints64P: []*int64{&i64_1, &i64_2},

		Uint:     1,
		Uints:    []uint{1, 2},
		Uint8P:   &u8_1,
		Uints8P:  []*uint8{&u8_1, &u8_2},
		Uint16:   5,
		Uints16:  []uint16{5, 6},
		Uint32:   7,
		Uints32:  []uint32{7, 8},
		Uint64P:  &u64_1,
		Uints64P: []*uint64{&u64_1, &u64_2},

		Float32:  1.1,
		Floats32: []float32{1.1, 2.2},
		Float64:  3.3,
		Floats64: []float64{3.3, 4.4},

		Duration:  time.Minute,
		Durations: []time.Duration{time.Minute, 2 * time.Minute},

		URLP:   url1,
		URLSP:  []*url.URL{url1, url2},
		TimeP:  &t1,
		TimesP: []*time.Time{&t1, &t2},
		Nested: &Nested{
			NestedString: "nested",
			NestedNum:    42,
		},

		Unsupported: map[string]string{
			"foo": "bar",
		},
	}

	return x, v
}

func TestDump(t *testing.T) {
	x, v := dumpTestValues()

	m, err := Dump(v)
	assert.NoError(t, err, "dump failed")
	assert.Equal(t, x, m, "unexpected result")

	// pointer to struct
	m, err = Dump(&v)
	assert.NoError(t, err, "dump failed")
	assert.Equal(t, x, m, "unexpected result")
}

func TestExport(t *testing.T) {
	x, v := dumpTestValues()

	require.NoError(t, Export(v), "export failed")
	for k, v := range x {
		assert.Equalf(t, v, os.Getenv(k), "unexpected %q", k)
	}
	os.Clearenv()
}

func TestIgnoreZeroValues(t *testing.T) {
	x := map[string]string{
		"STRING":     "present",
		"INT":        "1",
		"NESTED_NUM": "2",
	}

	v := DumpTarget{
		String: "present",
		Int:    1,
		Ints:   []int{},
		Nested: &Nested{
			NestedNum: 2,
		},
	}

	m, err := Dump(v, IgnoreZeroValues)
	assert.NoError(t, err, "dump failed")
	assert.Equal(t, x, m, "unexpected result")
}

func TestDump_invalidTarget(t *testing.T) {
	invalid := []interface{}{
		"string",
		[]string{},
		map[string]string{},
		int(10),
	}

	for _, v := range invalid {
		_, err := Dump(v, IgnoreZeroValues)
		assert.EqualError(t, err, "not a struct", "dump accepted invalid target")
	}
}

func TestDump_badFields(t *testing.T) {
	invalid := []interface{}{
		BadTarget{Oops: "oops"},
		struct {
			BadTarget badMarshaller
		}{
			BadTarget: "oops",
		},
		struct {
			Oops []badMarshaller
		}{
			[]badMarshaller{
				"oops",
				"oops",
			},
		},
		struct {
			Oops BadTarget
		}{
			Oops: BadTarget{Oops: "oops"},
		},
	}

	for _, v := range invalid {
		_, err := Dump(v)
		assert.EqualError(t, err, "oops", "BadTarget succeeded")
	}
}

func TestVarNameFunc(t *testing.T) {
	fun := func(name string) string {
		name = VarName(name)
		name = strings.ToLower(name)
		return strings.ReplaceAll(name, "_", "-")
	}

	x := map[string]string{
		"nested-string": "nested",
		"nested-num":    "42",

		"duration":  "1m0s",
		"durations": "1m0s,2m0s",
	}

	v := DumpTarget{
		Duration:  time.Minute,
		Durations: []time.Duration{time.Minute, 2 * time.Minute},
		Nested: &Nested{
			NestedString: "nested",
			NestedNum:    42,
		},
	}

	vars, err := Dump(v, VarNameFunc(fun), IgnoreZeroValues)
	assert.NoError(t, err, "dump failed")
	assert.Equal(t, x, vars, "unexpected vars")
}

func TestExport_invalidTarget(t *testing.T) {
	invalid := []interface{}{
		"string",
		[]string{},
		map[string]string{},
		int(10),
	}

	for _, v := range invalid {
		err := Export(v, IgnoreZeroValues)
		assert.EqualError(t, err, "not a struct", "dump accepted invalid target")
	}
}
