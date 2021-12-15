# mail
Golang SMTP电子邮件包

# 安装
```go
go get github.com/admpub/mail
```

# 例子
```go
import (
  "fmt"
  "github.com/admpub/mail"
  "os"
)

func main() {
  conf := &mail.SMTPConfig{
      Username: "admpub",
      Password: "",
      Host:     "smtp.admpub.com",
      Port:     587,
      Secure:   "SSL",
  }
  c := mail.NewSMTPClient(conf)
  m := mail.NewMail()
  m.AddTo("hello@admpub.com") //或 "老弟 <hello@admpub.com>"
  m.AddFrom("hank@admpub.com") //或 "老哥 <hank@admpub.com>"
  m.AddSubject("Testing")
  m.AddText("Some text :)")
  filepath, _ := os.Getwd()
  m.AddAttachment(filepath + "/mail.go")
  if e := c.Send(m); e != nil {
    fmt.Println(e)
  } else {
    fmt.Println("发送成功")
  }
}
```