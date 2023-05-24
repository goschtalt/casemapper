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

		casemapper.ConfigIs("two_words",
			// Keys are the struct field names and values are the configuration
			// names.
			map[string]string{
				"CNAMEs": "cnames",
			},
		),

		// Normally you'll be including data from a file or something like that.
		// Here we want to use the built in options to avoid including additional
		// dependencies.
		goschtalt.AddValue("incoming", goschtalt.Root,
			&struct {
				CNAMEs string
				Header string
				Sally  string
			}{
				CNAMEs: "C Names",
				Header: "Content-Type: text/plain",
				Sally:  "Likes go",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	type Config struct {
		CNAMEs string
		Header string
		Sally  string
	}

	cfg, err := goschtalt.Unmarshal[Config](gs, goschtalt.Root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Mappings:\n")
	fmt.Print(c.String())

	fmt.Printf("\nConfig:\n")
	fmt.Printf("Config.CNAMEs: '%s'\n", cfg.CNAMEs)
	fmt.Printf("Config.Header: '%s'\n", cfg.Header)
	fmt.Printf("Config.Sally:  '%s'\n", cfg.Sally)
	// Output:
	// Mappings:
	// 'CNAMEs' --> 'cnames'
	// 'Header' --> 'header'
	// 'Sally'  --> 'sally'
	//
	// Config:
	// Config.CNAMEs: 'C Names'
	// Config.Header: 'Content-Type: text/plain'
	// Config.Sally:  'Likes go'
}
