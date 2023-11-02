package cloudbackup

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/once"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var (
	levelDBPool *dbPool
	levelDBOnce once.Once
	LevelDBDir  = `data/cache/backup-db`
)

func LevelDB() *dbPool {
	levelDBOnce.Do(initLevelDB)
	return levelDBPool
}

func initLevelDB() {
	levelDBPool = NewLevelDBPool()
}

func NewLevelDBPool() *dbPool {
	return &dbPool{mp: map[uint]*leveldb.DB{}}
}

type dbPool struct {
	mu sync.RWMutex
	mp map[uint]*leveldb.DB
}

func (t *dbPool) OpenDB(taskId uint) (*leveldb.DB, error) {
	t.mu.RLock()
	db := t.mp[taskId]
	t.mu.RUnlock()

	if db == nil {
		t.mu.Lock()
		defer t.mu.Unlock()
		var err error
		db, err = openLevelDB(taskId)
		if err != nil {
			return nil, err
		}
		t.mp[taskId] = db
	}
	return db, nil
}

func openLevelDB(taskId uint) (*leveldb.DB, error) {
	idKey := com.String(taskId)
	dbDir := filepath.Join(echo.Wd(), LevelDBDir)
	err := com.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(dbDir, idKey)
	return leveldb.OpenFile(dbFile, nil)
}

func removeLevelDB(taskId uint) error {
	idKey := com.String(taskId)
	dbFile := filepath.Join(echo.Wd(), LevelDBDir, idKey)
	if !com.FileExists(dbFile) {
		return nil
	}
	return os.RemoveAll(dbFile)
}

func (t *dbPool) CloseDB(taskId uint) {
	t.mu.Lock()
	if db, ok := t.mp[taskId]; ok {
		db.Close()
		delete(t.mp, taskId)
	}
	t.mu.Unlock()
}

func (t *dbPool) RemoveDB(taskId uint) error {
	t.mu.Lock()
	if db, ok := t.mp[taskId]; ok {
		db.Close()
		delete(t.mp, taskId)
	}
	err := removeLevelDB(taskId)
	t.mu.Unlock()
	return err
}

func (t *dbPool) CloseAllDB() {
	t.mu.Lock()
	for _, db := range t.mp {
		db.Close()
	}
	t.mu.Unlock()
}

func ParseDBValue(val []byte) (md5 string, startTs, endTs, fileModifyTs, fileSize int64) {
	parts := strings.Split(com.Bytes2str(val), `||`)
	md5 = parts[0]
	if len(parts) > 1 {
		com.SliceExtractCallback(parts[1:], func(v string) int64 {
			return param.AsInt64(v)
		}, &startTs, &endTs, &fileModifyTs, &fileSize)
	}
	return
}
