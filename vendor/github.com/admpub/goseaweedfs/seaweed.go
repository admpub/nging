package goseaweedfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/admpub/goseaweedfs/libs"
	"github.com/admpub/goseaweedfs/model"

	cache "github.com/patrickmn/go-cache"
)

var (
	cacheDuration = 10 * time.Minute

	// ErrFileNotFound return file not found error
	ErrFileNotFound = fmt.Errorf("File not found")
)

const (
	// ParamCollection http param to specify collection which files belong. According to SeaweedFS API.
	ParamCollection = "collection"

	// ParamTTL http param to specify time to live. According to SeaweedFS API.
	ParamTTL = "ttl"

	// ParamCount http param to specify how many file ids to reserve. According to SeaweedFS API.
	ParamCount = "count"

	// ParamAssignReplication http param to assign files with a specific replication type.
	ParamAssignReplication = "replication"

	// ParamAssignCount http param to specify how many file ids to reserve.
	ParamAssignCount = "count"

	// ParamAssignDataCenter http param to assign a specific data center
	ParamAssignDataCenter = "dataCenter"

	// ParamLookupVolumeID http param to specify volume ID for looking up.
	ParamLookupVolumeID = "volumeId"

	// ParamLookupPretty http param to make json response prettified or not. Default should not be set.
	ParamLookupPretty = "pretty"

	// ParamLookupCollection http param to specify known collection, this would make file look up/search faster.
	ParamLookupCollection = "collection"

	// ParamVacuumGarbageThreshold if your system has many deletions, the deleted file's disk space will not be synchronously re-claimed.
	// There is a background job to check volume disk usage. If empty space is more than the threshold,
	// default to 0.3, the vacuum job will make the volume readonly, create a new volume with only existing files,
	// and switch on the new volume. If you are impatient or doing some testing, vacuum the unused spaces this way.
	ParamVacuumGarbageThreshold = "GarbageThreshold"

	// ParamGrowReplication http param to specify a specific replication.
	ParamGrowReplication = "replication"

	// ParamGrowCount http param to specify number of empty volume to grow.
	ParamGrowCount = "count"

	// ParamGrowDataCenter http param to specify datacenter of growing volume.
	ParamGrowDataCenter = "dataCenter"

	// ParamGrowCollection http param to specify collection of files for growing.
	ParamGrowCollection = "collection"

	// ParamGrowTTL specify time to live for growing api. Refers to: https://github.com/chrislusf/seaweedfs/wiki/Store-file-with-a-Time-To-Live
	// 3m: 3 minutes
	// 4h: 4 hours
	// 5d: 5 days
	// 6w: 6 weeks
	// 7M: 7 months
	// 8y: 8 years
	ParamGrowTTL = "ttl"

	// admin operations
	// ParamAssignVolumeReplication = "replication"
	// ParamAssignVolume            = "volume"
	// ParamDeleteVolume            = "volume"
	// ParamMountVolume             = "volume"
	// ParamUnmountVolume           = "volume"
)

// Seaweed client containing almost features/operations to interact with SeaweedFS
type Seaweed struct {
	Master    string
	Filers    []*model.Filer
	Scheme    string
	ChunkSize int64
	Client    *libs.HTTPClient
	cache     *cache.Cache
}

// NewSeaweed create new seaweed with default
func NewSeaweed(scheme string, master string, filers []string, chunkSize int64, timeout time.Duration) *Seaweed {
	res := &Seaweed{
		Master:    master,
		Scheme:    scheme,
		Client:    libs.NewHTTPClient(timeout),
		cache:     cache.New(cacheDuration, cacheDuration*2),
		ChunkSize: chunkSize,
	}
	if filers != nil {
		res.Filers = make([]*model.Filer, len(filers))
		for i := range filers {
			res.Filers[i] = model.NewFiler(filers[i], res.Client)
		}
	}

	return res
}

