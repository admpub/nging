package forms

import (
	"encoding/xml"

	"github.com/webx-top/com/encoding/json"
)

func NewForms(f *Form) *Forms {
	return &Forms{
		Form: f,
	}
}

type Forms struct {
	*Form
}

// MarshalJSON allows type Pagination to be used with json.Marshal
func (f *Forms) MarshalJSON() ([]byte, error) {
	f.runBefore()
	return json.Marshal(f.Form)
}

// MarshalXML allows type Pagination to be used with xml.Marshal
func (f *Forms) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	f.runBefore()
	return e.EncodeElement(f.Form, start)
}
