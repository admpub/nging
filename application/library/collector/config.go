package collector

type PageConfig struct {
	Engine string //goquery,regexp
	Method string //GET,POST
	URL    string
	Rule   string

	/*
		{
			"GET":{
				"name":"value",
				"name":"value",
			},
			"POST":{
				"name":"value",
				"name":"value",
			},
			"Cookie":{
				"name":"value",
				"name":"value",
			}
		}
	*/
	Params map[string]map[string]string
}

func (p *PageConfig) SetParam(args ...string) *PageConfig {
	var method, name, value string
	for i, v := range args {
		switch i {
		case 0:
			method = v
			_, ok := p.Params[method]
			if !ok {
				p.Params[method] = map[string]string{}
			}
		case 1:
			name = v
		case 2:
			value = v
		}
	}
	if len(method) > 0 && len(name) > 0 {
		p.Params[method][name] = value
	}
	return p
}
