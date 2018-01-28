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

type testHost struct {
	ID           string `env:"-"`
	Hostname     string `env:"HOST"`
	Online       bool
	Port         uint
	Score        int
	FreeSpace    int64         `env:"SPACE"`
	PingInterval time.Duration `env:"PING"`
	PingAverage  float64
}

var (
	testID                 = "uid34"
	testHostname           = "test.example.com"
	testOnline             = true
	testPort         uint  = 3000
	testScore              = 10000
	testFreeSpace    int64 = 9876543210
	testPingInterval       = time.Second * 10
	testPingAverage        = 4.5
	// How many visible, non-ignored fields are in testHost
	fieldCount = 7
)

var testEnv = Env{
	"ID":           "not empty",
	"HOST":         testHostname,
	"ONLINE":       fmt.Sprintf("%v", testOnline),
	"PORT":         fmt.Sprintf("%d", testPort),
	"SCORE":        fmt.Sprintf("%d", testScore),
	"SPACE":        fmt.Sprintf("%d", testFreeSpace),
	"PING":         fmt.Sprintf("%s", testPingInterval),
	"PING_AVERAGE": fmt.Sprintf("%0.1f", testPingAverage),
}

func TestBind(t *testing.T) {

	withEnv(testEnv, func(e Env) {

		th := &testHost{}

		// Bind to the environment
		if err := Bind(th); err != nil {
			t.Fatalf("couldn't bind testHost: %v", err)
		}

		if th.ID != "" {
			t.Errorf("Non-empty ID. Got=%v", th.ID)
		}

		if th.Hostname != testHostname {
			t.Errorf("Bad Hostname. Expected=%v, Got=%v", testHostname, th.Hostname)
		}

		if th.Online != testOnline {
			t.Errorf("Bad Online. Expected=%v, Got=%v", testOnline, th.Online)
		}

		if th.Port != testPort {
			t.Errorf("Bad Port. Expected=%v, Got=%v", testPort, th.Port)
		}

		if th.Score != testScore {
			t.Errorf("Bad Score. Expected=%v, Got=%v", testScore, th.Score)
		}

		if th.FreeSpace != testFreeSpace {
			t.Errorf("Bad FreeSpace. Expected=%v, Got=%v", testFreeSpace, th.FreeSpace)
		}

		if th.PingAverage != testPingAverage {
			t.Errorf("Bad PingAverage. Expected=%v, Got=%v", testPingAverage, th.PingAverage)
		}

		if th.PingInterval != testPingInterval {
			t.Errorf("Bad PingInterval. Expected=%v, Got=%v", testPingInterval, th.PingInterval)
		}

		// Bind to invalid type
		sl := []string{}
		if err := Bind(&sl); err == nil {
			t.Error("Invalid Type accepted.")
		}

		// Unsupported field type
		st := struct {
			Host   string
			Online []string // slices aren't supported
		}{}

		if err := Bind(&st); err == nil {
			t.Error("Invalid field accepted.")
		}

		// Field not found error
		st2 := struct {
			Host string
			Port uint
		}{}

		binds, err := extract(&st2)
		if err != nil {
			t.Fatalf("couldn't bind to struct: %v", err)
		}
		// Change field names
		for _, bind := range binds {
			bind.Name = bind.Name + "2"
		}
		// Fail to load fields
		for _, bind := range binds {
			if err := bind.Load(); err == nil {
				t.Errorf("Accepted bad binding (%s)", bind.Name)
			}
		}
	})
}

// Populate a struct from environment variables.
func ExampleBind() {

	// Simple configuration struct
	type config struct {
		HostName     string        `env:"HOSTNAME"` // default = HOST_NAME
		SSL          bool          `env:"USE_SSL"`  // default = SSL
		Port         int           // leave as default (PORT)
		PingInterval time.Duration `env:"PING"` // default = PING_INTERVAL
		Online       bool          `env:"-"`    // ignore this field
	}

	// Set some values in the environment for test purposes
	os.Setenv("HOSTNAME", "api.example.com")
	os.Setenv("PORT", "443")
	os.Setenv("USE_SSL", "1")
	os.Setenv("PING", "5m")
	os.Setenv("ONLINE", "1") // will be ignored

	// Create a config and bind it to the environment
	c := &config{}
	if err := Bind(c); err != nil {
		// handle error...
	}

	// config struct now populated from the environment
	fmt.Println(c.HostName)
	fmt.Printf("%d\n", c.Port)
	fmt.Printf("%v\n", c.SSL)
	fmt.Printf("%v\n", c.Online)
	fmt.Printf("%v\n", c.PingInterval*4) // it's not a string!
	// Output:
	// api.example.com
	// 443
	// true
	// false
	// 20m0s

	for _, k := range []string{
		"HOSTNAME",
		"PORT",
		"USE_SSL",
		"ONLINE",
		"PING"} {

		os.Unsetenv(k)
	}
}

