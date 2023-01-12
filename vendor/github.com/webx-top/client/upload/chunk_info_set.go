package upload

import (
	"math"
	"net/url"
	"regexp"

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

var reContentRange = regexp.MustCompile(`bytes[ ]+([0-9]+)-([0-9]+)/([0-9]+)`)

func (c *ChunkInfo) ParseHeader(formValue func(string) string, header func(string) string) bool {
	chunkSize := formValue(`chunkSize`)
	if len(chunkSize) == 0 {
		chunkSize = header(`X-Chunk-Size`)
	}
	if len(chunkSize) == 0 {
		return false
	}
	c.Set(`fileChunkBytes`, chunkSize)
	if c.FileChunkBytes < 1 {
		return false
	}
	contentRange := header(`Content-Range`)
	return c.parseHeader(contentRange)
}

// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Range
func (c *ChunkInfo) parseHeader(contentRange string) bool {
	if len(contentRange) == 0 {
		return false
	}
	matches := reContentRange.FindAllStringSubmatch(contentRange, 1)
	if len(matches) == 0 {
		return false
	}
	c.Set(`chunkOffsetBytes`, matches[0][1])
	c.Set(`fileTotalBytes`, matches[0][3])
	c.ChunkEndBytes = param.AsUint64(matches[0][2])
	if c.ChunkEndBytes <= c.ChunkOffsetBytes {
		return false
	}
	if c.CurrentSize < 1 {
		c.CurrentSize = (c.ChunkEndBytes - c.ChunkOffsetBytes) + 1
	}
	if c.FileChunkBytes < 1 && c.ChunkEndBytes+1 != c.FileTotalBytes {
		c.FileChunkBytes = c.CurrentSize
	}
	if c.FileChunkBytes > 0 {
		c.FileTotalChunks = TotalChunks(c.FileTotalBytes, c.FileChunkBytes)
		c.ChunkIndex = (c.ChunkEndBytes + 1) / c.FileChunkBytes
	}
	return true
}

func TotalChunks(totalBytes uint64, chunkBytes uint64) uint64 {
	return uint64(math.Ceil(float64(totalBytes) / float64(chunkBytes)))
}

func (c *ChunkInfo) Set(key string, value interface{}) {
	switch key {
	case `fileUUID`:
		c.FileUUID = param.AsString(value)
	case `chunkIndex`: // * required
		c.ChunkIndex = param.AsUint64(value)
	case `fileTotalBytes`: // * required
		c.FileTotalBytes = param.AsUint64(value)
	case `fileChunkBytes`: // * required
		c.FileChunkBytes = param.AsUint64(value)
	case `fileTotalChunks`: // * required
		c.FileTotalChunks = param.AsUint64(value)
	case `chunkOffsetBytes`:
		c.ChunkOffsetBytes = param.AsUint64(value)
	case `chunkEndBytes`:
		c.ChunkEndBytes = param.AsUint64(value)
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
