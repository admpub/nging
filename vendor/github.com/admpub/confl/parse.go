// Copyright 2013 Apcera Inc. All rights reserved.

// Conf is a configuration file format used by gnatsd. It is
// a flexible format that combines the best of traditional
// configuration formats and newer styles such as JSON and YAML.
package confl

// The format supported is less restrictive than today's formats.
// Supports mixed Arrays [], nested Maps {}, multiple comment types (# and //)
// Also supports key value assigments using '=' or ':' or whiteSpace()
//   e.g. foo = 2, foo : 2, foo 2
// maps can be assigned with no key separator as well
// semicolons as value terminators in key/value assignments are optional
//
// see parse_test.go for more examples.

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Parser will return a map of keys to interface{}, although concrete types
// underly them. The values supported are string, bool, int64, float64, DateTime.
// Arrays and nested Maps are also supported.
type parser struct {
	mapping map[string]interface{}
	types   map[string]confType
	lx      *lexer

	// A list of keys in the order that they appear in the data.
	ordered []Key

	// the full key for the current hash in scope
	context Key

	// the base key name for everything except hashes
	currentKey string

	// rough approximation of line number
	approxLine int

	// The current scoped context, can be array or map
	ctx interface{}

	// stack of contexts, either map or array/slice stack
	ctxs []interface{}

	// Keys stack
	keys []string

	// A map of 'key.group.names' to whether they were created implicitly.
	implicits map[string]bool
}

type parseError string

func (pe parseError) Error() string {
	return string(pe)
}

func Parse(data string) (map[string]interface{}, error) {
	p, err := parse(data)
	if err != nil {
		return nil, err
	}
	return p.mapping, nil
}

func parse(data string) (p *parser, err error) {

	p = &parser{
		mapping: make(map[string]interface{}),
		lx:      lex(data),
		ctxs:    make([]interface{}, 0, 4),
		keys:    make([]string, 0, 4),
	}
	p.pushContext(p.mapping)

	for {
		it := p.next()
		if it.typ == itemEOF {
			break
		}
		if err := p.processItem(it); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *parser) panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf("Near line %d, key '%s': %s",
		p.approxLine, p.current(), fmt.Sprintf(format, v...))
	panic(parseError(msg))
}

func (p *parser) next() item {
	return p.lx.nextItem()
}

func (p *parser) bug(format string, v ...interface{}) {
	log.Fatalf("BUG: %s\n\n", fmt.Sprintf(format, v...))
}

func (p *parser) expect(typ itemType) item {
	it := p.next()
	p.assertEqual(typ, it.typ)
	return it
}

func (p *parser) assertEqual(expected, got itemType) {
	if expected != got {
		p.bug("Expected '%v' but got '%v'.", expected, got)
	}
}

func (p *parser) pushContext(ctx interface{}) {
	p.ctxs = append(p.ctxs, ctx)
	p.ctx = ctx
}

func (p *parser) popContext() interface{} {
	if len(p.ctxs) == 0 {
		panic("BUG in parser, context stack empty")
	}
	li := len(p.ctxs) - 1
	last := p.ctxs[li]
	p.ctxs = p.ctxs[0:li]
	p.ctx = p.ctxs[len(p.ctxs)-1]
	return last
}

func (p *parser) pushKey(key string) {
	p.keys = append(p.keys, key)
}

func (p *parser) popKey() string {
	if len(p.keys) == 0 {
		panic("BUG in parser, keys stack empty")
	}
	li := len(p.keys) - 1
	last := p.keys[li]
	p.keys = p.keys[0:li]
	return last
}

func (p *parser) processItem(it item) error {
	switch it.typ {
	case itemError:
		//panic("error")
		return fmt.Errorf("Parse error on line %d: '%s'", it.line, it.val)
	case itemKey:
		p.pushKey(it.val)
	case itemMapStart:
		newCtx := make(map[string]interface{})
		p.pushContext(newCtx)
	case itemMapEnd:
		p.setValue(p.popContext())
	case itemString:
		// FIXME(dlc) sanitize string?
		p.setValue(maybeRemoveIndents(it.val))
	case itemInteger:
		num, err := strconv.ParseInt(it.val, 10, 64)
		if err != nil {
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange {
				return fmt.Errorf("Integer '%s' is out of the range.", it.val)
			}
			return fmt.Errorf("Expected integer, but got '%s'.", it.val)
		}
		p.setValue(num)
	case itemFloat:
		num, err := strconv.ParseFloat(it.val, 64)
		if err != nil {
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange {
				return fmt.Errorf("Float '%s' is out of the range.", it.val)
			}
			return fmt.Errorf("Expected float, but got '%s'.", it.val)
		}
		p.setValue(num)
	case itemBool:
		switch it.val {
		case "true":
			p.setValue(true)
		case "false":
			p.setValue(false)
		default:
			return fmt.Errorf("Expected boolean value, but got '%s'.", it.val)
		}
	case itemDatetime:
		dt, err := time.Parse("2006-01-02T15:04:05Z", it.val)
		if err != nil {
			return fmt.Errorf(
				"Expected Zulu formatted DateTime, but got '%s'.", it.val)
		}
		p.setValue(dt)
	case itemArrayStart:
		var array []interface{}
		p.pushContext(array)
	case itemArrayEnd:
		array := p.ctx
		p.popContext()
		p.setValue(array)
	}

	return nil
}

