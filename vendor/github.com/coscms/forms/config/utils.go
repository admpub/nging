package config

const (
	STATIC   = "static"
	Disabled = "disabled"
	Readonly = "readonly"
)

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
		if len(elem.Name) > 0 && elem.Type != STATIC && !elem.HasAttr(Disabled, Readonly) {
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

func setDefaultValue(elements []*Element, languages []*Language, fieldDefaultValue func(string) string) {
	for _, elem := range elements {
		if elem.Type == `langset` {
			setDefaultValue(elem.Elements, elem.Languages, fieldDefaultValue)
			continue
		}
		if elem.Type == `fieldset` {
			setDefaultValue(elem.Elements, languages, fieldDefaultValue)
			continue
		}
		if len(elem.Value) > 0 {
			continue
		}
		if len(elem.Name) > 0 {
			if len(languages) == 0 {
				elem.Value = fieldDefaultValue(elem.Name)
				continue
			}
			for _, lang := range languages {
				elem.Value = fieldDefaultValue(lang.Name(elem.Name))
			}
		}
	}
}

func setValue(elements []*Element, languages []*Language, fieldValue func(string) string) {
	for _, elem := range elements {
		if elem.Type == `langset` {
			setValue(elem.Elements, elem.Languages, fieldValue)
			continue
		}
		if elem.Type == `fieldset` {
			setValue(elem.Elements, languages, fieldValue)
			continue
		}
		if len(elem.Name) > 0 {
			if len(languages) == 0 {
				elem.Value = fieldValue(elem.Name)
				continue
			}
			for _, lang := range languages {
				elem.Value = fieldValue(lang.Name(elem.Name))
			}
		}
	}
}

func getValue(elements []*Element, languages []*Language, fieldValue func(string, string) error) (err error) {
	for _, elem := range elements {
		if elem.Type == `langset` {
			getValue(elem.Elements, elem.Languages, fieldValue)
			continue
		}
		if elem.Type == `fieldset` {
			getValue(elem.Elements, languages, fieldValue)
			continue
		}
		if len(elem.Name) > 0 {
			if len(languages) == 0 {
				if err = fieldValue(elem.Name, elem.Value); err != nil {
					return
				}
				continue
			}
			for _, lang := range languages {
				if err = fieldValue(lang.Name(elem.Name), elem.Value); err != nil {
					return
				}
			}
		}
	}
	return
}
