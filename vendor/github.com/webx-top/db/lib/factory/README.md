# factory
为了方便使用，将原包中的常用功能进行再次封装。支持的目标功能有：

1. 支持主从式读写分离，支持多主多从
2. 支持同时使用多种数据库且每一种数据库都支持读写分离
3. 支持缓存数据结果
4. 便捷的分页查询功能
5. 更加的易于使用

以上。

# 用法

```go
package main
import (
	"fmt"
	"log"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "blog",
	User:     "root",
	Password: "root",
}

type Post struct {
    Id      int     `db:"id,omitempty"`
	UserId  int     `db:"user_id"`
    Title   string  `db:"title"`
    Group   string  `db:"group"`
    Views   int     `db:"views"`
}

type PostCollection struct {
	Post *Post `db:",inline"`
	User *User `db:",inline"`
}

type User struct {
    Id      int     `db:"id,omitempty"`
    Name   string   `db:"name"`
}

func main() {
	database, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	factory.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	factory.SetDebug(true)
	defer factory.Default().CloseAll()

	var posts []*Post

    err = factory.All(factory.NewParam().SetCollection(`post`).SetPage(1).SetSize(10).SetRecv(&posts))
	// 生成SQL：SELECT * FROM `webx_post` LIMIT 10

	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
}
```

## 查询多行数据 (使用All方法)

### 方法 1.
```go
err = factory.All(factory.NewParam().SetCollection(`post`).SetRecv(&posts))
```


也可以附加更多条件（后面介绍的所有方法均支持这种方式）：
```go
 err = factory.NewParam().SetCollection(`post`).SetRecv(&posts).SetArgs(db.Cond{`title LIKE`:`%test%`}).SetMiddleware(func(r db.Result)db.Result{
     return r.OrderBy(`-id`).Group(`group`)
 }).All()
 // 生成SQL：SELECT * FROM `webx_post` WHERE (`title` LIKE "%test%") GROUP BY `group` ORDER BY `id` DESC
 // OrderBy(`-id`) 表示按字段id的值进行降序排列，如果没有前面的负号“-”则按照升序排列
```

### 方法 2.
```go
err = factory.NewParam().SetCollection(`post`).SetRecv(&posts).All()
```

### 方法 3.
```go
post := &Post{}
res := factory.NewParam().SetCollection(`post`).Result()
//defer res.Close() //操作结束后别忘了执行关闭操作
for res.Next(post) {
	log.Printf("%q (ID: %d)\n", post.Title, post.Id)
}
res.Close() //操作结束后别忘了执行关闭操作
```

### 关联查询

#### 手动 Join 查询
```go
m := []*PostCollection{}
err = factory.NewParam().SetCollection(`post`,`a`).SetCols(db.Raw(`a.*`)).AddJoin(`LEFT`, `user`, `b`, `b.id=a.user_id`).Select().All(&m)
```

#### 自动关联查询
```go
type PostWithUser struct {
	*Post
	User *User `db:"-,relation=id:user_id"`
}
m := []*PostWithUser{}
err = factory.NewParam().SetCollection(`post`).Select().All(&m)
```

自动关联查询是通过在结构体 Tag 属性的 `db` 属性中定义 `relation` 参数来实现。  
其完整的格式如下：

relation=`User表列名`:`Post表列名`|`Post表管道函数1`|`Post表管道函数2`,table=`Post表名称`,dbconn=`查询User表采用的数据库连接名称`,columns=`User中的列1`&`User中的列2`,where=`User中的列1`:`User中的列1值`&`User中的列2`:`User中的列2值`,orderby=`User中的列1`&`User中的列2`

其中，__relation=`User表列名`:`Post表列名`__ 是必需的，其余的为可选项。

##### 各个参数的详细说明

