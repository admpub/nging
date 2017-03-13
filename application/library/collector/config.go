package collector

type PageConfig struct {
	URL         string
	MatchRule   string
	FetchValues map[string]string //{title:'news'}
}
