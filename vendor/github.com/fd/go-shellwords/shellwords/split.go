package shellwords

import "fmt"
import "unicode"

// Split line of Shell words.
func Split(line string) ([]string, error) {
	p := parser{
		runes: []rune(line),
		word:  make([]rune, 0, 4*1024),
		words: make([]string, 0, 1024),
		i:     -1,
	}

	err := p.parse()
	if err != nil {
		return nil, err
	}

	return p.words, nil
}

type parser struct {
	words []string
	word  []rune
	runes []rune
	rune  rune
	i     int
}

func (p *parser) next() {
	if p.i >= (len(p.runes) - 1) {
		p.rune = -1
	} else {
		p.i++
		p.rune = p.runes[p.i]
	}
}

func (p *parser) parse() error {
	p.next() // init
	p.consume_whitespace()
	consumed := false

	for p.rune != -1 {
		err := p.capture_word()
		if err != nil {
			return err
		}
		consumed = true

		if p.consume_whitespace() {
			p.words = append(p.words, string(p.word))
			p.word = p.word[:0]
			consumed = false
		}
	}

	if consumed {
		p.words = append(p.words, string(p.word))
		p.word = p.word[:0]
	}

	return nil
}

func (p *parser) consume_whitespace() bool {
	found := false
	for unicode.IsSpace(p.rune) {
		found = true
		p.next()
	}
	return found
}

func (p *parser) capture_word() error {
	switch p.rune {
	case '\'': // single quote
		return p.capture_sq_word()
	case '"': // double quote
		return p.capture_dq_word()
	case '\\': // escape
		return p.capture_escape()
	default:
		return p.capture_simple_word()
	}
	panic("")
}

func (p *parser) capture_sq_word() error {
	p.next() // skip '

	for p.rune != '\'' && p.rune >= 0 {
		p.word = append(p.word, p.rune)
		p.next()
	}

	if p.rune != '\'' {
		return fmt.Errorf("Expected a `'` (single quote)")
	}

	p.next() // skip '
	return nil
}

func (p *parser) capture_dq_word() error {
	p.next() // skip "

	for p.rune != '"' && p.rune >= 0 {
		if p.rune == '\\' {
			p.next()
			if p.rune < 0 {
				return fmt.Errorf("Unexpected end of shell word")
			}
			p.word = append(p.word, p.rune)
			p.next()
		} else {
			p.word = append(p.word, p.rune)
			p.next()
		}
	}

	if p.rune != '"' {
		return fmt.Errorf("Expected a `\"` (double quote)")
	}

	p.next() // skip "
	return nil
}

func (p *parser) capture_escape() error {
	p.next() // skip \
	p.word = append(p.word, p.rune)
	p.next()
	return nil
}

func (p *parser) capture_simple_word() error {
	for !unicode.IsSpace(p.rune) && p.rune != '\\' && p.rune != '"' && p.rune != '\'' && p.rune >= 0 {
		p.word = append(p.word, p.rune)
		p.next()
	}

	return nil
}
