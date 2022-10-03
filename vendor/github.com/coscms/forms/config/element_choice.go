package config

type Choice struct {
	Group   string   `json:"group"`
	Option  []string `json:"option"` //["value","text"]
	Checked bool     `json:"checked"`
}

func (c *Choice) Clone() *Choice {
	r := &Choice{
		Group:   c.Group,
		Option:  make([]string, len(c.Option)),
		Checked: c.Checked,
	}
	copy(r.Option, r.Option)
	return r
}
