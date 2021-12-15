package upload

import (
	"net/url"

	"github.com/webx-top/echo/param"
)

func (c *ChunkInfo) getFormField(field string) string {
	if c.Mapping == nil {
		return field
	}
	if v, y := c.Mapping[field]; y {
		return v
	}
	return field
}

func (c *ChunkInfo) BatchSet(m param.Store) {
	for key, val := range m {
		c.Set(key, val)
	}
}

func (c *ChunkInfo) Set(key string, value interface{}) {
	switch key {
	case `fileUUID`:
		c.FileUUID = param.AsString(value)
	case `chunkIndex`:
		c.ChunkIndex = param.AsUint64(value)
	case `fileTotalBytes`:
		c.FileTotalBytes = param.AsUint64(value)
	case `fileChunkBytes`:
		c.FileChunkBytes = param.AsUint64(value)
	case `fileTotalChunks`:
		c.FileTotalChunks = param.AsUint64(value)
	case `chunkOffsetBytes`:
		c.ChunkOffsetBytes = param.AsUint64(value)
	}
}

func (c *ChunkInfo) BatchSetByURLValues(m url.Values) {
	c.CallbackBatchSet(m.Get)
}

func (c *ChunkInfo) CallbackBatchSet(cb func(string) string) {
	c.FileUUID = cb(c.getFormField(`fileUUID`))
	c.ChunkIndex = param.AsUint64(cb(c.getFormField(`chunkIndex`)))
	c.FileTotalBytes = param.AsUint64(cb(c.getFormField(`fileTotalBytes`)))
	c.FileChunkBytes = param.AsUint64(cb(c.getFormField(`fileChunkBytes`)))
	c.FileTotalChunks = param.AsUint64(cb(c.getFormField(`fileTotalChunks`)))
	c.ChunkOffsetBytes = param.AsUint64(cb(c.getFormField(`chunkOffsetBytes`)))
}
