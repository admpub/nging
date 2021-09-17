package navigate

type NavigateType string

const (
	Left   NavigateType = `left`
	Top    NavigateType = `top`
	Right  NavigateType = `right`
	Bottom NavigateType = `bottom`
)

func NewCollection() *Collection {
	return &Collection{
		Backend:  &Navigates{},
		Frontend: &Navigates{},
	}
}

type Collection struct {
	Backend  *Navigates
	Frontend *Navigates
}

type Navigates map[NavigateType]*List

func (n *Navigates) Add(typ NavigateType, nav *List) {
	(*n)[typ] = nav
}

func (n *Navigates) Get(typ NavigateType) (nav *List) {
	nav = (*n)[typ]
	return
}

func (n *Navigates) Remove(typ NavigateType) bool {
	_, ok := (*n)[typ]
	if ok {
		delete(*n, typ)
	}
	return ok
}
