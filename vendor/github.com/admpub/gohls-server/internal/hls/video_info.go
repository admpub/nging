package hls

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var VideoSuffixes = []string{".mp4", ".rmvb", ".avi", ".mkv", ".flv", ".wmv", ".mov", ".mpg"}

func AddVideoSuffix(suffixes ...string) {
	for _, suffix := range suffixes {
		suffix = strings.TrimSpace(suffix)
		if len(suffix) == 0 {
			continue
		}
		if !strings.HasPrefix(suffix, `.`) {
			suffix = `.` + suffix
		}
		var exists bool
		for _, vs := range VideoSuffixes {
			if suffix == vs {
				exists = true
				break
			}
		}
		if !exists {
			VideoSuffixes = append(VideoSuffixes, suffix)
		}
	}
}

// TODO make mutex
var videoInfos = make(map[string]*VideoInfo)

type VideoInfo struct {
	Duration float64 `json:"duration"`
	//FileCreated      time.Time `json:"created"`
	FileLastModified time.Time `json:"lastModified"`
}

func FilenameLooksLikeVideo(name string) bool {
	name = strings.ToLower(name)
	for _, suffix := range VideoSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func GetRawFFMPEGInfo(path string) ([]byte, error) {
	log.Debugf("Executing ffprobe for %v", path)
	cmd := exec.Command(FFProbePath, "-v", "quiet", "-print_format", "json", "-show_format", path)
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing ffprobe for file '%v': %v", path, err)
	}
	return data, nil
}

func GetFFMPEGJson(path string) (map[string]interface{}, error) {
	data, cmderr := GetRawFFMPEGInfo(path)
	if cmderr != nil {
		return nil, cmderr
	}
	var info map[string]interface{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON from ffprobe output for file '%v':", path, err)
	}
	return info, nil
}

func GetVideoInformation(path string) (*VideoInfo, error) {
	if data, ok := videoInfos[path]; ok {
		if data == nil {
			return nil, fmt.Errorf("no video data available due to previous error for %v", path)
		}
		return data, nil
	}
	info, jsonerr := GetFFMPEGJson(path)
	if jsonerr != nil {
		videoInfos[path] = nil
		return nil, jsonerr
	}
	log.Debugf("ffprobe for %v returned", path, info)
	if _, ok := info["format"]; !ok {
		return nil, fmt.Errorf("ffprobe data for '%v' does not contain format info", path)
	}
	format := info["format"].(map[string]interface{})
	if _, ok := format["duration"]; !ok {
		return nil, fmt.Errorf("ffprobe format data for '%v' does not contain duration", path)
	}
	duration, perr := strconv.ParseFloat(format["duration"].(string), 64)
	if perr != nil {
		return nil, fmt.Errorf("Could not parse duration (%v) of '%v' ", format["duration"].(string), path, perr)
	}
	finfo, staterr := os.Stat(path)
	if staterr != nil {
		return nil, fmt.Errorf("Could not stat file '%v': %v", path, staterr)
	}
	var vi = &VideoInfo{duration, finfo.ModTime()}
	videoInfos[path] = vi
	return vi, nil
}
