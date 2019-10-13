package hls

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/admpub/gohls-server/internal/fileindex"
)

var streamRegexp = regexp.MustCompile(`^(.*)/([0-9]+)\.ts$`)

type StreamHandler struct {
	idx     fileindex.Index
	rootUri string
	encoder *Encoder
}

func NewStreamHandler(idx fileindex.Index, rootUri string) *StreamHandler {
	return &StreamHandler{idx, rootUri, NewEncoder("segments", 2)}
}

func (s *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Stream request: %v", r.URL.Path)
	s.idx.WaitForReady()
	matches := streamRegexp.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.Error(w, "Wrong path format", 400)
		return
	}

	entry, err := s.idx.Get(matches[1])
	if err != nil {
		ServeJson(404, err, w)
		return
	}

	res := int64(720)
	segment, _ := strconv.ParseInt(matches[2], 0, 64)
	file := entry.Path()
	er := NewEncodingRequest(file, segment, res)
	s.encoder.Encode(*er)

	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
	select {
	case data := <-er.data:
		w.Write(*data)
	case err := <-er.err:
		log.Errorf("Error encoding %v", err)
	case <-time.After(60 * time.Second):
		log.Errorf("Timeout encoding")
	}
}
