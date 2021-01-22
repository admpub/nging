package sync

import (
	"strings"

	"github.com/admpub/mysql-schema-sync/internal"
)

type Config struct {
	Sync            bool
	Drop            bool
	SourceDSN       string
	DestDSN         string
	AlterIgnore     string
	Tables          string
	SkipTables      string
	MailTo          string
	SQLPreprocessor func(string) string `json:"-" xml:"-"`
	Comparer        internal.Comparer   `json:"-" xml:"-"`
}

func (c *Config) ToConfig(mc *EmailConfig) (*internal.Config, error) {
	cfg := new(internal.Config)
	cfg.SourceDSN = c.SourceDSN
	cfg.DestDSN = c.DestDSN
	cfg.Sync = c.Sync
	cfg.Drop = c.Drop
	cfg.SetSQLPreprocessor(c.SQLPreprocessor)
	cfg.SetComparer(c.Comparer)
	c.AlterIgnore = strings.TrimSpace(c.AlterIgnore)
	c.AlterIgnore = strings.TrimLeft(c.AlterIgnore, `{`)
	c.AlterIgnore = strings.TrimRight(c.AlterIgnore, `}`)
	c.AlterIgnore = strings.TrimRight(c.AlterIgnore, `,`)
	if len(c.AlterIgnore) > 0 {
		to := &map[string]*internal.AlterIgnoreTable{}
		err := internal.ParseJSON(`{`+c.AlterIgnore+`}`, &to)
		if err != nil {
			return cfg, err
		}
		cfg.AlterIgnore = *to
	}

	if mc != nil {
		cfg.Email = &internal.EmailStruct{
			SendMailAble: mc.On,
			SMTPHost:     mc.SMTPHost,
			From:         mc.From,
			Password:     mc.Password,
			To:           mc.To,
		}
	}

	if len(c.MailTo) > 0 && cfg.Email != nil {
		cfg.Email.To = c.MailTo
	}

	if cfg.Tables == nil {
		cfg.Tables = []string{}
	}
	if len(c.Tables) > 0 {
		_ts := strings.Split(c.Tables, ",")
		for _, _name := range _ts {
			_name = strings.TrimSpace(_name)
			if len(_name) > 0 {
				cfg.Tables = append(cfg.Tables, _name)
			}
		}
	}
	if cfg.SkipTables == nil {
		cfg.SkipTables = []string{}
	}
	if len(c.SkipTables) > 0 {
		_ts := strings.Split(c.SkipTables, ",")
		for _, _name := range _ts {
			_name = strings.TrimSpace(_name)
			if len(_name) > 0 {
				cfg.SkipTables = append(cfg.SkipTables, _name)
			}
		}
	}
	return cfg, nil
}

type EmailConfig struct {
	On       bool
	SMTPHost string
	From     string
	Password string
	To       string
}
