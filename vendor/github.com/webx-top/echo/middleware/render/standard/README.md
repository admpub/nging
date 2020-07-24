# 模板引擎

## 特点

1. 支持继承
2. 支持包含子模板
3. 支持golang原生模板语法（详细用法可参阅[golang模板语法简明教程](http://www.admpub.com/blog/post-221.html)）
4. 自动注册模板函数 `hasBlock(blocks ...string) bool` 和 `hasAnyBlock(blocks ...string) bool`

    * `hasBlock(blocks ...string) bool` - 是否在扩展模板中设置了指定名称的Block
    * `hasAnyBlock(blocks ...string) bool`  - 是否在扩展模板中设置了指定名称中的任意一个Block

5. 支持多重继承

## 模板继承

用于模板继承的标签有：Block、Extend、Super
例如，有以下两个模板：

1. layout.html：

```html
{{Block "title"}}-- powered by webx{{/Block}}
{{Block "body"}}内容{{/Block}}
```

2. index.html：

```html	
{{Extend "layout"}}
{{Block "title"}}首页 {{Super}}{{/Block}}
{{Block "body"}}这是一个演示{{/Block}}
```

渲染模板index.html将会输出:

```html
首页 -- powered by webx
这是一个演示
```

注意：

> 1. Super标签只能在扩展模板（含Extend标签的模板）的Block标签内使用。
> 2. Extend标签 必须放在页面内容的起始位置才有效

因为最新增加了对多重继承的支持，所以，现在我们还可以创建一个模板`new.html`用来继承上面的`index.html`，比如`new.html`的内容为：

```html
{{Extend "index"}}
{{Block "body"}}这是一个多重继承的演示{{/Block}}
```

渲染这个新模板将会输出：

```html
首页 -- powered by webx
这是一个多重继承的演示
```

也就是说这个新模板具有这样的继承关系：new.html -> index.html -> layout.html (目前限制为最多不超过10级)

		
## 包含子模板
	
	例如，有以下两个模板：
	footer.html:
		
		www.webx.top
	
	index.html:
	
		前面的一些内容
		{{Include "footer"}}
		后面的一些内容
		
	渲染模板index.html将会输出:
	
		前面的一些内容
		www.webx.top
		后面的一些内容
		
	也可以在循环中包含子模板，例如：
	
		{{range .list}}
		{{Include "footer"}}
		{{end}}
		
因为本模板引擎缓存了模板对象，所以它并不会多次读取模板内容，在循环体内也能高效的工作。
	
Include标签也能在Block标签内部使用，例如：
	
		{{Block "body"}}
		这是一个演示
		{{Include "footer"}}
		{{/Block}}

另外，Include标签也支持嵌套。


点此查看[完整例子](https://github.com/webx-top/echo/tree/master/middleware/render/example)
