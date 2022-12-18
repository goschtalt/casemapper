# casemapper

[![Build Status](https://github.com/goschtalt/casemapper/actions/workflows/ci.yml/badge.svg)](https://github.com/goschtalt/casemapper/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/goschtalt/casemapper/coverage.svg?branch=main)](http://codecov.io/github/goschtalt/casemapper?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/goschtalt/casemapper)](https://goreportcard.com/report/github.com/goschtalt/casemapper)
[![GitHub Release](https://img.shields.io/github/release/goschtalt/casemapper.svg)](https://github.com/goschtalt/casemapper/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goschtalt/casemapper)](https://pkg.go.dev/github.com/goschtalt/casemapper)


An optional set of [naming conventions](https://en.wikipedia.org/wiki/Naming_convention_(programming)) mapping options for goschtalt.

Fairly often we want to use structures for configuration that we can't easily alter.
Fortunately [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) provides an
easy way to inject arbitrary remappers/handlers, however that code is fairly boiler
plate and clunky to repeat.  Hence the value of this library.  With the
[janos/casbab](https://github.com/janos/casbab) and a bit of glue code, it's now
easy to automatically convert from config files that may be `snake_case` into go
structures that follow the general conventions.

When unmarshaling simply slide in this option:

```go
	casemapper.From("two_words")
```

With more complex structures there will likely be exceptions needed.  For example,
if the input data has the field `HTTP-Header` but the normal format expects that
value as `Http-Header` no problem; there is an optional `map[string]string` parameter
that maps the **configuration field** to the **structure field**.

```go
	casemapper.From("Two_Words",
		map[string]string{
			"Configuration-Field": "StructureField",
			"HTTP-Header": "HTTPHeader",
		},
    )
```

With the ability to overwrite, you can easily remap whatever you want.  You can
even prevent fields from being populated just as if they included the `mapstructure:"-"`
tag.

```go
	casemapper.From("Two_Words",
		map[string]string{
			"Configuration-Field": "-",
		},
    )
```

References
----------
- https://github.com/mitchellh/mapstructure
- https://github.com/janos/casbab
- https://github.com/goschtalt/goschtalt
