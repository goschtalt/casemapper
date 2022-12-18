// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/goschtalt/goschtalt"
	"resenje.org/casbab"
)

// Popular names of commonly used naming conventions.
const (
	FlatCase           = "twowords"
	UpperCase          = "TWOWORDS"
	CamelCase          = "twoWords"
	PascalCase         = "TwoWords"
	SnakeCase          = "two_words"
	ScreamingSnakeCase = "TWO_WORDS"
	TitleCase          = "Two_Words"
	KebabCase          = "two-words"
	ScreamingKebabCase = "TWO-WORDS"
	TrainCase          = "Two-Words"
)

// From provides a strict mapper that expects configuration values in the format
// specified and maps them to structures.  An optional adjustments map provides
// a simple way to specify any mappings that aren't automatically determined
// correctly.
//
// Since the names of the different formatting styles are not standardized, only
// a few of the common ones have consts defined.  The complete list is below:
//
//   - twowords
//   - TWOWORDS
//   - twoWords
//   - TwoWords
//   - two_words
//   - TWO_WORDS
//   - two_Words
//   - Two_Words
//   - two-words
//   - TWO-WORDS
//   - Two-Words
//   - two-Words
//
// adjustments are in the form of incoming name is the key and the structure
// field name is the value.
func From(format string, adjustments ...map[string]string) goschtalt.DecoderConfigOption {
	adjustments = append(adjustments, map[string]string{})
	return goschtalt.MatchName(from(format, nil, adjustments[0]))
}

// DebugFrom is identical to From() except that it accepts a io.Writer to output
// the comparisons to for debugging purposes.
func DebugFrom(w io.Writer, format string, adjustments ...map[string]string) goschtalt.DecoderConfigOption {
	adjustments = append(adjustments, map[string]string{})
	return goschtalt.MatchName(from(format, w, adjustments[0]))
}

func from(format string, w io.Writer, adjustments map[string]string) func(string, string) bool {
	var toCase func(string) string

	switch format {
	case "twowords":
		toCase = allLower
	case "TWOWORDS":
		toCase = allUpper
	case "twoWords":
		toCase = casbab.Camel
	case "TwoWords":
		toCase = casbab.Pascal
	case "two_words":
		toCase = casbab.Snake
	case "TWO_WORDS":
		toCase = casbab.ScreamingSnake
	case "two_Words":
		toCase = lowerCamelSnake
	case "Two_Words":
		toCase = casbab.CamelSnake
	case "two-words":
		toCase = casbab.Kebab
	case "TWO-WORDS":
		toCase = casbab.ScreamingKebab
	case "Two-Words":
		toCase = casbab.CamelKebab
	case "two-Words":
		toCase = lowerCamelKebab
	}

	cmp := func(key, field string) bool {
		if v, match := adjustments[key]; match {
			if v == field {
				return true
			}
			if v == "-" {
				return false
			}
		}

		// Convert the field to match the desired type since casbab will auto
		// convert from any input into the desired output.
		return key == toCase(field)
	}

	if w == nil {
		return cmp
	}

	return func(key, field string) bool {
		rv := cmp(key, field)
		equality := "!="
		if rv {
			equality = "=="
		}
		fmt.Fprintf(w, "(key) '%s' %s '%s' (struct field)\n", key, equality, field)
		return rv
	}
}

func allLower(s string) string {
	return strings.Join(strings.Split(casbab.Lower(s), " "), "")
}

func allUpper(s string) string {
	return strings.Join(strings.Split(casbab.Screaming(s), " "), "")
}

func lowerCamelSnake(s string) string {
	if len(s) == 0 {
		return ""
	}

	r := []rune(casbab.CamelSnake(s))
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func lowerCamelKebab(s string) string {
	if len(s) == 0 {
		return ""
	}

	r := []rune(casbab.CamelKebab(s))
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
