//
// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2018-01-28
//

/*

Package env reads environment variables and populates structs from them.

Supports boolean, string, int, float64 and time.Duration.


Reading variables

Read environment variables with the Get* functions:

	// String
	s := env.Get("HOME")
	// String with fallback
	s = env.Get("NON_EXISTENT_VAR", "default") // -> default

	// Int
	i := env.GetInt("SOME_COUNT")
	// Int with fallback
	i = env.GetInt("NON_EXISTENT_VAR", 10) // -> 10

	// Duration (e.g. "1m", "2h30m", "5.3s")
	d := env.GetDuration("SOME_TIME")
	// Duration with fallback
	d = env.GetDuration("NON_EXISTENT_VAR", time.Minute * 120) // -> 2h0m


Populating structs

Populate a struct from the environment by passing it to Bind():

	type options struct {
		Hostname string
		Port     int
	}

	o := &options{}
	if err := env.Bind(o); err != nil {
		// handle error...
	}

	fmt.Println(o.Hostname) // -> value of HOSTNAME environment variable
	fmt.Println(o.Port)     // -> value of PORT environment variable

Use tags to specify a variable name or ignore a field:

	type options {
		HostName string `env:"HOSTNAME"` // default would be HOST_NAME
		Port     int
		Online   bool `env:"-"` // ignored
	}


Tags

Add `env:"..."` tags to your struct fields to bind them to specific
environment variables or ignore them. `env:"-"` tells Bind() to
ignore the field:

	type options {
		UserName      string
		LastUpdatedAt string `env:"-"` // Not loaded from environment
	}

Add `env:"VARNAME"` to bind a field to the variable VARNAME:

	type options {
		UserName string `env:"USERNAME"`   // default = USER_NAME
		APIKey   string `env:"APP_SECRET"` // default = API_KEY
	}


Customisation

Variables are retrieved via implementors of the Env interface, which
Bind() accepts as a second, optional parameter.

So you can pass a custom Env implementation to Bind() to populate
structs from a source other than environment variables.

See examples/docopt to see a custom Env implementation used to
populate a struct from docopt command-line options.


Licence

This library is released under the MIT Licence.

*/
package env
