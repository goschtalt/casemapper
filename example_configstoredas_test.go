// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper_test

import (
	"fmt"

	"github.com/goschtalt/casemapper"
	"github.com/goschtalt/goschtalt"
)

func ExampleConfigStoredAs() {
	gs, err := goschtalt.New(
		goschtalt.AutoCompile(),

		casemapper.ConfigStoredAs("two_words"),

		// Normally you'll be including data from a file or something like that.
		// Here we want to use the built in options to avoid including additional
		// dependencies.
		goschtalt.AddValue("incoming", "",
			&struct {
				Name   string `goschtalt:"FirstName"` // Rename via tags...
				Header string
				Sally  string
			}{
				Name:   "Gopher",
				Header: "Content-Type: text/plain",
				Sally:  "Likes go",
			},

			// or rename via mapping, which is useful if you can't change
			// the target structure, but want the configuration to be a
			// specific key.
			goschtalt.Keymap(
				map[string]string{
					"header": "HTTPHeader", // Convert Http back to HTTP
				},
			),
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

	cfg, err := goschtalt.Unmarshal[Config](gs, "",
		// You can also remap when you unmarshal if that's simpler
		goschtalt.Keymap(
			map[string]string{
				"frog": "sally",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config.HTTPHeader: '%s'\n", cfg.HTTPHeader)
	fmt.Printf("Config.FirstName:  '%s'\n", cfg.FirstName)
	fmt.Printf("Config.Frog:       '%s'\n", cfg.Frog)
	// Output:
	// Config.HTTPHeader: 'Content-Type: text/plain'
	// Config.FirstName:  'Gopher'
	// Config.Frog:       'Likes go'
}
