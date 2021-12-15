package sync

import (
	"github.com/admpub/mysql-schema-sync/internal"
)

func Sync(c *Config, mc *EmailConfig, dbOperators ...internal.DBOperator) (sta *internal.Statics, err error) {
	cfg, err := c.ToConfig(mc)
	if err != nil {
		return nil, err
	}
	return internal.CheckSchemaDiff(cfg, dbOperators...), nil
}
