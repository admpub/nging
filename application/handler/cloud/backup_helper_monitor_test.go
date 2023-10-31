package cloud

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestGrowSliceSize(t *testing.T) {
	oldValParts := []string{`md5`}
	if len(oldValParts) < 5 {
		temp := make([]string, 5)
		copy(temp, oldValParts)
		oldValParts = temp
	}
	assert.Equal(t, `md5`, oldValParts[0])
	assert.Equal(t, 5, len(oldValParts))
}

func TestMonitorBackup(t *testing.T) {
	dir := `./testdata/backup`
	err := com.MkdirAll(dir, os.ModePerm)
	assert.NoError(t, err)
	cfg := dbschema.NgingCloudBackup{
		Id:            100000,
		SourcePath:    filepath.Join(dir),
		DestPath:      `/backup`,
		DestStorage:   0,
		StorageEngine: `mock`,
		//Delay:             5,
		WaitFillCompleted: `Y`,
	}
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := monitorBackupStart(cfg, false)
		assert.NoError(t, err)
	}()

	go func() {
		defer wg.Done()
		filePath := filepath.Join(dir, `test.txt`)
		fp, err := os.Create(filePath)
		assert.NoError(t, err)
		for i := 0; i < 1000; i++ {
			fp.WriteString(fmt.Sprintf("~~~~~~~~~~~~~~~~~~~~~~>%d\n", i))
			fmt.Printf("write:%d\n", i)
			//time.Sleep(time.Millisecond * 10)
		}
		fp.Close()
		time.Sleep(30 * time.Second)
		monitorBackupStop(cfg.Id)
		cloudbackup.LevelDB().RemoveDB(cfg.Id)
	}()
	wg.Wait()
	result := buf.String()
	assert.Equal(t, 1, strings.Count(result, `StorageMock: Put`))
}