// Grow pre-Allocate Volumes
func (c *Seaweed) Grow(count int, collection, replication, dataCenter string) error {
	args := url.Values{}
	if count > 0 {
		args.Set(ParamGrowCount, strconv.Itoa(count))
	}
	if collection != "" {
		args.Set(ParamGrowCollection, collection)
	}
	if replication != "" {
		args.Set(ParamGrowReplication, replication)
	}
	if dataCenter != "" {
		args.Set(ParamGrowDataCenter, dataCenter)
	}

	return c.GrowArgs(args)
}

// GrowArgs pre-Allocate volumes with args
func (c *Seaweed) GrowArgs(args url.Values) (err error) {
	_, _, err = c.Client.Get(c.Scheme, c.Master, "/vol/grow", args)
	return
}

// Lookup volume ID
func (c *Seaweed) Lookup(volID string, args url.Values) (result *model.LookupResult, err error) {
	if item, exist := c.cache.Get(volID); !exist || item == nil {
		if result, err = c.doLookup(volID, args); err == nil {
			c.cache.Set(volID, result, cacheDuration)
		}
	} else {
		switch item.(type) {
		case *model.LookupResult:
			result, err = item.(*model.LookupResult), nil
			return
		}

		if result, err = c.doLookup(volID, args); err == nil {
			c.cache.Set(volID, result, cacheDuration)
		}
	}

	return
}

// LookupNoCache lookup by volume id without get from caching first, but set cache in the end of process.
func (c *Seaweed) LookupNoCache(volID string, args url.Values) (result *model.LookupResult, err error) {
	if result, err = c.doLookup(volID, args); err == nil {
		c.cache.Set(volID, result, cacheDuration)
	}
	return
}

func (c *Seaweed) doLookup(volID string, args url.Values) (result *model.LookupResult, err error) {
	if args == nil {
		args = make(url.Values)
	}
	args.Set(ParamLookupVolumeID, volID)

	jsonBlob, _, err := c.Client.PostForm(libs.MakeURL(c.Scheme, c.Master, "/dir/lookup", nil), args)
	if err != nil {
		return nil, err
	}

	result = &model.LookupResult{}
	if err = json.Unmarshal(jsonBlob, result); err != nil {
		return
	}

	if result.Error != "" {
		err = errors.New(result.Error)
		return
	}

	return
}

// LookupServerByFileID lookup server by fileID
func (c *Seaweed) LookupServerByFileID(fileID string, args url.Values, readonly bool) (server string, err error) {
	var parts []string
	if strings.Contains(fileID, ",") {
		parts = strings.Split(fileID, ",")
	} else {
		parts = strings.Split(fileID, "/")
	}

	if len(parts) != 2 { // wrong file id format
		return "", errors.New("Invalid fileID " + fileID)
	}

	lookup, lookupError := c.Lookup(parts[0], args)
	if lookupError != nil {
		err = lookupError
		return
	} else if len(lookup.VolumeLocations) == 0 {
		err = ErrFileNotFound
		return
	}

	if readonly {
		server = lookup.VolumeLocations.RandomPickForRead().URL
	} else {
		server = lookup.VolumeLocations.Head().URL
	}

	return
}

// LookupFileID lookup file by id
func (c *Seaweed) LookupFileID(fileID string, args url.Values, readonly bool) (fullURL string, err error) {
	u, err := c.LookupServerByFileID(fileID, args, readonly)
	if err != nil {
		return
	}

	fullURL = libs.MakeURL(c.Scheme, u, fileID, nil)
	return
}

// LookupVolumeIDs find volume locations by cache and actual lookup
func (c *Seaweed) LookupVolumeIDs(volIDs []string) (result map[string]*model.LookupResult, err error) {
	result = make(map[string]*model.LookupResult)

	//
	unknownVolIDs := make([]string, len(volIDs))
	n := 0

	//check vid cache first
	for _, vid := range volIDs {
		if item, exist := c.cache.Get(vid); exist && item != nil {
			result[vid] = item.(*model.LookupResult)
		} else {
			unknownVolIDs[n] = vid
			n++
		}
	}

	if n == 0 {
		return
	}

	//only query unknown_vids
	args := url.Values{}
	for i := 0; i < n; i++ {
		args.Add("volumeId", unknownVolIDs[i])
	}

	jsonBlob, _, err := c.Client.PostForm(libs.MakeURL(c.Scheme, c.Master, "/vol/lookup", nil), args)
	if err != nil {
		return
	}

	ret := make(map[string]*model.LookupResult)
	if err = json.Unmarshal(jsonBlob, &ret); err != nil {
		return
	}

	for k, v := range ret {
		result[k] = v
		c.cache.Set(k, v, cacheDuration)
	}

	err = nil
	return
}

