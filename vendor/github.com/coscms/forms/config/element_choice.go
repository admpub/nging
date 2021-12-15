package config

type Choice struct {
	Group   string   `json:"group"`
	Option  []string `json:"option"` //["value","text"]
	Checked bool     `json:"checked"`
}
