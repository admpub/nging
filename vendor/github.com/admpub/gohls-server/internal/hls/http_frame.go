package hls

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/admpub/gohls-server/internal/fileindex"
)

type FrameHandler struct {
	idx        fileindex.Index
	rootUri    string
	cmdHandler *HttpCommandHandler
}

func NewFrameHandler(idx fileindex.Index, rootUri string) *FrameHandler {
	return &FrameHandler{idx, rootUri, NewHttpCommandHandler(2, "frames")}
}

func (s *FrameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	time := 30
	if tint, err := strconv.Atoi(t); err == nil {
		time = tint
	}
	s.idx.WaitForReady()
	entry, err := s.idx.Get(r.URL.Path)
	if err != nil {
		ServeJson(404, err, w)
		return
	}
	path := entry.Path()
	args := []string{
		"-timelimit", "15",
		"-loglevel", "error",
		"-ss", fmt.Sprintf("%v.0", time),
		"-i", path,
		"-vf", "scale=320:-1",
		"-frames:v", "1",
		"-f", "image2",
		"-",
	}
	if err := s.cmdHandler.ServeCommand(FFMPEGPath, args, calculateCommandHash(FFMPEGPath, args), w); err != nil {
		log.Errorf("Problem serving screenshot: %v", err)
	}
}