// GC force Garbage Collection
func (c *Seaweed) GC(threshold float64) (err error) {
	args := url.Values{
		"garbageThreshold": []string{strconv.FormatFloat(threshold, 'f', -1, 64)},
	}

	if _, _, err = c.Client.Get(c.Scheme, c.Master, "/vol/vacuum", args); err != nil {
		// TODO: handle response later
		return
	}

	return
}

// Status check System Status
func (c *Seaweed) Status() (result *model.SystemStatus, err error) {
	data, _, err := c.Client.Get(c.Scheme, c.Master, "/dir/status", nil)
	if err != nil {
		return
	}

	result = &model.SystemStatus{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

// ClusterStatus get cluster status
func (c *Seaweed) ClusterStatus() (result *model.ClusterStatus, err error) {
	data, _, err := c.Client.Get(c.Scheme, c.Master, "/cluster/status", nil)
	if err != nil {
		return
	}

	result = &model.ClusterStatus{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

// Assign do assign api
func (c *Seaweed) Assign(args url.Values) (result *model.AssignResult, err error) {
	if args == nil {
		args = make(url.Values)
	}

	jsonBlob, _, err := c.Client.PostForm(libs.MakeURL(c.Scheme, c.Master, "/dir/assign", nil), args)
	if err != nil {
		return nil, err
	}

	result = &model.AssignResult{}
	if err = json.Unmarshal(jsonBlob, result); err != nil {
		err = fmt.Errorf("/dir/assign result JSON unmarshal error:%v, json:%s", err, string(jsonBlob))
		return
	} else if result.Count <= 0 {
		err = errors.New(result.Error)
		return
	}

	return
}

// Submit file directly to master
func (c *Seaweed) Submit(filePath string, collection, ttl string) (result *model.SubmitResult, err error) {
	fp, err := model.NewFilePart(filePath)
	if err != nil {
		return
	}
	fp.Collection = collection
	fp.TTL = ttl

	return c.SubmitFilePart(fp, url.Values{})
}

// SubmitFilePart directly to master
func (c *Seaweed) SubmitFilePart(f *model.FilePart, args url.Values) (result *model.SubmitResult, err error) {
	data, _, err := c.Client.Upload(libs.MakeURL(c.Scheme, c.Master, "/submit", args), f.FileName, f.Reader, f.IsGzipped, f.MimeType)
	if err != nil {
		return
	}

	result = &model.SubmitResult{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

// Upload file by reader
func (c *Seaweed) Upload(fileReader io.Reader, fileName string, size int64, collection, ttl string) (fp *model.FilePart, fileID string, err error) {
	fp = model.NewFilePartFromReader(fileReader, fileName, size)
	fp.Collection, fp.TTL = collection, ttl

	_, fileID, err = c.UploadFilePart(fp)
	return
}

// UploadFile upload file with full file dir/path
func (c *Seaweed) UploadFile(filePath string, collection, ttl string) (cm *model.ChunkManifest, fp *model.FilePart, fileID string, err error) {
	fp, err = model.NewFilePart(filePath)
	if err != nil {
		return
	}
	fp.Collection, fp.TTL = collection, ttl

	cm, fileID, err = c.UploadFilePart(fp)
	return
}

// UploadFilePart upload a file part
func (c *Seaweed) UploadFilePart(f *model.FilePart) (cm *model.ChunkManifest, fileID string, err error) {
	if f.FileID == "" {
		args := make(url.Values)
		if f.Collection != "" {
			args.Set(ParamCollection, f.Collection)
		}
		if f.TTL != "" {
			args.Set(ParamTTL, f.TTL)
		}
		args.Set(ParamAssignCount, "1")

		res, err := c.Assign(args)
		if err != nil {
			return nil, "", err
		}
		f.Server, f.FileID = res.URL, res.FileID
	}

	if f.Server == "" {
		if f.Server, err = c.LookupServerByFileID(f.FileID, url.Values{ParamCollection: []string{f.Collection}}, false); err != nil {
			return
		}
	}

	if closer, ok := f.Reader.(io.Closer); ok { // closing after read content
		defer closer.Close()
	}

	baseName := path.Base(f.FileName)
	if c.ChunkSize > 0 && f.FileSize > c.ChunkSize {
		chunks := f.FileSize/c.ChunkSize + 1

		cm = &model.ChunkManifest{
			Name:   baseName,
			Size:   f.FileSize,
			Mime:   f.MimeType,
			Chunks: make([]*model.ChunkInfo, chunks),
		}
		args := url.Values{ParamCollection: []string{f.Collection}}

		for i := int64(0); i < chunks; i++ {
			_, id, count, e := c.uploadChunk(f, baseName+"_"+strconv.FormatInt(i+1, 10))
			if e != nil { // delete all uploaded chunks
				c.DeleteChunks(cm, args)
				return nil, "", e
			}

			cm.Chunks[i] = &model.ChunkInfo{
				Offset: i * c.ChunkSize,
				Size:   int64(count),
				Fid:    id,
			}
		}

		if err = c.uploadManifest(f, cm); err != nil { // delete all uploaded chunks
			c.DeleteChunks(cm, args)
		}
	} else {
		args := url.Values{}
		if f.ModTime != 0 {
			args.Set("ts", strconv.FormatInt(f.ModTime, 10))
		}

		_, _, err = c.Client.Upload(libs.MakeURL(c.Scheme, f.Server, f.FileID, args), baseName, f.Reader, f.IsGzipped, f.MimeType)
	}

	if err == nil {
		fileID = f.FileID
	}

	return
}

// BatchUploadFiles batch upload files
func (c *Seaweed) BatchUploadFiles(files []string, collection, ttl string) ([]*model.SubmitResult, error) {
	fps, e := model.NewFileParts(files)
	if e != nil {
		return nil, e
	}

	return c.BatchUploadFileParts(fps, collection, ttl)
}

// BatchUploadFileParts upload multiple file parts at once
func (c *Seaweed) BatchUploadFileParts(files []*model.FilePart, collection string, ttl string) ([]*model.SubmitResult, error) {
	results := make([]*model.SubmitResult, len(files))
	for index, file := range files {
		results[index] = &model.SubmitResult{
			FileName: file.FileName,
		}
	}

	args := make(url.Values)
	if collection != "" {
		args.Set(ParamCollection, collection)
	}
	if ttl != "" {
		args.Set(ParamTTL, ttl)
	}
	args.Set(ParamAssignCount, strconv.Itoa(len(files)))

	ret, err := c.Assign(args)
	if err != nil {
		for index := range files {
			results[index].Error = err.Error()
		}
		return results, err
	}

	wg := sync.WaitGroup{}
	for index, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int, file *model.FilePart) {
			file.FileID = ret.FileID
			if index > 0 {
				file.FileID = file.FileID + "_" + strconv.Itoa(index)
			}
			file.Server = ret.URL
			file.Collection = collection

			if _, _, err := c.UploadFilePart(file); err != nil {
				results[index].Error = err.Error()
			}

			results[index].Size = file.FileSize
			results[index].FileID = file.FileID
			results[index].FileURL = ret.PublicURL + "/" + file.FileID

			wg.Done()
		}(&wg, index, file)
	}
	wg.Wait()

	return results, nil
}

// Replace with file reader
func (c *Seaweed) Replace(fileID string, fileReader io.Reader, fileName string, size int64, collection, ttl string, deleteFirst bool) (err error) {
	fp := model.NewFilePartFromReader(fileReader, fileName, size)
	fp.Collection, fp.TTL = collection, ttl
	fp.FileID = fileID

	_, err = c.ReplaceFilePart(fp, deleteFirst)
	return
}

// ReplaceFile replace file by fileID with local filePath
func (c *Seaweed) ReplaceFile(fileID, filePath string, deleteFirst bool) error {
	fp, e := model.NewFilePart(filePath)
	if e != nil {
		return e
	}
	fp.FileID = fileID

	_, e = c.ReplaceFilePart(fp, deleteFirst)
	return e
}

// ReplaceFilePart replace file part
func (c *Seaweed) ReplaceFilePart(f *model.FilePart, deleteFirst bool) (fileID string, err error) {
	if deleteFirst && f.FileID != "" {
		c.DeleteFile(f.FileID, url.Values{ParamCollection: []string{f.Collection}})
	}

	_, fileID, err = c.UploadFilePart(f)
	return
}

func (c *Seaweed) uploadChunk(f *model.FilePart, filename string) (assignResult *model.AssignResult, fileID string, size int64, err error) {
	// Assign first to get file id and url for uploading
	assignResult, err = c.Assign(url.Values{
		ParamCollection:  []string{f.Collection},
		ParamTTL:         []string{f.TTL},
		ParamAssignCount: []string{"1"},
	})
	if err != nil {
		return
	}

	fileID = assignResult.FileID

	// now do upload
	dat, _, err := c.Client.Upload(
		libs.MakeURL(c.Scheme, assignResult.URL, assignResult.FileID, nil),
		filename, io.LimitReader(f.Reader, c.ChunkSize),
		false, "application/octet-stream")
	if err != nil {
		return
	}

	// parsing response data
	uploadResult := model.UploadResult{}
	if err = json.Unmarshal(dat, &uploadResult); err != nil {
		return
	}
	size = uploadResult.Size

	return
}

func (c *Seaweed) uploadManifest(f *model.FilePart, manifest *model.ChunkManifest) error {
	buf, e := manifest.Marshal()
	if e != nil {
		return e
	}
	bufReader := bytes.NewReader(buf)

	args := url.Values{}
	if f.ModTime != 0 {
		args.Set("ts", strconv.FormatInt(f.ModTime, 10))
	}
	args.Set("cm", "true")

	_, _, e = c.Client.Upload(libs.MakeURL(c.Scheme, f.Server, f.FileID, args), manifest.Name, bufReader, false, "application/json")
	return e
}

func (c *Seaweed) DownloadFile(fileID string, args url.Values) (string, []byte, error) {
	if args == nil {
		args = make(url.Values)
	}
	fileURL, err := c.LookupFileID(fileID, args, true)
	if err != nil {
		return "", nil, err
	}
	fileName, rc, err := c.Client.DownloadFromURL(fileURL)
	if err != nil {
		return "", nil, err
	}
	fileData, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", nil, err
	}
	defer rc.Close()
	return fileName, fileData, nil

}

// DeleteChunks concurrently delete chunks
func (c *Seaweed) DeleteChunks(cm *model.ChunkManifest, args url.Values) (err error) {
	if cm == nil || len(cm.Chunks) == 0 {
		return nil
	}

	result := make(chan bool, len(cm.Chunks))
	for _, ci := range cm.Chunks {
		go func(fileID string) {
			result <- c.DeleteFile(fileID, args) == nil
		}(ci.Fid)
	}

	isOk := true
	for i := 0; i < len(cm.Chunks); i++ {
		if r := <-result; !r {
			isOk = false
		}
	}

	if !isOk {
		err = errors.New("Not all chunks deleted.")
		return
	}

	return nil
}

// DeleteFile delete file by fileID
func (c *Seaweed) DeleteFile(fileID string, args url.Values) (err error) {
	fileURL, err := c.LookupFileID(fileID, args, false)
	if err != nil {
		return fmt.Errorf("Failed to lookup %s:%v", fileID, err)
	}

	if _, err = c.Client.Delete(fileURL); err != nil {
		err = fmt.Errorf("Failed to delete %s:%v", fileURL, err)
		return
	}

	return nil
}
