package scron

type Cron struct {
	PoolSize int    `json:"poolSize"`
	Template string `json:"template"` //发信模板
}
