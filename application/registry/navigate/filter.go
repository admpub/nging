package navigate

type Checker interface {
	Check(string) bool
}

func NewFilter(checker Checker) *Filter {
	return &Filter{
		Checker: checker,
	}
}

type Filter struct {
	Checker
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *Filter) FilterNavigate(navList *List) List {
	var result List
	if navList == nil {
		return result
	}
	for _, nav := range *navList {
		children := r.filterNavigateChidren(nav.Action, nav, nav.Children)
		if children == nil {
			continue
		}
		navCopy := *nav
		navCopy.Children = children
		result = append(result, &navCopy)
	}
	return result
}

func (r *Filter) filterNavigateChidren(permPath string, parent *Item, children *List) *List {
	if children == nil {
		if !parent.Unlimited && !r.Check(permPath) {
			return nil
		}
		return &List{}
	}
	newChildren := List{}
	for _, child := range *children {
		var perm string
		if len(child.Action) > 0 {
			perm = permPath + `/` + child.Action
		} else {
			perm = permPath
		}
		list := r.filterNavigateChidren(perm, child, child.Children)
		if list == nil {
			continue
		}
		childCopy := *child
		childCopy.Children = list
		newChildren = append(newChildren, &childCopy)
	}
	if parent.Unlimited {
		return &newChildren
	}
	if len(newChildren) == 0 && !r.Check(permPath) {
		return nil
	}
	return &newChildren
}
