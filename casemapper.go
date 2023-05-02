// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

// casemapper bundles up the tedious work of tagging structures with a specific
// naming convention case for configuration.
//
// # Extended Character Sets
//
// This mapping is based on strings and runes so extended character sets are
// supported by default.
package casemapper

import (
	"errors"
	"fmt"
	"sort"
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

var (
	ErrDuplicate = errors.New("duplicate adjustment")
)

var fmtToFunc = map[string]func(string) string{
	"twowords":  allLower,
	"TWOWORDS":  allUpper,
	"twoWords":  casbab.Camel,
	"TwoWords":  casbab.Pascal,
	"two_words": casbab.Snake,
	"TWO_WORDS": casbab.ScreamingSnake,
	"two_Words": lowerCamelSnake,
	"Two_Words": casbab.CamelSnake,
	"two-words": casbab.Kebab,
	"TWO-WORDS": casbab.ScreamingKebab,
	"two-Words": lowerCamelKebab,
	"Two-Words": casbab.CamelKebab,
}

type casemapper struct {
	toCase      func(string) string
	adjustments map[string]string
}

func (c casemapper) Map(in string) string {
	out, found := c.adjustments[in]
	if !found {
		out = c.toCase(in)
	}
	return out
}

// ConfigIs provides a strict field/key mapper that converts the config
// values from the specified nomenclature into the go structure name.
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
// This option provides a goschtalt.KeymapMapper based option that will convert
// every input string, ending the chain 100% of the time.
//
// To make adjustments pass in a map (or many) with keys being the golang
// structure field names and values being the configuration name.
//
// The map keys and values must be unique and the inverse mapping must also be
// unique or an error is returned.
func ConfigIs(format string, structToConfig ...map[string]string) goschtalt.Option {
	sToC, cToS, err := merge(structToConfig)
	if err != nil {
		return goschtalt.WithError(err)
	}

	if toCase, found := fmtToFunc[format]; found {
		return goschtalt.Options(
			goschtalt.DefaultUnmarshalOptions(
				goschtalt.KeymapMapper(&casemapper{
					toCase:      toCase,
					adjustments: cToS,
				}),
			),
			goschtalt.DefaultValueOptions(
				goschtalt.KeymapMapper(&casemapper{
					toCase:      toCase,
					adjustments: sToC,
				}),
			),
		)
	}

	keys := make([]string, 0, len(fmtToFunc))
	for k := range fmtToFunc {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return goschtalt.WithError(
		fmt.Errorf("%w, '%s' unknown format by casemapper.ConfigIs().  Known formats: %s",
			goschtalt.ErrInvalidInput,
			format,
			strings.Join(keys, ", "),
		),
	)
}

func merge(in []map[string]string) (map[string]string, map[string]string, error) {
	sToC := make(map[string]string, len(in))
	cToS := make(map[string]string, len(in))
	for i := range in {
		for k, v := range in[i] {
			_, a := sToC[k]
			_, b := cToS[v]
			if a || b {
				return nil, nil, fmt.Errorf("%w, '%s' is duplicated.", ErrDuplicate, k)
			}
			sToC[k] = v
			cToS[v] = k
		}
	}
	return sToC, cToS, nil
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
