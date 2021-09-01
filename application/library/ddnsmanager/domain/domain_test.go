package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestDomain(t *testing.T) {
	domains := parseDomainArr([]string{
		`a.b.c.test.com.cn`,
		`w.webx.top`,
		`dl.eget.io`,
		`webx.top`,
	})
	com.Dump(domains)
	expected := []*Domain{
		{
			DomainName:   "test.com.cn",
			SubDomain:    "a.b.c",
			UpdateStatus: "",
		},
		{
			DomainName:   "webx.top",
			SubDomain:    "w",
			UpdateStatus: "",
		},
		{
			DomainName:   "eget.io",
			SubDomain:    "dl",
			UpdateStatus: "",
		},
		{
			DomainName:   "webx.top",
			SubDomain:    "",
			UpdateStatus: "",
		},
	}
	assert.Equal(t, expected, domains)
}