func TestExtract(t *testing.T) {

	th := &testHost{}
	data := map[string]string{
		"Hostname":     "HOST",
		"Online":       "ONLINE",
		"Port":         "PORT",
		"Score":        "SCORE",
		"FreeSpace":    "SPACE",
		"PingInterval": "PING",
		"PingAverage":  "PING_AVERAGE",
	}

	binds, err := extract(th)
	if err != nil {
		t.Fatalf("couldn't extract testHost: %v", err)
	}

	if len(binds) != fieldCount {
		t.Errorf("Bad Bindings count. Expected=%d, Got=%d",
			fieldCount, len(binds))
	}

	x := map[string]string{}
	for _, bind := range binds {
		x[bind.Name] = bind.EnvVar
	}

	if err := testMapsEqual(x, data); err != nil {
		t.Fatalf("extract failed: %v", err)
	}
}

// TestVarName tests the envvar name algorithm.
func TestVarName(t *testing.T) {
	data := []struct {
		in, out string
	}{
		{"URL", "URL"},
		{"Name", "NAME"},
		{"LastName", "LAST_NAME"},
		{"URLEncoding", "URL_ENCODING"},
		{"LongBeard", "LONG_BEARD"},
		{"HTML", "HTML"},
		{"etc", "ETC"},
	}

	for _, td := range data {
		v := VarName(td.in)
		if v != td.out {
			t.Errorf("Bad VarName (%s). Expected=%v, Got=%v",
				td.in, td.out, v)
		}
	}
}

// Example output of VarName.
func ExampleVarName() {
	// single-case words are upper-cased
	fmt.Println(VarName("URL"))
	fmt.Println(VarName("name"))
	// words that start with fewer than 3 uppercase chars are
	// upper-cased
	fmt.Println(VarName("Folder"))
	fmt.Println(VarName("MTime"))
	// but with 3+ uppercase chars, the last is treated as the first
	// char of the next word
	fmt.Println(VarName("VIPath"))
	fmt.Println(VarName("URLEncoding"))
	fmt.Println(VarName("SSLPort"))
	// camel-case words are split on the case changes
	fmt.Println(VarName("LastName"))
	fmt.Println(VarName("LongHorse"))
	fmt.Println(VarName("loginURL"))
	fmt.Println(VarName("newHomeAddress"))
	fmt.Println(VarName("PointA"))
	// digits are considered as the end of a word, not the start
	fmt.Println(VarName("b2B"))
	// Output:
	// URL
	// NAME
	// FOLDER
	// MTIME
	// VI_PATH
	// URL_ENCODING
	// SSL_PORT
	// LAST_NAME
	// LONG_HORSE
	// LOGIN_URL
	// NEW_HOME_ADDRESS
	// POINT_A
	// B2_B
}

func slicesEqual(one, two []string) bool {
	if len(one) != len(two) {
		return false
	}
	for i, v1 := range one {
		if v1 != two[i] {
			return false
		}
	}
	return true
}

func testMapsEqual(a, b map[string]string) error {
	if len(a) != len(b) {
		return fmt.Errorf("different lengths (%d != %d)", len(a), len(b))
	}

	for k, v := range a {

		v2, ok := b[k]
		if !ok {
			return fmt.Errorf("key %v missing in b", k)

		} else if v2 != v {
			return fmt.Errorf("%s is different: %#v != %#v", k, v, v2)
		}
	}

	return nil
}

func TestIsCamelCase(t *testing.T) {
	data := []struct {
		s string
		v bool
	}{
		{"", false},
		{"URL", false},
		{"url", false},
		{"Url", false},
		{"HomeAddress", true},
		{"myHomeAddress", true},
		{"PlaceA", true},
		{"myPlaceB", true},
		{"myB", true},
		{"my2B", true},
		{"B2B", false},
		{"SSLPort", true},
	}

	for _, td := range data {
		b := isCamelCase(td.s)
		if b != td.v {
			t.Errorf("Bad CamelCase (%s). Expected=%v, Got=%v", td.s, td.v, b)
		}

	}
}

func TestSplitCamelCase(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"HomeAddress", "HOME_ADDRESS"},
		{"homeAddress", "HOME_ADDRESS"},
		{"loginURL", "LOGIN_URL"},
		{"SSLPort", "SSL_PORT"},
		{"HomeAddress", "HOME_ADDRESS"},
		{"myHomeAddress", "MY_HOME_ADDRESS"},
		{"PlaceA", "PLACE_A"},
		{"myPlaceB", "MY_PLACE_B"},
		{"myB", "MY_B"},
		{"my2B", "MY2_B"},
		{"URLEncoding", "URL_ENCODING"},
	}

	for _, td := range data {
		s := splitCamelCase(td.in)
		if s != td.out {
			t.Errorf("Bad SplitCamel (%s). Expected=%v, Got=%v", td.in, td.out, s)
		}

	}
}
