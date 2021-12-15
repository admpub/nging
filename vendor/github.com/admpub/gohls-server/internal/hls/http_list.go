package hls

import (
	"net/http"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/admpub/gohls-server/internal/fileindex"
)

type ListResponseVideo struct {
	Name string     `json:"name"`
	Path string     `json:"path"`
	Info *VideoInfo `json:"info"`
}

type ListResponseFolder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ListResponse struct {
	Error   error                  `json:"error"`
	Name    string                 `json:"name"`
	Path    string                 `json:"path"`
	Parents *[]*ListResponseFolder `json:"parents"`
	Folders []*ListResponseFolder  `json:"folders"`
	Videos  []*ListResponseVideo   `json:"videos"`
}

type ListHandler struct {
	idx     fileindex.Index
	name    string
	rootUri string
}

func NewListHandler(idx fileindex.Index, name string, rootUri string) *ListHandler {
	return &ListHandler{idx, name, rootUri}
}

func (s *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.idx.WaitForReady()
	videos := make([]*ListResponseVideo, 0)
	folders := make([]*ListResponseFolder, 0)

	parents := make([]*ListResponseFolder, 0)
	response := &ListResponse{nil, s.name, s.rootUri, &parents, folders, videos}

	if r.URL.Path != "" {
		entry, err := s.idx.Get(r.URL.Path)
		if err != nil || entry == nil {
			ServeJson(404, err, w)
			return
		}

		curr := entry
		for curr.ParentId() != "" {
			curr, err = s.idx.Get(curr.ParentId())
			if err != nil {
				ServeJson(500, err, w)
				return
			}
			parents = append(parents, &ListResponseFolder{curr.Name(), path.Join(s.rootUri, curr.Id())})
		}
		parents = append(parents, &ListResponseFolder{s.name, s.rootUri})

		response.Path = path.Join(s.rootUri, entry.Id())
		response.Name = entry.Name()
	}

	parents = append(parents, &ListResponseFolder{"Home", ""})

	files, err := s.idx.List(r.URL.Path)
	if err != nil {
		ServeJson(404, err, w)
		return
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") || strings.HasPrefix(f.Name(), "$") {
			continue
		}
		if FilenameLooksLikeVideo(f.Path()) {
			vinfo, err := GetVideoInformation(f.Path())
			if err != nil {
				log.Errorf("Could not read video information of %v: %v", f.Path(), err)
				continue
			}
			video := &ListResponseVideo{f.Name(), path.Join(s.rootUri, f.Id()), vinfo}
			videos = append(videos, video)
		}
		if f.IsDir() {
			folder := &ListResponseFolder{f.Name(), path.Join(s.rootUri, f.Id())}
			folders = append(folders, folder)
		}
	}
	response.Videos = videos
	response.Folders = folders
	ServeJson(200, response, w)
}
