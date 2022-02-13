package navigate

type NavigateType string

const (
	Left   NavigateType = `left`
	Top    NavigateType = `top`
	Right  NavigateType = `right`
	Bottom NavigateType = `bottom`
)

func NewCollection(baseProject string) *Collection {
	return &Collection{
		Backend: &ProjectNavigates{
			Navigates:   &Navigates{},
			baseProject: baseProject,
			projects:    map[string]*Navigates{},
		},
		Frontend: &Navigates{},
	}
}

type ProjectNavigates struct {
	*Navigates
	baseProject string
	projects    map[string]*Navigates
}

func (p *ProjectNavigates) Project(project string) *Navigates {
	if p.baseProject == project {
		return p.Navigates
	}
	nav, ok := p.projects[project]
	if !ok {
		nav = &Navigates{}
		p.projects[project] = nav
	}
	return nav
}

type Collection struct {
	Backend  *ProjectNavigates
	Frontend *Navigates
}

type Navigates map[NavigateType]*List

func (n *Navigates) Add(typ NavigateType, nav *List) {
	(*n)[typ] = nav
}

func (n *Navigates) AddItems(typ NavigateType, index int, items ...*Item) {
	nav := n.Get(typ)
	if nav == nil {
		nav = &List{}
		(*n)[typ] = nav
	}
	nav.Add(index, items...)
}

func (n *Navigates) AddTopItems(index int, items ...*Item) {
	n.AddItems(Top, index, items...)
}

func (n *Navigates) AddLeftItems(index int, items ...*Item) {
	n.AddItems(Left, index, items...)
}

func (n *Navigates) AddRightItems(index int, items ...*Item) {
	n.AddItems(Right, index, items...)
}

func (n *Navigates) AddBottomItems(index int, items ...*Item) {
	n.AddItems(Bottom, index, items...)
}

func (n *Navigates) Get(typ NavigateType) (nav *List) {
	nav = (*n)[typ]
	return
}

func (n *Navigates) GetTop() *List {
	return n.Get(Top)
}

func (n *Navigates) GetLeft() *List {
	return n.Get(Left)
}

func (n *Navigates) GetRight() *List {
	return n.Get(Right)
}

func (n *Navigates) GetBottom() *List {
	return n.Get(Bottom)
}

func (n *Navigates) Remove(typ NavigateType) bool {
	_, ok := (*n)[typ]
	if ok {
		delete(*n, typ)
	}
	return ok
}
