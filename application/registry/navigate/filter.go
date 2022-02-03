package navigate

import "github.com/webx-top/echo"

type Checker interface {
	Check(echo.Context, string) bool
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
func (r *Filter) FilterNavigate(ctx echo.Context, navList *List) List {
	var result List
	if navList == nil {
		return result
	}
	for _, nav := range *navList {
		children := r.filterNavigateChidren(ctx, nav.Action, nav, nav.Children)
		if children == nil {
			continue
		}
		navCopy := *nav
		navCopy.Children = children
		result = append(result, &navCopy)
	}
	return result
}

func (r *Filter) filterNavigateChidren(ctx echo.Context, permPath string, parent *Item, children *List) *List {
	if children == nil {
		if !parent.Unlimited && !r.Check(ctx, permPath) {
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
		list := r.filterNavigateChidren(ctx, perm, child, child.Children)
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
	if len(newChildren) == 0 && !r.Check(ctx, permPath) {
		return nil
	}
	return &newChildren
}
