//
// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2018-01-28
//

/*

Command docopt demonstrates binding docopt to a struct via env.

It uses a simple implementation of env.Env to map environment
variable names to docopt options (e.g. LOG_LEVEL -> --log-level),
and to cast the limited types returned by docopt to strings.

	USERNAME=bob go run docopt_example.go --cache-path /var/cache
	// Output:
	// CachePath=/var/cache
	// Username=bob
	// ...
	// ...

*/
package main

import (
	"fmt"
	"log"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"go.deanishe.net/env"
)

var usage = `usage: docopt [options]

Options:
    --cache-path=<dir>  path to cache stuff in
	--username=<name>   your username
	--server=<url>      the server URL
	--debug             show debugging info
`

type options struct {
	CachePath string
	Username  string
	Server    string
	Debug     bool
}

// doptenv adapts docopt's options (map[string]interface{}) to
// env.Bind() by implementing env.Env.
type doptenv struct {
	dopts map[string]interface{}
}

// Lookup implements env.Env, so you can pass it to env.Bind().
//
// env calls this with the name of the environment variables,
// so translate them to options, e.g. CACHE_PATH -> --cache-path,
// before lookup.
//
// docopt options contain 3 types: string, bool and nil.
func (env *doptenv) Lookup(key string) (string, bool) {

	opt := envToOpt(key)

	// Key is not present or wasn't passed
	if v, ok := env.dopts[opt]; !ok || v == nil {
		return "", false
	}

	if b, ok := env.dopts[opt].(bool); ok {
		if b {
			return "true", true
		}
		return "false", true
	}
	if s, ok := env.dopts[opt].(string); ok {
		return s, true
	}

	return "", false
}

// Translate an environment variable to command-line option.
func envToOpt(key string) string {
	s := strings.ToLower(key)
	s = strings.Replace(s, "_", "-", -1)
	return "--" + s
}

func main() {

	// Parse command-line options with docopt
	dopts, err := docopt.Parse(usage, nil, true, "", true)
	if err != nil {
		log.Fatal(err)
	}

	// Default options
	opts := &options{
		CachePath: "/tmp/cache",
	}

	// Update defaults from env
	if err := env.Bind(opts); err != nil {
		log.Fatal(err)
	}

	// Update config from docopt options
	e := &doptenv{dopts}
	if err := env.Bind(opts, e); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CachePath=%s\n", opts.CachePath)
	fmt.Printf("Username=%s\n", opts.Username)
	fmt.Printf("Server=%s\n", opts.Server)
	fmt.Printf("Debug=%v\n", opts.Debug)
}
