package echo

import (
	"strings"

	"github.com/webx-top/echo/param"
)

var (
	acceptFlagComma      = ','
	acceptFlagSemicolon  = ';'
	acceptFlagSlash      = '/'
	acceptFlagPlus       = '+'
	acceptFlagDot        = '.'
	acceptFlagWhitespace = ' '
	acceptFlagEqual      = '='
)

func NewAccepts(accept string) *Accepts {
	a := &Accepts{Raw: accept}
	return a
}

type Accepts struct {
	Raw    string
	Type   []*Accept
	Params param.StringMap
}

func (a *Accepts) Simple(n int) *Accepts {
	accept := strings.SplitN(a.Raw, `;`, 2)[0]
	accepts := strings.SplitN(accept, `,`, n)
	a.Type = make([]*Accept, len(accepts))
	for k, r := range accepts {
		c := NewAccept()
		c.Raw = strings.TrimSpace(r)
		c.Mime = c.Raw
		a.Type[k] = c
	}
	return a
}

func (a *Accepts) Advance() *Accepts {
	var (
		r        []rune
		c        = NewAccept()
		isVendor bool
		subtype  []rune
		rootype  []rune
		stared   bool
	)
	cleanUp := func() {
		if c == nil {
			return
		}
		c.Subtype = append(c.Subtype, string(subtype))
		c.Mime = c.Type + `/` + string(r)
		c.Raw = string(rootype)
		isVendor = false
		stared = false
		a.Type = append(a.Type, c)
		c = NewAccept()
		r = []rune{}
		subtype = []rune{}
		rootype = []rune{}
	}
	//application/vnd.example.v2+json
	//application/json, text/javascript, */*; q=0.01
	for _, v := range a.Raw {
		if v == acceptFlagWhitespace {
			continue
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
			c = nil
			break
		}
		r = append(r, v)
		if stared {
			subtype = append(subtype, v)
		}
		rootype = append(rootype, v)
	}
	cleanUp()
	return a
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
