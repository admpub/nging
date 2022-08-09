package ip2region

import (
	"sync"

	"github.com/admpub/ip2region/v2/binding/golang/xdb"
)

type Ip2Region struct {
	dbFile string
	dbBuff []byte
	mu     sync.RWMutex
}

func New(path string) (*Ip2Region, error) {
	cBuff, err := xdb.LoadContentFromFile(path)
	if err != nil {
		return nil, err
	}

	//searcher, err := xdb.NewWithFileOnly(path)
	return &Ip2Region{
		dbFile: path,
		dbBuff: cBuff,
	}, nil
}

func (a *Ip2Region) DBBuff() []byte {
	a.mu.RLock()
	buff := a.dbBuff
	a.mu.RUnlock()
	return buff
}

func (a *Ip2Region) Close() error {
	return nil
}

func (a *Ip2Region) MemorySearch(ipStr string) (ipInfo IpInfo, err error) {
	searcher := searcherPool.Get().(*xdb.Searcher)
	searcher.SetContentBuff(a.DBBuff())
	var result string
	result, err = searcher.SearchByStr(ipStr)
	searcher.Close()
	searcherPool.Put(searcher)
	if err != nil {
		return
	}
	ipInfo = getIpInfo(result)
	return
}
