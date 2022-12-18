// SPDX-FileCopyrightText: 2022 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package casemapper

import (
	"bytes"
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

func TestFrom(t *testing.T) {
	tests := []struct {
		field  string
		fmt    string
		key    string
		adjust map[string]string
		want   bool
	}{
		{field: "FooBar22Tarball", fmt: "twowords", key: "foobar22tarball", want: true},
		{field: "FooBar22Tarball", fmt: "TWOWORDS", key: "FOOBAR22TARBALL", want: true},
		{field: "FooBar22Tarball", fmt: "twoWords", key: "fooBar22Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "TwoWords", key: "FooBar22Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "two_words", key: "foo_bar_22_tarball", want: true},
		{field: "FooBar22Tarball", fmt: "TWO_WORDS", key: "FOO_BAR_22_TARBALL", want: true},
		{field: "FooBar22Tarball", fmt: "two_Words", key: "foo_Bar_22_Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "Two_Words", key: "Foo_Bar_22_Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "two-words", key: "foo-bar-22-tarball", want: true},
		{field: "FooBar22Tarball", fmt: "TWO-WORDS", key: "FOO-BAR-22-TARBALL", want: true},
		{field: "FooBar22Tarball", fmt: "two-Words", key: "foo-Bar-22-Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "Two-Words", key: "Foo-Bar-22-Tarball", want: true},
		{field: "FooBar22Tarball", fmt: "twowords", key: "Foo-Bar-22-Tarball", want: false},
		{
			field: "HTTPHeader",
			fmt:   "Two-Words",
			key:   "HTTP-Header",
			want:  true,
			adjust: map[string]string{
				"HTTP-Header": "HTTPHeader",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.field+"_"+tc.fmt, func(t *testing.T) {
			assert := assert.New(t)

			fn := from(tc.fmt, nil, tc.adjust)
			assert.Equal(tc.want, fn(tc.key, tc.field))
		})
	}
}

func TestIntegration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	type foo struct {
		FooBar     string `test1:"Foo-Bar" test2:"foo_bar"`
		HTTPHeader string `test1:"HTTP-Header" test2:"http_header"`
	}

	gs, err := goschtalt.New(
		goschtalt.AutoCompile(),
		goschtalt.AddValue("record", "test1",
			&foo{
				FooBar:     "test1: foobar",
				HTTPHeader: "test1: http header",
			},
			goschtalt.TagName("test1"),
		),
		goschtalt.AddValue("record", "test2",
			&foo{
				FooBar:     "test2: foobar",
				HTTPHeader: "test2: http header",
			},
			goschtalt.TagName("test2"),
		),
	)
	require.NoError(err)
	require.NotNil(gs)

	foo1, err := goschtalt.Unmarshal[foo](gs, "test1",
		From("Two-Words",
			map[string]string{
				"HTTP-Header": "HTTPHeader",
			},
		))

	assert.NoError(err)
	assert.NotNil(foo1)

	assert.Equal("test1: foobar", foo1.FooBar)
	assert.Equal("test1: http header", foo1.HTTPHeader)

	foo2, err := goschtalt.Unmarshal[foo](gs, "test2", From("two_words"))

	assert.NoError(err)
	assert.NotNil(foo2)

	assert.Equal("test2: foobar", foo2.FooBar)
	assert.Equal("test2: http header", foo2.HTTPHeader)

	foo3, err := goschtalt.Unmarshal[foo](gs, "test2",
		From("two_words",
			map[string]string{
				"http_header": "-",
				"foo_bar":     "-",
			},
		))

	assert.NoError(err)
	assert.NotNil(foo3)

	assert.Equal("", foo3.FooBar)
	assert.Equal("", foo3.HTTPHeader)
}

func TestDebugIntegration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	type foo struct {
		HTTPHeader string `train:"HTTP-Header"`
	}

	gs, err := goschtalt.New(
		goschtalt.AutoCompile(),
		goschtalt.AddValue("record", goschtalt.Root,
			&foo{
				HTTPHeader: "http header",
			},
			goschtalt.TagName("train"),
		),
	)
	require.NoError(err)
	require.NotNil(gs)

	var w bytes.Buffer
	fooVar, err := goschtalt.Unmarshal[foo](gs, goschtalt.Root,
		DebugFrom(&w, "Two-Words",
			map[string]string{
				"HTTP-Header": "HTTPHeader",
			},
		))

	assert.NoError(err)
	assert.NotNil(fooVar)

	assert.Equal("(key) 'HTTP-Header' == 'HTTPHeader' (struct field)\n", w.String())
	assert.Equal("http header", fooVar.HTTPHeader)
	w.Reset()

	type bar struct {
		Name string
	}

	barVar, err := goschtalt.Unmarshal[bar](gs, goschtalt.Root, DebugFrom(&w, "two_words"))

	assert.NoError(err)
	assert.NotNil(barVar)

	assert.Equal("(key) 'HTTP-Header' != 'Name' (struct field)\n", w.String())
	assert.Equal("", barVar.Name)
}
