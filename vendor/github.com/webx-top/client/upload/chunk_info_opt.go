package upload

type ChunkInfoOpter func(*ChunkInfo)

func OptChunkOffsetBytes(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.ChunkOffsetBytes = m
	}
}

func OptChunkEndBytes(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.ChunkEndBytes = m
	}
}

func OptChunkIndex(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.ChunkIndex = m
	}
}

func OptChunkCurrentSize(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.CurrentSize = m
	}
}

func OptChunkFileTotalBytes(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FileTotalBytes = m
	}
}

func OptChunkFileTotalChunks(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FileTotalChunks = m
	}
}

func OptChunkFileChunkBytes(m uint64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FileChunkBytes = m
	}
}

func OptChunkFileUUID(m string) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FileUUID = m
	}
}

func OptChunkFileName(m string) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FileName = m
	}
}

func OptChunkInfoMapping(m map[string]string) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.Mapping = m
	}
}

func OptChunkFormField(m string) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.FormField = m
	}
}

// OptChunkSpeedBps 速度限制:每秒接收字节数
func OptChunkSpeedBps(speedBps int64) ChunkInfoOpter {
	return func(c *ChunkInfo) {
		c.SpeedBps = speedBps
	}
}