* `table`: 指定表名称(在无法通过左侧的字段名自动推导出表名称的情况下使用)
* `dbconn`: 指定查询数据库采用的连接名称(不指定的情况下采用默认连接)
* `columns`: 指定要查询的列(默认为查询所有列)。如果左侧的User字段的数据类型为 `map[string]interface{}` 或者 `echo.H` , 还需要分别对每一个字段指明数据类型，例如：columns=col1:string&col2:int , 要注意的是仅支持指定 github.com/webx-top/echo/param 包中 AsType() 函数支持的数据类型
* `where`: 指定附加条件，条件的列名称和列值用冒号间隔“:”。多个条件之间用“&”分隔。列的值支持采用urlencode的结果值(按照格式`urldecode(编码后的值)`来设置)，列值如果包含前缀`~`则作为原始sql语句(不会经过任何编码处理)
* `orderby`: 指定排序方式。多个字段之间用“&”分隔。支持采用urlencode的结果值(按照格式`urldecode(编码后的值)`来设置)，如果包含前缀`~`则作为原始sql语句
* `offset`: 指定偏移值（小于或等于0的值会忽略）
* `limit`: 指定查询数量（小于或等于0的值会忽略）
* `groupby`: 指定 group by。多个字段之间用“&”分隔。支持采用urlencode的结果值(按照格式`urldecode(编码后的值)`来设置)，如果包含前缀`~`则作为原始sql语句

##### 内置管道函数
管道函数用于定义关联查询前需要对关联值进行预先处理的逻辑。如果管道函数返回 `nil` 值，则取消关联查询，这可以用来自定义判断和控制是否关联查询。
* `neq`: 不等于。用法：neq(col:100) // 即User表中的列col不等于100的时候才进行关联查询
* `eq`: 等于。用法：eq(col:100) // 即User表中的列col等于100的时候才进行关联查询
* `split`: 拆分为切片。用法：split // 即User表中的关联列的值用 strings.Split(value,",") 拆分成切片后再用 IN 查询
* `gtZero`: 大于0。用法：gtZero // 即User表中的关联列的值大于0时才进行关联查询
* `notEmpty`: 非空字符串。用法：notEmpty // 即User表中的关联列的值不为空字符串时才进行关联查询
##### 注册自定义管道函数：
```go
import (
	"github.com/webx-top/db/lib/sqlbuilder"
)

func init(){
	// 方式1：直接注册管道函数
	// 生成管道函数 myPipe
	sqlbuilder.PipeList.Add("myPipe",func(row reflect.Value, val interface{}) interface{}{
		// 你的代码
		return val
	})

	// 方式2：注册管道函数生成器
	// 生成管道函数 myPipeGenerator(params)
	sqlbuilder.PipeGeneratorList.Add("myPipeGenerator",func(params string) sqlbuilder.Pipe {
		// 你的代码
		return func(row reflect.Value, v interface{}) interface{} {
			// 你的代码
			return v
		}
	})
}
```

## 查询分页数据 (使用List方法)

### 方法 1.
```go
var countFn func()int64
countFn, err = factory.List(factory.NewParam().SetCollection(`post`).SetRecv(&posts).SetPage(1).SetSize(10))
```

### 方法 2.
```go
countFn, err = factory.NewParam().SetCollection(`post`).SetRecv(&posts).SetPage(1).SetSize(10).List()
```

## 查询一行数据 (使用One方法)

### 方法 1.
```go
var post Post
err = factory.One(factory.NewParam().SetCollection(`post`).SetRecv(&post))
```

### 方法 2.
```go
var post Post
err = factory.NewParam().SetCollection(`post`).SetRecv(&post).One()
```

## 插入数据 (使用Insert方法)

### 方法 1.
```go
var post Post
post=Post{
    Title:`test title`,
}
err = factory.Insert(factory.NewParam().SetCollection(`post`).SetSend(&post))
```

### 方法 2.
```go
var post Post
post=Post{
    Title:`test title`,
}
err = factory.NewParam().SetCollection(`post`).SetSend(&post).Insert()
```

## 更新数据 (使用Update方法)

### 方法 1.
```go
var post Post
post=Post{
    Title:`test title`,
}
err = factory.Update(factory.NewParam().SetCollection(`post`).SetSend(&post).SetArgs("id",1))
```

### 方法 2.
```go
var post Post
post=Post{
    Title:`test title`,
}
err = factory.NewParam().SetCollection(`post`).SetSend(&post).SetArgs("id",1).Update()
```

## 删除数据 (使用Delete方法)

### 方法 1.
```go
err = factory.Delete(factory.NewParam().SetCollection(`post`).SetArgs("id",1))
```

### 方法 2.
```go
err = factory.NewParam().SetCollection(`post`).SetArgs("id",1).Update()
```

## 使用事务

### 方法 1.
```go
param = factory.NewParam().SetCollection(`post`).SetTxMW(func(t *factory.Transaction) (err error) {
	param := factory.NewParam().SetCollection(`post`).SetSend(map[string]int{
		"views": 1,
	}).SetArgs("id", 1)
	err = t.Update(param)
	// err=fmt.Errorf(`failured`)
	// 当返回 nil 时，自动执行Commit，否则自动执行Rollback
	return
})
factory.Tx(param)
```

