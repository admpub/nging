package hls

import (
	"net/http"
	"path"

	"github.com/admpub/gohls-server/internal/fileindex"
)

type InfoHandler struct {
	idx     fileindex.Index
	title   string
	rootUri string
}

func NewInfoHandler(idx fileindex.Index, title string, rootUri string) *InfoHandler {
	return &InfoHandler{idx, title, rootUri}
}

func (s *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.idx.WaitForReady()

	folderName := ""
	folderPath := ""

	entry, err := s.idx.Get(r.URL.Path)
	if err != nil {
		ServeJson(404, err, w)
		return
	}
	if !FilenameLooksLikeVideo(entry.Path()) {
		ServeJson(404, "Not found", w)
		return
	}

	vinfo, err := GetVideoInformation(entry.Path())
	if err != nil {
		ServeJson(500, err, w)
		return
	}

	if entry.ParentId() != "" {
		folder, err := s.idx.Get(entry.ParentId())
		if err != nil {
			ServeJson(404, err, w)
			return
		}
		folderName = folder.Name()
		folderPath = ""
	}

	videos := make([]*ListResponseVideo, 0)
	folders := make([]*ListResponseFolder, 0)
	parents := make([]*ListResponseFolder, 0)
	response := &ListResponse{nil, folderName, folderPath, &parents, folders, videos}

	curr := entry
	for curr.ParentId() != "" {
		curr, err = s.idx.Get(curr.ParentId())
		if err != nil {
			ServeJson(500, err, w)
			return
		}
		parents = append(parents, &ListResponseFolder{curr.Name(), path.Join(s.rootUri, curr.Id())})
	}

	parents = append(parents, &ListResponseFolder{s.title, s.rootUri})

	parents = append(parents, &ListResponseFolder{"Home", ""})

	videos = append(videos, &ListResponseVideo{entry.Name(), path.Join(s.rootUri, entry.Id()), vinfo})
	response.Videos = videos
	response.Parents = &parents

	ServeJson(200, response, w)

}
