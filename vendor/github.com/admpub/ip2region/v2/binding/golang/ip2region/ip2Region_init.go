package ip2region

import (
	"sync"

	"github.com/admpub/ip2region/v2/binding/golang/xdb"
)

func (a *Ip2Region) Reload(newPath ...string) error {
	path := a.dbFile
	if len(newPath) > 0 && len(newPath[0]) > 0 {
		path = newPath[0]
	}
	cBuff, err := xdb.LoadContentFromFile(path)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.dbBuff = cBuff
	a.mu.Unlock()
	return nil
}

var searcherPool = sync.Pool{
	New: func() interface{} {
		return &xdb.Searcher{}
	},
}
