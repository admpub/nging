package goseaweedfs

// UploadResult contains upload result after put file to SeaweedFS
// Raw response: {"name":"go1.8.3.linux-amd64.tar.gz","size":82565628,"error":""}
type UploadResult struct {
	Name  string `json:"name,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Error string `json:"error,omitempty"`
}

// AssignResult contains assign result.
// Raw response: {"fid":"1,0a1653fd0f","url":"localhost:8899","publicUrl":"localhost:8899","count":1,"error":""}
type AssignResult struct {
	FileID    string `json:"fid,omitempty"`
	URL       string `json:"url,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
	Count     uint64 `json:"count,omitempty"`
	Error     string `json:"error,omitempty"`
}

// SubmitResult result of submit operation.
type SubmitResult struct {
	FileName string `json:"fileName,omitempty"`
	FileURL  string `json:"fileUrl,omitempty"`
	FileID   string `json:"fid,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ClusterStatus result of getting status of cluster
type ClusterStatus struct {
	IsLeader bool
	Leader   string
	Peers    []string
}

// SystemStatus result of getting status of system
type SystemStatus struct {
	Topology Topology
	Version  string
	Error    string
}

// Topology result of topology stats request
type Topology struct {
	DataCenters []*DataCenter
	Free        int
	Max         int
	Layouts     []*Layout
}

// DataCenter stats of a datacenter
type DataCenter struct {
	Free  int
	Max   int
	Racks []*Rack
}

// Rack stats of racks
type Rack struct {
	DataNodes []*DataNode
	Free      int
	Max       int
}

// DataNode stats of data node
type DataNode struct {
	Free      int
	Max       int
	PublicURL string `json:"PublicUrl"`
	URL       string `json:"Url"`
	Volumes   int
}

// Layout of replication/collection stats. According to https://github.com/chrislusf/seaweedfs/wiki/Master-Server-API
type Layout struct {
	Replication string
	Writables   []uint64
}
