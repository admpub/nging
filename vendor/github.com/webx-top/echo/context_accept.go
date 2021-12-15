package echo

import (
	"strconv"
	"strings"
)

const (
	acceptFlagComma      = ','
	acceptFlagSemicolon  = ';'
	acceptFlagSlash      = '/'
	acceptFlagPlus       = '+'
	acceptFlagDot        = '.'
	acceptFlagWhitespace = ' '
)

func NewAccepts(accept string) *Accepts {
	return &Accepts{Raw: accept}
}

type Accepts struct {
	Raw     string
	Accepts []*AcceptQuality
}

func (a *Accepts) Simple(n int) *Accepts {
	accept := strings.SplitN(a.Raw, `;`, 2)[0]
	accepts := strings.SplitN(accept, `,`, n)
	t := &AcceptQuality{Quality: 1, Type: make([]*Accept, len(accepts))}
	a.Accepts = []*AcceptQuality{t}
	for k, r := range accepts {
		c := NewAccept()
		c.Raw = strings.TrimSpace(r)
		c.Mime = c.Raw
		t.Type[k] = c
	}
	return a
}

func (a *Accepts) Advance() *Accepts {
	var (
		r        []rune
		c        = NewAccept()
		types    []*Accept
		isVendor bool
		subtype  []rune
		rootype  []rune
		quality  []rune
		stared   bool
		foundSem bool
	)
	cleanUp := func() {
		if c == nil || len(c.Type) == 0 || len(r) == 0 {
			return
		}
		c.Subtype = append(c.Subtype, string(subtype))
		c.Mime = c.Type + `/` + string(r)
		c.Raw = string(rootype)
		isVendor = false
		stared = false
		types = append(types, c)
		c = NewAccept()
		r = []rune{}
		subtype = []rune{}
		rootype = []rune{}
		if len(quality) > 0 {
			quality = []rune{}
		}
	}
	combine := func() {
		if len(types) == 0 {
			return
		}
		var q float64
		if len(quality) > 0 {
			q, _ = strconv.ParseFloat(strings.TrimPrefix(string(quality), `q=`), 32)
		}
		a.Accepts = append(a.Accepts, &AcceptQuality{
			Quality: float32(q),
			Type:    types[:],
		})
		types = []*Accept{}
	}
	//application/vnd.example.v2+json
	//application/json, text/javascript, */*; q=0.01
	//text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
	for _, v := range a.Raw {
		if v == acceptFlagWhitespace {
			continue
		}
		if foundSem {
			if v != acceptFlagComma {
				quality = append(quality, v)
				continue
			}
			combine()
			quality = []rune{}
			foundSem = false
		}
		if v == acceptFlagSlash { // /
			c.Type = string(r)
			r = []rune{}
			stared = true
			subtype = []rune{}
			rootype = append(rootype, v)
			continue
		}
		if v == acceptFlagDot { // .
			s := string(r)
			if s == `vnd` {
				isVendor = true
			} else {
				if isVendor {
					c.Vendor = append(c.Vendor, s)
				} else {
					subtype = append(subtype, v)
				}
			}
			r = []rune{}
			rootype = append(rootype, v)
			continue
		}
		if v == acceptFlagPlus { // +
			if isVendor {
				c.Vendor = append(c.Vendor, string(r))
			} else {
				c.Subtype = append(c.Subtype, string(subtype))
			}
			isVendor = false
			r = []rune{}
			subtype = []rune{}
			rootype = append(rootype, v)
			continue
		}
		if v == acceptFlagComma { // ,
			cleanUp()
			continue
		}
		if v == acceptFlagSemicolon { // ;
			cleanUp()
			foundSem = true
			continue
		}
		r = append(r, v)
		if stared {
			subtype = append(subtype, v)
		}
		rootype = append(rootype, v)
	}
	cleanUp()
	combine()
	return a
}

type AcceptQuality struct {
	Quality float32 // maximum quality: 1
	Type    []*Accept
}

type Accept struct {
	Raw     string
	Type    string
	Subtype []string
	Mime    string
	Vendor  []string
}

func NewAccept() *Accept {
	return &Accept{}
}
