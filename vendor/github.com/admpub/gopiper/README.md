# gopiper

[TOC]

## 介绍

gopiper提供一种通过配置规则的方式将网页源码【网页源码类型可以为html/json/text】提取结果为json序列化的数据格式。

比如豆瓣电影的一个网页[https://movie.douban.com/subject/26580232/]

可以将页面提取成一个json对象

```json
{
	"name": "看不见的客人",
	"pic": "https://img3.doubanio.com/view/movie_poster_cover/lpst/public/p2498971355.webp",
	"score": 8.7,
	"director": "奥里奥尔·保罗",
	"actor": ["马里奥·卡萨斯","阿娜·瓦格纳","何塞·科罗纳多","巴巴拉·莱涅","弗兰塞斯克·奥雷利亚"]
}
```

## 规则描述

规则被描述成一个可嵌套的JSON结构（包含子结构），json结构如下：

```json
{
	"name": "结果名",
	"selector": "节点选择器",
	"type": "规则类型",
	"filter": "过滤处理函数",
	"subitem": [
	    //子规则嵌套, 只有规则类型为map或array
	],
}
```

其Go表示的Struct如下：

```go
type PipeItem struct {
	Name     string     `json:"name,omitempty"`       // 结果名称[map子结构有效]
	Selector string     `json:"selector,omitempty"`   // 节点选择器
	Type     string     `json:"type"`                 // 规则类型
	Filter   string     `json:"filter,omitempty"`     // 过滤器或结果函数处理
	SubItem  []PipeItem `json:"subitem,omitempty"`    // 嵌套子结构
}
```

### 规则类型

规则类型主要可分为三种，一种是map类型（结果值为json对象），一种是array类型(结果值为json数组），另外一种为单值类型（字符串、数值等）

#### map类型

#### array类型

#### 值类型


### 选择器

### 过滤器函数

### 规则案例

豆瓣电影页面提取规则: http://movie.douban.com/subject/25850640/ 

```json
{
	"type": "map",
	"selector": "",
	"subitem": [
		{
			"type": "string",
			"selector": "title",
			"name": "name",
			"filter": "trimspace|replace((豆瓣))|trim( )"
		},
		{
			"type": "string",
			"selector": "#content .gtleft a.bn-sharing@attr[data-type]",
			"name": "fenlei"
		},
		{
			"type": "string",
			"selector": "#content .gtleft a.bn-sharing@attr[data-pic]",
			"name": "thumbnail"
		},
		{
			"type": "string-array",
			"selector": "#info span.attrs a[rel=v\\:directedBy]",
			"name": "direct"
		},
		{
			"type": "string-array",
			"selector": "#info span a[rel=v\\:starring]",
			"name": "starring"
		},
		{
			"type": "string-array",
			"selector": "#info span[property=v\\:genre]",
			"name": "type"
		},
		{
			"type": "string-array",
			"selector": "#related-pic .related-pic-bd a:not(.related-pic-video) img@attr[src]",
			"name": "imgs",
			"filter": "join($)|replace(albumicon,photo)|split($)"
		},
		{
			"type": "string-array",
			"selector": "#info span[property=v\\:initialReleaseDate]",
			"name": "releasetime"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">单集片长:</span> ([\\w\\W]+?)<br/>",
			"name": "longtime"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">制片国家/地区:</span> ([\\w\\W]+?)<br/>",
			"name": "country",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">语言:</span> ([\\w\\W]+?)<br/>",
			"name": "language",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "int",
			"selector": "regexp:<span class=\"pl\">集数:</span> (\\d+)<br/>",
			"name": "episode"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">又名:</span> ([\\w\\W]+?)<br/>",
			"name": "alias",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "string",
			"selector": "#link-report span.hidden, #link-report span[property=v\\:summary]|last",
			"name": "brief",
			"filter": "trimspace|split(\n)|trimspace|wraphtml(p)|join"
		},
		{
			"type": "float",
			"selector": "#interest_sectl .rating_num",
			"name": "score"
		},
		{
			"type": "string",
			"selector": "#content h1 span.year",
			"name": "year",
			"filter": "replace(()|replace())|intval"
		},
		{
			"type": "string",
			"selector": "#comments-section > .mod-hd h2 a",
			"name": "comment",
			"filter": "replace(全部)|replace(条)|trimspace|intval"
		}
	]
}

```

## 用法


