// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper_test

import (
	"fmt"

	"github.com/goschtalt/casemapper"
	"github.com/goschtalt/goschtalt"
)

func ExampleFrom() {
	gs, err := goschtalt.New(
		goschtalt.AutoCompile(),

		// Normally you'll be including data from a file or something like that.
		// Here we want to use the built in options to avoid including additional
		// dependencies.
		goschtalt.AddValue("incoming", "",
			&struct {
				Name   string `train:"First-Name"`
				Header string `train:"HTTP-Header"`
			}{
				Name:   "Gopher",
				Header: "Content-Type: text/plain",
			},
			goschtalt.TagName("train"),
		),
	)
	if err != nil {
		panic(err)
	}

	type Config struct {
		HTTPHeader string
		FirstName  string
	}

	cfg, err := goschtalt.Unmarshal[Config](gs, "",
		casemapper.From("Two-Words",
			map[string]string{
				"HTTP-Header": "HTTPHeader",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config.HTTPHeader: '%s'\n", cfg.HTTPHeader)
	fmt.Printf("Config.FirstName:  '%s'\n", cfg.FirstName)
	// Output:
	// Config.HTTPHeader: 'Content-Type: text/plain'
	// Config.FirstName:  'Gopher'
}
