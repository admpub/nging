package collector

import (
	"github.com/admpub/nging/v5/application/library/cron"
	"github.com/admpub/nging/v5/application/library/module"

	"github.com/nging-plugins/collector/application/handler"
	"github.com/nging-plugins/collector/application/library/setup"
)

const ID = `collector`

var Module = module.Module{
	TemplatePath: map[string]string{
		ID: `collector/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	CronJobs: []*cron.Jobx{
		{
			Name:         `collect_page`,
			RunnerGetter: handler.CollectPageJob,
			Example:      `>collect_page:1`,
			Description:  `网页采集`,
		},
	},
	DBSchemaVer: 0.1000,
}
