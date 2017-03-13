package collector

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

func NewGoQuery(url string, reader io.Reader) (*GoQuery, error) {
	q := &GoQuery{
		URL: url,
	}
	var err error
	q.Document, err = goquery.NewDocumentFromReader(reader)
	return q, err
}

type GoQuery struct {
	URL       string
	Document  *goquery.Document
	Selection *goquery.Selection
}

func (g *GoQuery) Parse() error {
	return nil
}
