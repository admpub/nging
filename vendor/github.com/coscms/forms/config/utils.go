package config

func getNames(elements []*Element, languages []*Language) []string {
	var names []string
	for _, elem := range elements {
		if elem.Type == `langset` {
			names = append(names, getNames(elem.Elements, elem.Languages)...)
			continue
		}
		if elem.Type == `fieldset` {
			names = append(names, getNames(elem.Elements, languages)...)
			continue
		}
		if len(elem.Name) > 0 {
			if len(languages) == 0 {
				names = append(names, elem.Name)
			} else {
				for _, lang := range languages {
					names = append(names, lang.Name(elem.Name))
				}
			}
		}
	}
	return names
}
