package confl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	u "github.com/araddon/gou"
)

type encodeError struct{ error }

var (
	KeyEqElement              = ":"
	errArrayMixedElementTypes = errors.New("can't encode array with mixed element types")
	errArrayNilElement        = errors.New("can't encode array with nil element")
	errNonString              = errors.New("can't encode a map with non-string key type")
	errAnonNonStruct          = errors.New("can't encode an anonymous field that is not a struct")
	errArrayNoTable           = errors.New("array element can't contain a table")
	errNoKey                  = errors.New("top-level values must be a Go map or struct")
	errAnything               = errors.New("") // used in testing
	_                         = u.EMPTY
)

var quotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)

// Marshall a go struct into bytes
func Marshal(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encoder controls the encoding of Go values to a document to some
// io.Writer.
//
// The indentation level can be controlled with the Indent field.
type Encoder struct {
	// A single indentation level. By default it is two spaces.
	Indent string

	// hasWritten is whether we have written any output to w yet.
	hasWritten   bool
	w            *bufio.Writer
	KeyEqElement string
}

// NewEncoder returns a encoder that encodes Go values to the io.Writer
// given. By default, a single indentation level is 2 spaces.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:            bufio.NewWriter(w),
		Indent:       "  ",
		KeyEqElement: KeyEqElement,
	}
}

// Encode writes a representation of the Go value to the underlying
// io.Writer. If the value given cannot be encoded to a valid document,
// then an error is returned.
//
// The mapping between Go values and values should be precisely the same
// as for the Decode* functions. Similarly, the TextMarshaler interface is
// supported by encoding the resulting bytes as strings. (If you want to write
// arbitrary binary data then you will need to use something like base64 since
// does not have any binary types.)
//
// When encoding hashes (i.e., Go maps or structs), keys without any
// sub-hashes are encoded first.
//
// If a Go map is encoded, then its keys are sorted alphabetically for
// deterministic output. More control over this behavior may be provided if
// there is demand for it.
//
// Encoding Go values without a corresponding representation---like map
// types with non-string keys---will cause an error to be returned. Similarly
// for mixed arrays/slices, arrays/slices with nil elements, embedded
// non-struct types and nested slices containing maps or structs.
// (e.g., [][]map[string]string is not allowed but []map[string]string is OK
// and so is []map[string][]string.)
func (enc *Encoder) Encode(v interface{}) error {
	rv := eindirect(reflect.ValueOf(v))
	if err := enc.safeEncode(Key([]string{}), rv); err != nil {
		return err
	}
	return enc.w.Flush()
}

func (enc *Encoder) safeEncode(key Key, rv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if terr, ok := r.(encodeError); ok {
				err = terr.error
				return
			}
			panic(r)
		}
	}()
	enc.encode(key, rv)
	return nil
}

func (enc *Encoder) encode(key Key, rv reflect.Value) {
	// Special case. Time needs to be in ISO8601 format.
	// Special case. If we can marshal the type to text, then we used that.
	// Basically, this prevents the encoder for handling these types as
	// generic structs (or whatever the underlying type of a TextMarshaler is).
	switch rv.Interface().(type) {
	case time.Time, TextMarshaler:
		enc.keyEqElement(key, rv)
		return
	}

	k := rv.Kind()
	//u.Debugf("key:%v len:%v  val=%v", key.String(), len(key), k.String())
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		enc.keyEqElement(key, rv)
	case reflect.Array, reflect.Slice:
		if typeEqual(confArrayHash, confTypeOfGo(rv)) {
			enc.eArrayOfTables(key, rv)
		} else {
			enc.keyEqElement(key, rv)
		}
	case reflect.Interface:
		if rv.IsNil() {
			return
		}
		enc.encode(key, rv.Elem())
	case reflect.Map:
		if rv.IsNil() {
			return
		}
		enc.eTable(key, rv)
	case reflect.Ptr:
		if rv.IsNil() {
			return
		}
		enc.encode(key, rv.Elem())
	case reflect.Struct:
		enc.eTable(key, rv)
	default:
		panic(e("Unsupported type for key '%s': %s", key, k))
	}
}

