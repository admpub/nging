package goseaweedfs

import "math/rand"

// VolumeLocation location of volume responsed from master API. According to https://github.com/chrislusf/seaweedfs/wiki/Master-Server-API
type VolumeLocation struct {
	URL       string `json:"url,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
}

// VolumeLocations returned VolumeLocations (volumes)
type VolumeLocations []*VolumeLocation

// Head get first location in list
func (c VolumeLocations) Head() *VolumeLocation {
	if len(c) == 0 {
		return nil
	}

	return c[0]
}

// RandomPickForRead random pick a location for further read request
func (c VolumeLocations) RandomPickForRead() *VolumeLocation {
	if len(c) == 0 {
		return nil
	}

	return c[rand.Intn(len(c))]
}

// LookupResult the result of looking up volume. According to https://github.com/chrislusf/seaweedfs/wiki/Master-Server-API
type LookupResult struct {
	VolumeLocations VolumeLocations `json:"locations,omitempty"`
	Error           string          `json:"error,omitempty"`
}
