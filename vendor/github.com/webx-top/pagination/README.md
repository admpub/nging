# 分页组件

支持两种分页模式：

1. 页码分页模式
2. 偏移值分页模式

## 分页方式

### 按照页码分页

```go
import (
    "github.com/webx-top/echo"
    "github.com/webx-top/pagination"
)

func index(ctx echo.Context) error {
    p := pagination.New(ctx)
    totalRows := 1000
    page := ctx.Formx("page").Int()
    if page < 1 {
        page = 1
    }
    links := 10
    limit := ctx.Formx("size").Int()
    if limit < 0 {
        limit = 20
    }
    p.SetAll(`tmpl/pagination`, totalRows, page, links, limit)
    // 自动根据当前网址生成分页网址
    p.SetURL(nil, `_pjax`)
    /* 或者手动指定:
    q := ctx.Request().URL().Query()
    q.Del(`page`)
    q.Del(`rows`)
    q.Del(`_pjax`)
    p.SetURL(ctx.Request().URL().Path() + `?` + q.Encode() + `&page={page}&rows={rows}`)
    // p.SetPage(page).SetRows(totalRows))
    */
    ctx.Set(`pagination`, p)
    return ctx.Render(`list`, nil)
}

```

### 按照偏移值分页

```go
import (
    "github.com/webx-top/echo"
    "github.com/webx-top/pagination"
)

func index(ctx echo.Context) error {
    p := pagination.New(ctx)
    offset := ctx.Form(`offset`)
    prevOffset := ctx.Form(`prev`)

    var nextOffset string
    // TODO：获取nextOffset
    // nextOffset = "你的下一页偏移值"

    // 自动根据当前网址生成分页网址
    p.SetURL(nil, `_pjax`)
    /* 或者手动指定:
    q := ctx.Request().URL().Query()
    q.Del(`offset`)
    q.Del(`prev`)
    q.Del(`_pjax`)
    p.SetURL(ctx.Request().URL().Path() + `?` + q.Encode() + `&offset={curr}&prev={prev}`)
    */
    p.SetPosition(prevOffset, nextOffset, offset)
    ctx.Set(`pagination`, p)
    return ctx.Render(`list`, nil)
}

```

## 模板中输出分页

list.html

```html
{{Stored.pagination.Render}}
<!-- 或者指定分页链接的模板：{{Stored.pagination.Render "pagination_cursor"}} -->
```
