```go

import (
	"log"
	"testing"
)

type User struct {
	Name string
	ID   int
}

func TestScan(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	m := map[string]interface{}{}
	m["Name"] = "test"
	m["ID"] = 100
	user := new(User)
	err := Scan(user, m)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("==>", user.ID, user.Name)
}

```