package containerfs

import "time"

type FileInfo struct {
	FileMode string
	Name     string
	ModTime  time.Time
	Size     int64
	User     string
	Group    string
	IsDir    bool
}
