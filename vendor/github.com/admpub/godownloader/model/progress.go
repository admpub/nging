package model

import "time"

type DownloadProgress struct {
	From          int64 // 分片起始字节位置
	To            int64 // 分片终止字节位置
	Pos           int64 // 已下载占整个文件的位置 From<=Pos<=To
	BytesInSecond int64
	Speed         int64
	Lsmt          time.Time
	IsPartial     bool // 是否分片下载
}

func (dp DownloadProgress) IsCompleted() bool {
	return dp.Pos >= dp.To
}

func (dp *DownloadProgress) ResetProgress() {
	dp.Pos = dp.From
	dp.Speed = 0
}
