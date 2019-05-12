package model

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"sort"
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

// UnGzipData unzip data in form of stream/bytes
func UnGzipData(input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(input)
	r, _ := gzip.NewReader(buf)
	defer r.Close()
	output, err := ioutil.ReadAll(r)
	return output, err
}

// LoadChunkManifest load chunk manifest from bytes, support gzipped bytes
func LoadChunkManifest(buffer []byte, isGzipped bool) (*ChunkManifest, error) {
	if isGzipped {
		var err error
		if buffer, err = UnGzipData(buffer); err != nil {
			return nil, err
		}
	}

	cm := ChunkManifest{}
	if e := json.Unmarshal(buffer, &cm); e != nil {
		return nil, e
	}

	sort.Slice(cm.Chunks, func(i, j int) bool {
		return cm.Chunks[i].Offset < cm.Chunks[j].Offset
	})

	return &cm, nil
}
