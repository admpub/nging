package goseaweedfs

import (
	"encoding/json"
)

// ChunkInfo chunk information. According to https://github.com/chrislusf/seaweedfs/wiki/Large-File-Handling.
type ChunkInfo struct {
	Fid    string `json:"fid"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
}

// ChunkManifest chunk manifest. According to https://github.com/chrislusf/seaweedfs/wiki/Large-File-Handling.
type ChunkManifest struct {
	Name   string       `json:"name,omitempty"`
	Mime   string       `json:"mime,omitempty"`
	Size   int64        `json:"size,omitempty"`
	Chunks []*ChunkInfo `json:"chunks,omitempty"`
}

// Marshal marshal whole chunk manifest
func (c *ChunkManifest) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
