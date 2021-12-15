package hls

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/admpub/gohls-server/internal/fileindex"
)

// Encodes a string like Javascript's encodeURIComponent()
func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

type PlaylistHandler struct {
	idx          fileindex.Index
	rootUri      string
	segmentsPath string
}

func NewPlaylistHandler(idx fileindex.Index, rootUri string, segmentsPath string) *PlaylistHandler {
	return &PlaylistHandler{idx, rootUri, segmentsPath}
}

func (s *PlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Playlist request: %v", r.URL.Path)
	s.idx.WaitForReady()
	entry, err := s.idx.Get(r.URL.Path)
	if err != nil {
		ServeJson(404, err, w)
		return
	}
	file := entry.Path()

	vinfo, err := GetVideoInformation(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	duration := vinfo.Duration
	protocol := "http"
	if r.Header.Get("X-Forwarded-Proto") != "" {
		protocol = r.Header.Get("X-Forwarded-Proto")
	}
	baseurl := fmt.Sprintf("%v://%v", protocol, r.Host)

	id, err := urlEncoded(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/vnd.apple.mpegurl"}
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprint(w, "#EXT-X-TARGETDURATION:"+fmt.Sprintf("%.f", hlsSegmentLength)+"\n")
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segmentIndex := 0

	for leftover > 0 {
		if leftover > hlsSegmentLength {
			fmt.Fprintf(w, "#EXTINF: %f,\n", hlsSegmentLength)
		} else {
			fmt.Fprintf(w, "#EXTINF: %f,\n", leftover)
		}
		fmt.Fprintf(w, baseurl+s.segmentsPath+"%v/%v.ts\n", id, segmentIndex)
		segmentIndex++
		leftover = leftover - hlsSegmentLength
	}
	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}
