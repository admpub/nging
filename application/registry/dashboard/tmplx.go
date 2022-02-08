package dashboard

import "github.com/webx-top/echo"

func NewTmplx(tmpl string, handle ...func(echo.Context) error) *Tmplx {
	var content func(echo.Context) error
	if len(handle) > 0 {
		content = handle[0]
	}
	return &Tmplx{Tmpl: tmpl, content: content}
}

type Tmplx struct {
	Tmpl    string //模板文件
	content func(echo.Context) error
}

func (c *Tmplx) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *Tmplx) SetContentGenerator(content func(echo.Context) error) *Tmplx {
	c.content = content
	return c
}

func (c *Tmplx) SetTmpl(tmpl string) *Tmplx {
	c.Tmpl = tmpl
	return c
}

type Tmplxs []*Tmplx

func (c *Tmplxs) Ready(ctx echo.Context) error {
	for _, blk := range *c {
		if blk != nil {
			if err := blk.Ready(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove 删除元素
func (c *Tmplxs) Remove(index int) {
	if index < 0 {
		*c = (*c)[0:0]
		return
	}
	size := c.Size()
	if size > index {
		if size > index+1 {
			*c = append((*c)[0:index], (*c)[index+1:]...)
		} else {
			*c = (*c)[0:index]
		}
	}
}

func (c *Tmplxs) Add(index int, list ...*Tmplx) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		list = append(list, (*c)[index])
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Add(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

// Set 设置元素
func (c *Tmplxs) Set(index int, list ...*Tmplx) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Set(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

func (c *Tmplxs) Size() int {
	return len(*c)
}

func (c *Tmplxs) Search(cb func(*Tmplx) bool) int {
	for index, footer := range *c {
		if cb(footer) {
			return index
		}
	}
	return -1
}

func (c *Tmplxs) FindTmpl(tmpl string) int {
	return c.Search(func(footer *Tmplx) bool {
		return footer.Tmpl == tmpl
	})
}

func (c *Tmplxs) RemoveByTmpl(tmpl string) {
	index := c.FindTmpl(tmpl)
	if index > -1 {
		c.Remove(index)
	}
}
