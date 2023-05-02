// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper

import (
	"strings"
	"testing"

	"github.com/goschtalt/goschtalt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLowerCamelSnake(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "foo_bar_tarball",
			want: "foo-Bar-Tarball",
		}, {
			in:   "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)

			got := lowerCamelKebab(tc.in)
			assert.Equal(tc.want, got)
		})
	}
}

func TestLowerCamelKebab(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "foo-bar-tarball",
			want: "foo_Bar_Tarball",
		}, {
			in:   "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)

			got := lowerCamelSnake(tc.in)
			assert.Equal(tc.want, got)
		})
	}
}

func TestAllLower(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "foo-bar-tarball",
			want: "foobartarball",
		}, {
			in:   "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)

			got := allLower(tc.in)
			assert.Equal(tc.want, got)
		})
	}
}

func TestAllUpper(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "foo-bar-tarball",
			want: "FOOBARTARBALL",
		}, {
			in:   "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			assert := assert.New(t)

			got := allUpper(tc.in)
			assert.Equal(tc.want, got)
		})
	}
}

func TestMapping(t *testing.T) {
	for key, mapper := range fmtToFunc {
		t.Run(key, func(t *testing.T) {
			assert := assert.New(t)

			got := mapper("two_words")
			assert.Equal(key, got)
		})
	}
}

func TestInvalidMapping(t *testing.T) {
	assert := assert.New(t)

	gs, err := goschtalt.New(goschtalt.AutoCompile(),
		ConfigStoredAs("two_words",
			map[string]string{
				"Foo": "foo",
			},
			map[string]string{
				"Foo": "foo",
			},
		))
	assert.Nil(gs)
	assert.Error(err)
}

func TestMapMerging(t *testing.T) {
	tests := []struct {
		description string
		in          []map[string]string
		sToC        map[string]string
		cToS        map[string]string
		expectErr   error
	}{
		{
			description: "single item",
			in: []map[string]string{
				map[string]string{
					"a": "b",
				},
			},
			sToC: map[string]string{
				"a": "b",
			},
			cToS: map[string]string{
				"b": "a",
			},
		}, {
			description: "multiple in the array",
			in: []map[string]string{
				map[string]string{
					"A": "a",
				},
				map[string]string{
					"B": "b",
				},
			},
			sToC: map[string]string{
				"A": "a",
				"B": "b",
			},
			cToS: map[string]string{
				"a": "A",
				"b": "B",
			},
		}, {
			description: "invalid, duplicated config name",
			in: []map[string]string{
				map[string]string{
					"A": "a",
				},
				map[string]string{
					"B": "b",
				},
				map[string]string{
					"C": "b",
				},
			},
			expectErr: ErrDuplicate,
		}, {
			description: "invalid, duplicated field name",
			in: []map[string]string{
				map[string]string{
					"A": "a",
				},
				map[string]string{
					"B": "b",
				},
				map[string]string{
					"B": "c",
				},
			},
			expectErr: ErrDuplicate,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			a, b, e := merge(tc.in)

			if tc.expectErr == nil {
				assert.Equal(a, tc.sToC)
				assert.Equal(b, tc.cToS)
				assert.NoError(e)
				return
			}

			assert.Nil(a)
			assert.Nil(b)
			assert.ErrorIs(e, tc.expectErr)
		})
	}
}

func TestUnknown(t *testing.T) {
	expected := `Known formats: TWO-WORDS, TWOWORDS, TWO_WORDS, Two-Words, TwoWords, Two_Words, two-Words, two-words, twoWords, two_Words, two_words, twowords`

	assert := assert.New(t)
	require := require.New(t)

	gs, err := goschtalt.New(goschtalt.AutoCompile(), ConfigStoredAs("Invalid"))
	require.Nil(gs)
	require.Error(err)
	assert.True(strings.Contains(err.Error(), expected))
}
