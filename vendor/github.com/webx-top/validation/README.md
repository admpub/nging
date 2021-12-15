validation
==============

validation is a form validation for a data validation and error collecting using Go.

## Installation and tests

Install:
```go
go get github.com/webx-top/validation
```

Test:
```go
go test github.com/webx-top/validation
```

## Example

Direct Use:
```go
import (
	"github.com/webx-top/validation"
	"log"
)
type User struct {
	Name string
	Age int
}
func main() {
	u := User{"man", 40}
	valid := validation.New()
	valid.Required(u.Name, "name")
	valid.MaxSize(u.Name, 15, "nameMax")
	valid.Range(u.Age, 0, 140, "age")
	if valid.HasError() {
		// validation does not pass
		// print invalid message
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}
	// or use like this
	if v := valid.Max(u.Age, 140); !v.Ok {
		log.Println(v.Error.Key, v.Error.Message)
	}
}
```
Struct Tag Use:
```go
import (
	"github.com/webx-top/validation"
)
// validation function follow with "valid" tag
// functions divide with ";"
// parameters in parentheses "()" and divide with ","
// Match function's pattern string must in "//"
type User struct {
	Id   int
	Name string `valid:"required;match(/^(test)?\\w*@;com$/)"`
	Age  int    `valid:"required;range(1, 140)"`
}
type Profile struct {
	Id   int
	Email string `valid:"required;match(/^\\w+@coscms\\.com$/)"`
	Addr  string `valid:"required"`
}
type NotValid struct {
	A string
	B string
}
type Group struct {
	Id   int
	User
	*Profile
	NotValid `valid:"-"` //valid标签设为“-”，意味着跳过此项不查询其成员
}
func main() {
	valid := validation.New()
	u := User{Name: "test", Age: 40}
	b, err := valid.Valid(u) //检查所有字段
	//b, err := valid.Valid(u, "Name", "Age") //检查指定字段：Name和Age
	if err != nil {
		// handle error
	}
	if !b {
		// validation does not pass
		// blabla...
	}
	valid.Clear()
	u := Group{
		User:           User{Name: "test", Age: 40},
		Profile:        &Profile{Email:"test@coscms.com",Addr:"address"},
		NotValid:       NotValid{},
	}
	b, err := valid.Valid(u) //检查所有字段
	//b, err := valid.Valid(u, "User.Name", "Profile.Email") //检查指定字段
	if err != nil {
		// handle error
	}
	if !b {
		// validation does not pass
		// blabla...
	}
}
```
Struct Tag Functions:
```
-
required
min(min int)
max(max int)
range(min, max int)
minSize(min int)
maxSize(max int)
length(length int)
alpha
numeric
alphaNumeric
match(pattern string)
alphaDash
email
ip
base64
mobile
tel
phone
zipCode
```

## LICENSE

BSD License http://creativecommons.org/licenses/BSD/
