package cron

var (
	// systemJobs 系统Job
	systemJobs = map[string]Jobx{}
)

func Register(name string, fn RunnerGetter, example string, description string) {
	AddSystemJob(name, fn, example, description)
}

// AddSystemJob 添加系统Job
func AddSystemJob(name string, fn RunnerGetter, example string, description string) {
	systemJobs[name] = Jobx{
		Example:      example,
		Description:  description,
		RunnerGetter: fn,
	}
}

type Jobx struct {
	Example      string //">funcName:param"
	Description  string
	RunnerGetter RunnerGetter
}
