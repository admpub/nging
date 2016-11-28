package echo

import "encoding/xml"

type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: ``,
		Local: `Map`,
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: ``, Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}

// ToData conversion to *Data
func (h H) ToData() *Data {
	var info, zone, data interface{}
	if v, y := h["Data"]; y {
		data = v
	}
	if v, y := h["Zone"]; y {
		zone = v
	}
	if v, y := h["Info"]; y {
		info = v
	}
	var code int
	if v, y := h["Code"]; y {
		if c, y := v.(int); y {
			code = c
		}
	}
	return &Data{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
}

func (h H) DeepMerge(source H) {
	for k, value := range source {
		var (
			destinationValue interface{}
			ok               bool
		)
		if destinationValue, ok = h[k]; !ok {
			h[k] = value
			continue
		}
		sourceM, sourceOk := value.(H)
		destinationM, destinationOk := destinationValue.(H)
		if sourceOk && sourceOk == destinationOk {
			destinationM.DeepMerge(sourceM)
		} else {
			h[k] = value
		}
	}
}