### 方法 2.
```go
param = factory.NewParam().Begin().SetCollection(`post`)
err:=param.SetSend(map[string]int{"views": 1}).SetArgs("id", 1).Update()
if err!=nil {
    param.End(err)
    return
}
err=factory.NewParam().TransFrom(param).SetCollection(`post`).SetSend(map[string]int{"views": 2}).SetArgs("id", 1).Update()
if err!=nil {
    param.End(err)
    return
}
err=factory.NewParam().TransFrom(param).SetCollection(`post`).SetSend(map[string]int{"views": 3}).SetArgs("id", 1).Update()
param.End(err)
```

# 自动生成数据表的结构体(struct)
进入目录`github.com/webx-top/db/cmd/dbgenerator`执行命令
```go
go build -o dbgenerator.exe
dbgenerator.exe -u <数据库用户名> -p <数据库密码> -h <数据库主机名> -e <数据库类型> -d <数据库名> -o <文件保存目录> -pre <数据表前缀> -pkg <生成的包名>
```

支持的参数：

* -u <数据库用户名> 默认为`root`
* -p <数据库密码> 默认为空
* -h <数据库主机名> 默认为`localhost`
* -e <数据库类型> 默认为`mysql`
* -d <数据库名> 默认为`blog`
* -o <文件保存目录> 默认为`dbschema`
* -pre <数据表前缀> 默认为空
* -pkg <生成的包名> 默认为`dbschema`
* -autoTime <自动生成时间戳的字段> 默认为`update(*:updated)/insert(*:created)`：
  
  >  即在更新任意数据表时，自动设置表的updated字段；
  >  在新增数据到任意数据表时，自动设置表的created字段。  
  >  括号内的格式：`<表1>:<字段1>,<字段2>,<...字段N>;<表2>:<字段1>,<字段2>,<...字段N>`

本命令会自动生成各个表的结构体和所有表的相关信息

生成的文件范例：[blog结构体文件](https://github.com/webx-top/db/tree/master/_tools/test/dbschema)

每一个生成的结构体中都自带了以下常用方法便于我们使用：

* 获取事务 `Trans() *factory.Transaction`
* 设置事务 `Use(trans *factory.Transaction) factory.Model`
* 新参数对象 `NewParam() *factory.Param` 
* 设置默认参数对象 `SetParam(param *factory.Param) factory.Model` 
* 获取参数对象 `Param() *factory.Param` (如果有默认参数对象则使用默认，否则新建参数对象)
* 复制列表数据结果集 `Objects() []*结构体名` 
* 新建列表数据结果集 `NewObjects() *[]*结构体名` 
* 查询一行 `Get(mw func(db.Result) db.Result, args ...interface{}) error`
* 分页查询 `List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)`
* 根据偏移量查询 `ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)`
* 添加数据 `Insert() (interface{}, error)`
* 修改数据 `Update(mw func(db.Result) db.Result, args ...interface{}) error`
* 删除数据 `Delete(mw func(db.Result) db.Result, args ...interface{}) error`
* 统计行数 `Count(mw func(db.Result) db.Result, args ...interface{}) error`

> 如果数据库中的字段含有注释，并且注释内容是以反引号``` `...` ```这样的样式开头，  
> 那么反引号内的内容会作为是否在该表结构体字段上的db标签添加`omitempty`和`pk`的依据。  
> 例如：数据表user的字段username注释为``` `omitempty`用户名 ```，则该结构体就会生成这样：
```go
...
type User struct {
	Username string `db:"username,omitempty" bson:"username,omitempty" comment:"用户名" json:"username" xml:"username"`
}
```
> 又例如：数据表user的字段username注释为``` `omitempty,pk`用户名 ```，则该结构体就会生成这样：
```go
...
type User struct {
	Username string `db:"username,omitempty,pk" bson:"username,omitempty,pk" comment:"用户名" json:"username" xml:"username"`
}
```

我们还可以根据生成的数据表信息来验证表或字段的类型和合法性，比如：

* 验证表是否存在 `factory.ValidTable(tableName)`
* 验证表中的字段是否存在 `factory.ValidField(tableName,fieldName)`
* 获取某个表中某个字段的信息  `factory.Fields[tableName][fieldName]`


