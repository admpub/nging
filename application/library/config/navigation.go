package config

var Navigations = []*Navigation{}

type Navigation struct {
	Children []*Navigation
	URL      string
	Name     string
	Label    string
	Pjax     string
	Attrs    map[string]string
}

func NewNavigation() *Navigation {
	return &Navigation{
		Children: []*Navigation{},
		Attrs:    map[string]string{},
	}
}
