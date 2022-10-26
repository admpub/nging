package settings

import "github.com/admpub/nging/v5/application/dbschema"

type Config struct {
	Group   string
	Items   map[string]*dbschema.NgingConfig
	Forms   []*SettingForm
	Encoder Encoder
	Decoder Decoder
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Apply() {
	if c.Encoder != nil {
		RegisterEncoder(c.Group, c.Encoder)
	}
	if c.Decoder != nil {
		RegisterDecoder(c.Group, c.Decoder)
	}
	AddDefaultConfig(c.Group, c.Items)
	Register(c.Forms...)
}
