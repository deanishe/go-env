
go-env
======

Easily access environment variables from Go, and bind them to structs.


<!-- vim-markdown-toc GFM -->

* [Usage](#usage)
    * [Access environment variables](#access-environment-variables)
    * [Binding](#binding)
    * [Customisation](#customisation)
* [Installation](#installation)
* [Documentation](#documentation)
* [Licence](#licence)

<!-- vim-markdown-toc -->

Usage
-----

Import path is `github.com/deanishe/go-env`, import name is `env`.

You can directly access environment variables, or populate your structs from them using struct tags and `env.Bind()`.


### Access environment variables

Read `int`, `float64`, `duration` and `string` values from environment variables, with optional fallback values for unset variables.

```go
import "github.com/deanishe/go-env"

// Get value for key or return empty string
s := env.Get("SHELL")

// Get value for key or return a specified default
s := env.Get("DOES_NOT_EXIST", "fallback value")

// Get an int (or 0 if SOME_NUMBER is unset or empty)
i := env.GetInt("SOME_NUMBER")

// Int with a fallback
i := env.GetInt("SOME_UNSET_NUMBER", 10)
```


### Binding

You can also populate a struct directly from the environment by appropriately tagging it and calling `env.Bind()`:

```go
// Simple configuration struct
type config struct {
    HostName string `env:"HOSTNAME"` // default would be HOST_NAME
    Port     int    // leave as default (PORT)
    SSL      bool   `env:"USE_SSL"` // default would be SSL
    Online   bool   `env:"-"`       // ignore this field
}

// Set some values in the environment for test purposes
os.Setenv("HOSTNAME", "api.example.com")
os.Setenv("PORT", "443")
os.Setenv("USE_SSL", "1")
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
// Output:
// api.example.com
// 443
// true
// false
```


### Customisation ###

Variables are retrieved via implementors of the `env.Env` interface (which `env.Bind()` accepts as a second, optional parameter):

```go
type Env interface {
	// Lookup retrieves the value of the variable named by key.
	//
	// It follows the same semantics as os.LookupEnv(). If a variable
	// is unset, the boolean will be false. If a variable is set, the
	// boolean will be true, but the variable may still be an empty
	// string.
	Lookup(key string) (string, bool)
}
```

So you can pass a custom `Env` implementation to `Bind()` in order to populate structs from a source other than environment variables.

See [examples/docopt][docopt] to see how to implement a custom `Env` that populates a struct from `docopt` command-line options.


Installation
------------

```bash
go get github.com/deanishe/go-env
```


Documentation
-------------

Read the documentation on [GoDoc][godoc].


Licence
-------

This library is released under the [MIT Licence][mit].

[mit]: ./LICENCE.txt
[docopt]: ./examples/docopt/docopt_example.go
[godoc]: https://godoc.org/github.com/deanishe/go-env

