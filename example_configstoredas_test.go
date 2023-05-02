// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper_test

import (
	"fmt"

	"github.com/goschtalt/casemapper"
	"github.com/goschtalt/goschtalt"
	"github.com/goschtalt/goschtalt/pkg/debug"
)

func ExampleConfigIs() {
	c := debug.Collect{}
	gs, err := goschtalt.New(
		goschtalt.AutoCompile(),
		goschtalt.DefaultUnmarshalOptions(goschtalt.KeymapReport(&c)),
		goschtalt.DefaultValueOptions(goschtalt.KeymapReport(&c)),

		casemapper.ConfigStoredAs("two_words",
			// Keys are the struct field names and values are the configuration
			// names.
			map[string]string{
				"Header": "http_header",
				"Sally":  "frog",
			},
		),

		// Normally you'll be including data from a file or something like that.
		// Here we want to use the built in options to avoid including additional
		// dependencies.
		goschtalt.AddValue("incoming", goschtalt.Root,
			&struct {
				Name   string `goschtalt:"FirstName"` // Rename via tags...
				Header string
				Sally  string
			}{
				Name:   "Gopher",
				Header: "Content-Type: text/plain",
				Sally:  "Likes go",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	type Config struct {
		HTTPHeader string
		FirstName  string
		Frog       string
	}

	cfg, err := goschtalt.Unmarshal[Config](gs, goschtalt.Root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Mappings:\n")
	fmt.Print(c.String())

	fmt.Printf("\nConfig:\n")
	fmt.Printf("Config.HTTPHeader: '%s'\n", cfg.HTTPHeader)
	fmt.Printf("Config.FirstName:  '%s'\n", cfg.FirstName)
	fmt.Printf("Config.Frog:       '%s'\n", cfg.Frog)
	// Output:
	// Mappings:
	// 'FirstName'  --> 'first_name'
	// 'Frog'       --> 'frog'
	// 'HTTPHeader' --> 'http_header'
	// 'Header'     --> 'http_header'
	// 'Sally'      --> 'frog'
	//
	// Config:
	// Config.HTTPHeader: 'Content-Type: text/plain'
	// Config.FirstName:  'Gopher'
	// Config.Frog:       'Likes go'
}
