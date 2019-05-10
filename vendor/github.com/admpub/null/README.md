## null [![GoDoc](https://godoc.org/github.com/admpub/null?status.svg)](https://godoc.org/github.com/admpub/null) [![Coverage](http://gocover.io/_badge/github.com/admpub/null)](http://gocover.io/github.com/admpub/null)

`null` is a library with reasonable options for dealing with nullable SQL and
JSON values.

Types in `null` will only be considered null on null input, and will JSON
encode to `null`.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this
library in place of `sql.NullXXX`. All types also implement:
`encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`,
`json.Unmarshaler` and `sql.Scanner`.

---

### Installation

Null used to be versioned with gopkg.in, so once you upgrade to v8 and beyond
please stop using gopkg.in and ensure you're using `vgo`, `dep` or vendoring to
version null.

`go get -u "github.com/admpub/null"`

### Usage

`import "github.com/admpub/null"`

| Type | Description | Notes |
|------|-------------|-------|
| `null.JSON` | Nullable `[]byte` | Will marshal to JSON null if Invalid. `[]byte{}` input will not produce an Invalid JSON, but `[]byte(nil)` will. This should be used for storing raw JSON in the database. Also has `null.JSON.Marshal` and `null.JSON.Unmarshal` helpers to marshal and unmarshal foreign objects. |
| `null.Bytes` | Nullable `[]byte` | `[]byte{}` input will not produce an Invalid Bytes, but `[]byte(nil)` will. This should be used for storing binary data (bytes in PSQL for example) in the database. |
| `null.String` | Nullable `string` | |
| `null.Byte` | Nullable `byte` | |
| `null.Bool` | Nullable `bool` | |
| `null.Time` | Nullable `time.Time | Marshals to JSON null if SQL source data is null. Uses `time.Time`'s marshaler. |
| `null.Float32` | Nullable `float32` | |
| `null.Float64` | Nullable `float64` | |
| `null.Int` | Nullable `int` | |
| `null.Int8` | Nullable `int8` | |
| `null.Int16` | Nullable `int16` | |
| `null.Int32` | Nullable `int32` | |
| `null.Int64` | Nullable `int64` | |
| `null.Uint` | Nullable `uint` | |
| `null.Uint8` | Nullable `uint8` | |
| `null.Uint16` | Nullable `uint16` | |
| `null.Uint32` | Nullable `int32` | |
| `null.Int64` | Nullable `uint64` | | |

### Bugs

`json`'s `",omitempty"` struct tag does not work correctly right now. It will
never omit a null or empty String. This might be [fixed
eventually](https://github.com/golang/go/issues/4357).


### License

BSD
