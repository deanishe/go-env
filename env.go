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
	"strconv"
	"time"
)

// Env is a string: string mapping.
type Env map[string]string

// Get returns the value for envvar "key".
// It accepts one optional "fallback" argument. If no
// envvar is set, returns fallback or an empty string.
func Get(key string, fallback ...string) string {

	var fb string

	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := os.LookupEnv(key)
	if !ok {
		return fb
	}
	return s
}

// GetString is a synonym for Get.
func GetString(key string, fallback ...string) string {
	return Get(key, fallback...)
}

// GetInt returns the value for envvar "key" as an int.
// It accepts one optional "fallback" argument. If no
// envvar is set, returns fallback or 0.
func GetInt(key string, fallback ...int) int {

	var fb int

	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := os.LookupEnv(key)
	if !ok {
		return fb
	}

	// log.Printf("[env] %s=%s", key, s)

	i, err := parseInt(s)
	if err != nil {
		return fb
	}

	return int(i)
}

// GetFloat returns the value for envvar "key" as a float.
// It accepts one optional "fallback" argument. If no
// envvar is set, returns fallback or 0.
func GetFloat(key string, fallback ...float64) float64 {

	var fb float64

	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := os.LookupEnv(key)
	if !ok {
		return fb
	}

	// log.Printf("[env] %s=%s", key, s)

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fb
	}

	return n
}

// GetDuration returns the value for envvar "key" as a time.Duration.
// It accepts one optional "fallback" argument. If no
// envvar is set, returns fallback or 0.
func GetDuration(key string, fallback ...time.Duration) time.Duration {

	var fb time.Duration

	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := os.LookupEnv(key)
	if !ok {
		return fb
	}

	// log.Printf("[env] %s=%s", key, s)

	d, err := time.ParseDuration(s)
	if err != nil {
		return fb
	}

	return d
}

// GetBool returns the value for envvar "key" as a boolean.
// It accepts one optional "fallback" argument. If no
// envvar is set, returns fallback or 0.
func GetBool(key string, fallback ...bool) bool {

	var fb bool

	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := os.LookupEnv(key)
	if !ok {
		return fb
	}

	// log.Printf("[env] %s=%s", key, s)

	b, err := strconv.ParseBool(s)
	if err != nil {
		return fb
	}

	return b
}

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(i), nil
	}

	// Try to parse as float, then convert
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %v", s)
	}
	return int(n), nil
}
