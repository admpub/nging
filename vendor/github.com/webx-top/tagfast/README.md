tagfast
=======

golang：优化结构体字段标签读取（只解析一次，以后都从缓存中读取。）

> 在使用标准库中的`typ.Field(i).Tag.Get("tag1")`时，每次都要解析一次，效率不高。  
> 使用本包的`tagfast.Value(typ, typ.Field(i), "tag1")`来代替`typ.Field(i).Tag.Get("tag1")`即可。

## 缓存更复杂的解析结果
```go
//`form:"checked(true);required(true)"`
tag, parse := tagfast.Tag(typ, typ.Field(i), "form")
fmt.Println(tag) // 输出：checked(true);required(false)

parser := func() interface{} {
    options := map[string]string{}
    for _, item := range strings.Split(tag, ";"){
        p := strings.IndexOf(item, "(")
        key := item[0:p]
        val := item[p+1:strings.IndexOf(item, ")")]
        options[key] = val
    }
    return options
}


if value, ok := parse.Parsed("form", parser).(map[string]string); ok {
    fmt.Println(value["checked"]) // 输出：true
    fmt.Println(value["required"]) // 输出：false
}
```


注意事项
=======
> 不要在同一个包内定义相同名称的结构体。  
> 由于是按照“包名称+结构体名称”来缓存Tag，第二个结构体如果跟第一个结构体同名，  
> 那么它们只会被当做同一个结构体来获取Tag，从而导致第二个具有相同名称的结构体获取Tag不正确。