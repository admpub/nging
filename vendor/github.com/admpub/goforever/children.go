package goforever

import (
	"encoding/json"
	"log"
)

//Children Child processes.
type Children map[string]*Process

//String Stringify
func (c Children) String() string {
	js, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(js)
}

//Keys Get child processes names.
func (c Children) Keys() []string {
	keys := []string{}
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

//Get child process.
func (c Children) Get(key string) *Process {
	if v, ok := c[key]; ok {
		return v
	}
	return nil
}

func (c Children) Stop(names ...string) {
	if len(names) < 1 {
		for name, p := range c {
			p.Stop()
			delete(c, name)
		}
		return
	}
	name := names[0]
	p := c.Get(name)
	p.Stop()
	delete(c, name)
}
