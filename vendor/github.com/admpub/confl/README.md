## Yet another Config Parser for go

This is a config parser most similar to Nginx, supports
Json format, but with line-breaks, comments, etc.   Also, like
Nginx is more lenient.

Uses same syntax as https://github.com/vstakhov/libucl

[![Code Coverage](https://codecov.io/gh/lytics/confl/branch/master/graph/badge.svg)](https://codecov.io/gh/lytics/confl)
[![GoDoc](https://godoc.org/github.com/lytics/confl?status.svg)](http://godoc.org/github.com/lytics/confl)
[![Build Status](https://travis-ci.org/lytics/confl.svg?branch=master)](https://travis-ci.org/lytics/confl)
[![Go ReportCard](https://goreportcard.com/badge/lytics/confl)](https://goreportcard.com/report/lytics/confl)



Use [SublimeText Nginx Plugin](https://github.com/brandonwamboldt/sublime-nginx) for formatting.

Credit to [BurntSushi/Toml](https://github.com/BurntSushi/toml) and [Apcera/Gnatsd](https://github.com/apcera/gnatsd/tree/master/conf) from which 
this was derived.

### Example

```
# nice, a config with comments!

# support the name = value format
title = "conf Example"
# support json semicolon
title2 : "conf example2"
# support omitting = or : because key starts a line
title3 "conf example"
# note, we do not have to have quotes
title4 = Without Quotes

# for Sections we can use brackets
hand {
  name = "Tyrion"
  organization = "Lannisters"
  bio = "Imp"                 // comments on fields
  dob = 1979-05-27T07:32:00Z  # dates, and more comments on fields
}

// Note, double-slash comment
// section name/value that is quoted and json valid, including commas
address : {
  "street"  : "1 Sky Cell",
  "city"    : "Eyre",
  "region"  : "Vale of Arryn",
  "country" : "Westeros"
}

# sections can omit the colon, equal before bracket 
seenwith {
  # nested section
  # can be spaces or tabs for nesting
  jaime : {
    season = season1
    episode = "episode1"
  }

  cersei = {
    season = season1
    episode = "episode1"
  }

}


# Line breaks are OK when inside arrays
seasons = [
  "season1",
  "season2",
  "season3",
  "season4",
  "???"
]


# long strings can use parens to allow multi-line
description (
    we possibly
    can have
    multi line text with a block paren
    block ends with end paren on new line
)



```

And the corresponding Go types are:

```go
type Config struct {
	Title       string
	Hand        HandOfKing
	Location    *Address `confl:"address"`
	Seenwith    map[string]Character
	Seasons     []string
	Description string
}

type HandOfKing struct {
	Name     string
	Org      string `json:"organization"`  // Reads either confl, or json attributes
	Bio      string
	DOB      time.Time
	Deceased bool
}

type Address struct {
	Street  string
	City    string
	Region  string
	ZipCode int
}

type Character struct {
	Episode string
	Season  string
}
```

Note that a case insensitive match will be tried if an exact match can't be
found.

A working example of the above can be found in `_examples/example.{go,conf}`.



### Examples

This package works similarly to how the Go standard library handles `XML`
and `JSON`. Namely, data is loaded into Go values via reflection.

For the simplest example, consider a file as just a list of keys
and values:

```
// Comments in Config
Age = 25
# another comment
Cats = [ "Cauchy", "Plato" ]
# now, using quotes on key
"Pi" = 3.14
Perfection = [ 6, 28, 496, 8128 ]
DOB = 1987-07-05T05:45:00Z
```

Which could be defined in Go as:

```go
type Config struct {
  Age int
  Cats []string
  Pi float64
  Perfection []int
  DOB time.Time 
}
```

And then decoded with:

```go
var conf Config
if err := confl.Unmarshal(byteData, &conf); err != nil {
  // handle error
}
```

You can also use struct tags if your struct field name doesn't map to a confl
key value directly:

```
some_key_NAME = "wat"
```

```go
type Config struct {
  ObscureKey string `confl:"some_key_NAME"`
}
```

### Using the `encoding.TextUnmarshaler` interface

Here's an example that automatically parses duration strings into 
`time.Duration` values:

```
song [
	{
		name = "Thunder Road"
		duration = "4m49s"
	},
	{
		name = "Stairway to Heaven"
		duration = "8m03s"
	}
]
```

Which can be decoded with:

```go
type song struct {
  Name     string
  Duration duration
}
type songs struct {
  Song []song
}
var favorites songs
if err := confl.Unmarshal(blob, &favorites); err != nil {
  log.Fatal(err)
}

for _, s := range favorites.Song {
  fmt.Printf("%s (%s)\n", s.Name, s.Duration)
}
```

And you'll also need a `duration` type that satisfies the 
`encoding.TextUnmarshaler` interface:

```go
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
```

