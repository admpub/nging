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

func (n *Navigates) AddItems(typ NavigateType, index int, items ...*Item) {
	nav := n.Get(typ)
	if nav != nil {
		nav.Add(index, items...)
	}
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
