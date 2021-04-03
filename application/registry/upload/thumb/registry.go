package thumb

type Registries map[string]Sizes

func (r *Registries) Get(subdir string) Sizes {
	v, y := (*r)[subdir]
	if !y {
		v, _ = (*r)[`*`]
	}
	return v
}

func (r *Registries) Add(subdir string, vs ...Size) *Registries {
	if _, y := (*r)[subdir]; y {
		(*r)[subdir] = append((*r)[subdir], vs...)
	} else {
		(*r)[subdir] = vs
	}
	return r
}

func (r *Registries) Set(subdir string, vs ...Size) *Registries {
	(*r)[subdir] = vs
	return r
}

var Registry = Registries{
	`*`: {DefaultSize},
}
