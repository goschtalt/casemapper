# casemapper

Functionality was moved into goschtalt v0.21.0 since this is nearly always needed.

[![Build Status](https://github.com/goschtalt/casemapper/actions/workflows/ci.yml/badge.svg)](https://github.com/goschtalt/casemapper/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/goschtalt/casemapper/coverage.svg?branch=main)](http://codecov.io/github/goschtalt/casemapper?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/goschtalt/casemapper)](https://goreportcard.com/report/github.com/goschtalt/casemapper)
[![GitHub Release](https://img.shields.io/github/release/goschtalt/casemapper.svg)](https://github.com/goschtalt/casemapper/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goschtalt/casemapper)](https://pkg.go.dev/github.com/goschtalt/casemapper)


An optional set of [naming conventions](https://en.wikipedia.org/wiki/Naming_convention_(programming)) mapping options for goschtalt.

Fairly often we want to use structures for configuration that can't easily be
altered.  Fortunately goschtalt provides an easy way to inject arbitrary mappers,
however that code is fairly boiler plate and clunky to repeat.  Hence the value of this
library.  With the [janos/casbab](https://github.com/janos/casbab) and a bit of
glue code, it's now easy to automatically convert to/from config files that may
be `snake_case` into go structures that follow the general conventions.

To set configuration from go structs (AddValue()) or read configuration values
from a configuration that is `snake_case` just use the option in place of
`goschtalt.KeymapMapper()`.

```go
	casemapper.ConfigIs("two_words")
```

With more complex structures there will likely be adjustments will be needed.
Adjustments are in the form of a map where keys are the golang struct field name
and values are the configuration key name.

For example:

```go
	casemapper.ConfigStoredAs("two_words",
		map[string]string{
            "CfgField": "configuration_field",
			"CNAMEs":   "cnames",
		},
    )
```

will use the normal snake case mapping except for `CfgField` and `CNAMEs` fields,
which are mapped to `configuration_field` and `cnames` respectively.

References
----------
- https://en.wikipedia.org/wiki/Naming_convention_(programming)
- https://github.com/goschtalt/goschtalt
- https://github.com/janos/casbab