// eElement encodes any value that can be an array element (primitives and
// arrays).
func (enc *Encoder) eElement(rv reflect.Value) {
	switch v := rv.Interface().(type) {
	case time.Time:
		// Special case time.Time as a primitive. Has to come before
		// TextMarshaler below because time.Time implements
		// encoding.TextMarshaler, but we need to always use UTC.
		enc.wf(v.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05Z"))
		return
	case TextMarshaler:
		// Special case. Use text marshaler if it's available for this value.
		if s, err := v.MarshalText(); err != nil {
			encPanic(err)
		} else {
			enc.writeQuoted(string(s))
		}
		return
	}
	switch rv.Kind() {
	case reflect.Bool:
		enc.wf(strconv.FormatBool(rv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.wf(strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		enc.wf(strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32:
		enc.wf(floatAddDecimal(strconv.FormatFloat(rv.Float(), 'f', -1, 32)))
	case reflect.Float64:
		enc.wf(floatAddDecimal(strconv.FormatFloat(rv.Float(), 'f', -1, 64)))
	case reflect.Array, reflect.Slice:
		enc.eArrayOrSliceElement(rv)
	case reflect.Interface:
		enc.eElement(rv.Elem())
	case reflect.String:
		enc.writeQuoted(rv.String())
	default:
		panic(e("Unexpected primitive type: %s", rv.Kind()))
	}
}

// all floats must have a decimal with at least one number on either side.
func floatAddDecimal(fstr string) string {
	if !strings.Contains(fstr, ".") {
		return fstr + ".0"
	}
	return fstr
}

func (enc *Encoder) writeQuoted(s string) {
	enc.wf("\"%s\"", quotedReplacer.Replace(s))
}

func (enc *Encoder) eArrayOrSliceElement(rv reflect.Value) {
	length := rv.Len()
	//u.Infof("arrayorslice?  %v", rv)
	enc.wf("[")
	for i := 0; i < length; i++ {
		elem := rv.Index(i)
		enc.eElement(elem)
		if i != length-1 {
			enc.wf(", ")
		}
	}
	enc.wf("]")
}

func (enc *Encoder) eArrayOfTables(key Key, rv reflect.Value) {
	if len(key) == 0 {
		encPanic(errNoKey)
	}
	//u.Debugf("eArrayOfTables  key=%s key.len=%v rv=%v", key, len(key), rv)
	panicIfInvalidKey(key, true)
	//enc.newline()
	//enc.wf("%s  [\n%s]]", enc.indentStr(key), key.String())
	newKey := key.insert("_")
	keyDelta := 0
	enc.wf("%s%s "+enc.KeyEqElement+" [", enc.indentStrDelta(key, -1), key[len(key)-1])
	for i := 0; i < rv.Len(); i++ {
		trv := rv.Index(i)
		if isNil(trv) {
			continue
		}
		enc.newline()
		enc.wf("%s{", enc.indentStrDelta(key, keyDelta))
		enc.newline()
		//enc.wf("%s{\n%s", enc.indentStr(key), key.String())
		//enc.newline()
		enc.eMapOrStruct(newKey, trv)
		//enc.newline()
		if i == rv.Len()-1 {
			enc.wf("%s}", enc.indentStrDelta(key, keyDelta))
		} else {
			enc.wf("%s},", enc.indentStrDelta(key, keyDelta))
		}
	}
	enc.newline()
	enc.wf("%s]", enc.indentStrDelta(key, -1))
	enc.newline()
}

func (enc *Encoder) eTable(key Key, rv reflect.Value) {
	if len(key) == 1 {
		// Output an extra new line between top-level tables.
		// (The newline isn't written if nothing else has been written though.)
		//enc.newline()
	}
	if len(key) > 0 {
		panicIfInvalidKey(key, true)
		//u.Infof("table?  %v  %v", key, rv)
		enc.wf("%s%s {", enc.indentStrDelta(key, -1), key[len(key)-1])
		enc.newline()
	}
	enc.eMapOrStruct(key, rv)

	if len(key) > 0 {
		enc.wf("%s}", enc.indentStrDelta(key, -1))
		enc.newline()
	}

}

func (enc *Encoder) eMapOrStruct(key Key, rv reflect.Value) {
	switch rv := eindirect(rv); rv.Kind() {
	case reflect.Map:
		enc.eMap(key, rv)
	case reflect.Struct:
		enc.eStruct(key, rv)
	default:
		panic("eTable: unhandled reflect.Value Kind: " + rv.Kind().String())
	}
}

func (enc *Encoder) eMap(key Key, rv reflect.Value) {
	rt := rv.Type()
	var convert func(string) (interface{}, error)
	switch rt.Key().Kind() {
	case reflect.String:
	case reflect.Int:
		convert = func(key string) (interface{}, error) {
			return strconv.Atoi(key)
		}
	case reflect.Int16:
		convert = func(key string) (interface{}, error) {
			r, e := strconv.ParseInt(key, 10, 16)
			return int16(r), e
		}
	case reflect.Int32:
		convert = func(key string) (interface{}, error) {
			r, e := strconv.ParseInt(key, 10, 32)
			return int32(r), e
		}
	case reflect.Int64:
		convert = func(key string) (interface{}, error) {
			return strconv.ParseInt(key, 10, 64)
		}
	case reflect.Uint:
		convert = func(key string) (interface{}, error) {
			r, e := strconv.Atoi(key)
			return uint(r), e
		}
	case reflect.Uint16:
		convert = func(key string) (interface{}, error) {
			r, e := strconv.ParseUint(key, 10, 16)
			return uint16(r), e
		}
	case reflect.Uint32:
		convert = func(key string) (interface{}, error) {
			r, e := strconv.ParseUint(key, 10, 32)
			return uint32(r), e
		}
	case reflect.Uint64:
		convert = func(key string) (interface{}, error) {
			return strconv.ParseUint(key, 10, 64)
		}
	default:
		encPanic(errNonString)
	}
	// Sort keys so that we have deterministic output. And write keys directly
	// underneath this key first, before writing sub-structs or sub-maps.
	var mapKeysDirect, mapKeysSub []string
	for _, mapKey := range rv.MapKeys() {
		//k := mapKey.String()
		k := fmt.Sprint(mapKey.Interface())
		//u.Infof("map key: %v", k)
		if typeIsHash(confTypeOfGo(rv.MapIndex(mapKey))) {
			//u.Debugf("found sub? %s  for %v", k, confTypeOfGo(rv.MapIndex(mapKey)))
			mapKeysSub = append(mapKeysSub, k)
		} else {
			mapKeysDirect = append(mapKeysDirect, k)
		}
	}

	var writeMapKeys = func(mapKeys []string) {
		sort.Strings(mapKeys)
		for _, mapKey := range mapKeys {
			//u.Infof("mapkey: %v", mapKey)
			var v interface{}
			if convert != nil {
				var e error
				v, e = convert(mapKey)
				if e != nil {
					encPanic(e)
				}
			} else {
				v = mapKey
			}
			mrv := rv.MapIndex(reflect.ValueOf(v))
			if isNil(mrv) {
				// Don't write anything for nil fields.
				continue
			}
			enc.encode(key.add(mapKey), mrv)
		}
	}
	writeMapKeys(mapKeysDirect)
	writeMapKeys(mapKeysSub)
}

func (enc *Encoder) eStruct(key Key, rv reflect.Value) {
	// Write keys for fields directly under this key first, because if we write
	// a field that creates a new table, then all keys under it will be in that
	// table (not the one we're writing here).
	rt := rv.Type()
	//var fieldsDirect, fieldsSub [][]int
	var addFields func(rt reflect.Type, rv reflect.Value)
	addFields = func(rt reflect.Type, rv reflect.Value) {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			// skip unexporded fields
			if len(f.PkgPath) > 0 {
				continue
			}
			keyName := f.Tag.Get("confl")
			if keyName == "-" {
				continue
			}
			if keyName == "" {
				keyName = f.Tag.Get("json")
				if keyName == "-" {
					continue
				}
			}
			if len(keyName) > 0 {
				attrs := strings.SplitN(keyName, `,`, 2)
				keyName = attrs[0]
				if keyName == "-" {
					continue
				}
			}
			frv := rv.Field(i)
			if f.Anonymous {
				frv := eindirect(frv)
				t := frv.Type()
				if t.Kind() != reflect.Struct {
					encPanic(errAnonNonStruct)
				}
				addFields(t, frv)
				continue
			}
			if keyName == "" {
				keyName = f.Name
			}
			if keyName != f.Name {
				attrs := strings.Split(keyName, `,`)
				keyName = attrs[0]
				if keyName == "-" {
					continue
				}
				if len(attrs) > 1 {
					var omitEmpty bool
					for _, attr := range attrs[1:] {
						switch attr {
						case `omitempty`:
							omitEmpty = isZero(frv)
						default:
						}
					}
					if omitEmpty {
						continue
					}
				}
			}
			enc.encode(key.add(keyName), frv)
		}
	}
	addFields(rt, rv)
}

// returns the Confl type name of the Go value's type. It is used to
// determine whether the types of array elements are mixed (which is forbidden).
// If the Go value is nil, then it is illegal for it to be an array element, and
// valueIsNil is returned as true.

// Returns the confl type of a Go value. The type may be `nil`, which means
// no concrete confl type could be found.
func confTypeOfGo(rv reflect.Value) confType {
	if isNil(rv) || !rv.IsValid() {
		return nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		return confBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return confInteger
	case reflect.Float32, reflect.Float64:
		return confFloat
	case reflect.Array, reflect.Slice:
		if typeEqual(confHash, confArrayType(rv)) {
			return confArrayHash
		}
		return confArray
	case reflect.Ptr, reflect.Interface:
		return confTypeOfGo(rv.Elem())
	case reflect.String:
		return confString
	case reflect.Map:
		return confHash
	case reflect.Struct:
		switch rv.Interface().(type) {
		case time.Time:
			return confDatetime
		case TextMarshaler:
			return confString
		default:
			return confHash
		}
	default:
		panic("unexpected reflect.Kind: " + rv.Kind().String())
	}
}

// returns the element type of a array. The type returned
// may be nil if it cannot be determined (e.g., a nil slice or a zero length
// slize). This function may also panic if it finds a type that cannot be
// expressed in (such as nil elements, heterogeneous arrays or directly
// nested arrays of tables).
func confArrayType(rv reflect.Value) confType {
	if isNil(rv) || !rv.IsValid() || rv.Len() == 0 {
		return nil
	}
	firstType := confTypeOfGo(rv.Index(0))
	if firstType == nil {
		encPanic(errArrayNilElement)
	}

	rvlen := rv.Len()
	for i := 1; i < rvlen; i++ {
		elem := rv.Index(i)
		switch elemType := confTypeOfGo(elem); {
		case elemType == nil:
			encPanic(errArrayNilElement)
		case !typeEqual(firstType, elemType):
			encPanic(errArrayMixedElementTypes)
		}
	}
	// If we have a nested array, then we must make sure that the nested
	// array contains ONLY primitives.
	// This checks arbitrarily nested arrays.
	if typeEqual(firstType, confArray) || typeEqual(firstType, confArrayHash) {
		nest := confArrayType(eindirect(rv.Index(0)))
		if typeEqual(nest, confHash) || typeEqual(nest, confArrayHash) {
			encPanic(errArrayNoTable)
		}
	}
	return firstType
}

func (enc *Encoder) newline() {
	if enc.hasWritten {
		enc.wf("\n")
	}
}

func (enc *Encoder) keyEqElement(key Key, val reflect.Value) {
	if len(key) == 0 {
		encPanic(errNoKey)
	}
	panicIfInvalidKey(key, false)
	//u.Infof("keyEqElement: %v", key[len(key)-1])
	enc.wf("%s%s "+enc.KeyEqElement+" ", enc.indentStrDelta(key, -1), key[len(key)-1])
	enc.eElement(val)
	enc.newline()
}

func (enc *Encoder) wf(format string, v ...interface{}) {
	if _, err := fmt.Fprintf(enc.w, format, v...); err != nil {
		encPanic(err)
	}
	enc.hasWritten = true
}

func (enc *Encoder) indentStr(key Key) string {
	return strings.Repeat(enc.Indent, len(key))
}
func (enc *Encoder) indentStrDelta(key Key, delta int) string {
	return strings.Repeat(enc.Indent, len(key)+delta)
}

func encPanic(err error) {
	panic(encodeError{err})
}

func eindirect(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return eindirect(v.Elem())
	default:
		return v
	}
}

func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func panicIfInvalidKey(key Key, hash bool) {
	if hash {
		for _, k := range key {
			if !isValidTableName(k) {
				encPanic(e("Key '%s' is not a valid table name. Table names "+
					"cannot contain '[', ']' or '.'.", key.String()))
			}
		}
	} else {
		if !isValidKeyName(key[len(key)-1]) {
			encPanic(e("Key '%s' is not a name. Key names "+
				"cannot contain whitespace.", key.String()))
		}
	}
}

func isValidTableName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r == '[' || r == ']' || r == '.' {
			return false
		}
	}
	return true
}

func isValidKeyName(s string) bool {
	if len(s) == 0 {
		return false
	}
	return true
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		for i := v.NumField() - 1; i >= 0; i-- {
			if len(vt.Field(i).PkgPath) > 0 {
				continue // Private field
			}
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}