func (p *parser) setValue(val interface{}) {
	// Test to see if we are on an array or a map

	// Array processing
	if ctx, ok := p.ctx.([]interface{}); ok {
		p.ctx = append(ctx, val)
		p.ctxs[len(p.ctxs)-1] = p.ctx
	}

	// Map processing
	if ctx, ok := p.ctx.(map[string]interface{}); ok {
		key := p.popKey()
		// FIXME(dlc), make sure to error if redefining same key?
		ctx[key] = val
	}
}

// setType sets the type of a particular value at a given key.
// It should be called immediately AFTER setValue.
//
// Note that if `key` is empty, then the type given will be applied to the
// current context (which is either a table or an array of tables).
func (p *parser) setType(key string, typ confType) {
	keyContext := make(Key, 0, len(p.context)+1)
	for _, k := range p.context {
		keyContext = append(keyContext, k)
	}
	if len(key) > 0 { // allow type setting for hashes
		keyContext = append(keyContext, key)
	}
	p.types[keyContext.String()] = typ
}

// addImplicit sets the given Key as having been created implicitly.
func (p *parser) addImplicit(key Key) {
	p.implicits[key.String()] = true
}

// removeImplicit stops tagging the given key as having been implicitly created.
func (p *parser) removeImplicit(key Key) {
	p.implicits[key.String()] = false
}

// isImplicit returns true if the key group pointed to by the key was created
// implicitly.
func (p *parser) isImplicit(key Key) bool {
	return p.implicits[key.String()]
}

// current returns the full key name of the current context.
func (p *parser) current() string {
	if len(p.currentKey) == 0 {
		return p.context.String()
	}
	if len(p.context) == 0 {
		return p.currentKey
	}
	return fmt.Sprintf("%s.%s", p.context, p.currentKey)
}

// for multi-line text comments lets remove the Indent
func maybeRemoveIndents(s string) string {
	if !strings.Contains(s, "\n") {
		return s
	}
	lines := strings.Split(s, "\n")
	indent := 0
findIndent:
	for idx, r := range lines[0] {
		switch r {
		case '\t', ' ':
			// keep consuming
		default:
			// first non-whitespace we are going to break
			// and use this as indent size.   This makes a variety of assumptions
			// - subsequent indents use same mixture of spaces/tabs
			indent = idx
			break findIndent
		}
	}

	for i, line := range lines {
		if len(line) >= indent {
			lines[i] = line[indent:]
		}
	}
	return strings.Join(lines, "\n")
}

var escapesReplacer = strings.NewReplacer(
	"\\b", "\u0008",
	"\\t", "\u0009",
	"\\n", "\u000A",
	"\\f", "\u000C",
	"\\r", "\u000D",
	"\\\"", "\u0022",
	"\\/", "\u002F",
	"\\\\", "\u005C",
)

func replaceEscapes(s string) string {
	return escapesReplacer.Replace(s)
}

func (p *parser) replaceUnicode(s string) string {
	indexEsc := func() int {
		return strings.Index(s, "\\u")
	}
	for i := indexEsc(); i != -1; i = indexEsc() {
		asciiBytes := s[i+2 : i+6]
		s = strings.Replace(s, s[i:i+6], p.asciiEscapeToUnicode(asciiBytes), -1)
	}
	return s
}

func (p *parser) asciiEscapeToUnicode(s string) string {
	hex, err := strconv.ParseUint(strings.ToLower(s), 16, 32)
	if err != nil {
		p.bug("Could not parse '%s' as a hexadecimal number, but the "+
			"lexer claims it's OK: %s", s, err)
	}

	// BUG(burntsushi)
	// I honestly don't understand how this works. I can't seem
	// to find a way to make this fail. I figured this would fail on invalid
	// UTF-8 characters like U+DCFF, but it doesn't.
	r := string(rune(hex))
	if !utf8.ValidString(r) {
		p.panicf("Escaped character '\\u%s' is not valid UTF-8.", s)
	}
	return string(r)
}
