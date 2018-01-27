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
	"testing"
	"time"
)

type testHost struct {
	ID           string `env:"-"`
	Hostname     string `env:"HOST"`
	Online       bool
	Port         int
	PingInterval time.Duration `env:"PING"`
	PingAverage  float64
}

var (
	testID           = "uid34"
	testHostname     = "test.example.com"
	testOnline       = true
	testPort         = 3000
	testPingInterval = time.Second * 10
	testPingAverage  = 4.5
	// How many visible, non-ignored fields are in testHost
	fieldCount = 5
)

var testEnv = Env{
	"ID":           "not empty",
	"HOST":         testHostname,
	"ONLINE":       fmt.Sprintf("%v", testOnline),
	"PORT":         fmt.Sprintf("%d", testPort),
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

		if th.PingAverage != testPingAverage {
			t.Errorf("Bad PingAverage. Expected=%v, Got=%v", testPingAverage, th.PingAverage)
		}

		if th.PingInterval != testPingInterval {
			t.Errorf("Bad PingInterval. Expected=%v, Got=%v", testPingInterval, th.PingInterval)
		}

	})
}

func TestExtract(t *testing.T) {

	th := &testHost{}
	data := map[string]string{
		"Hostname":     "HOST",
		"Online":       "ONLINE",
		"Port":         "PORT",
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

// TestEnvName tests the envvar name algorithm.
func TestEnvName(t *testing.T) {
	data := []struct {
		in, out string
	}{
		{"URL", "URL"},
		{"Name", "NAME"},
		{"LastName", "LAST_NAME"},
		{"URLEncoding", "URLENCODING"},
		{"LongBeard", "LONG_BEARD"},
		{"HTML", "HTML"},
		{"etc", "ETC"},
	}

	for _, td := range data {
		v := envName(td.in)
		if v != td.out {
			t.Errorf("Bad EnvName (%s).Expected=%v, Got=%v",
				td.in, td.out, v)
		}
	}
}

// TestSplitName tests the name-splitting algorithm.
func TestSplitName(t *testing.T) {
	data := []struct {
		in  string
		out []string
	}{
		{"URL", []string{"url"}},
		{"Name", []string{"name"}},
		{"LastName", []string{"last", "name"}},
		{"URLEncoding", []string{"urlencoding"}},
		{"LongBeard", []string{"long", "beard"}},
		{"HTML", []string{"html"}},
	}

	for _, td := range data {
		v := splitName(td.in)
		if !slicesEqual(td.out, v) {
			t.Errorf("Bad Split (%s). Expected=%v, Got=%v",
				td.in, td.out, v)
		}
	}
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
