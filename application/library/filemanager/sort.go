package filemanager

import "os"

type byFileType []os.FileInfo

func (s byFileType) Len() int { return len(s) }
func (s byFileType) Less(i, j int) bool {
	if s[i].IsDir() {
		if !s[j].IsDir() {
			return true
		}
	} else if s[j].IsDir() {
		if !s[i].IsDir() {
			return false
		}
	}
	return s[i].Name() < s[j].Name()
}
func (s byFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byModTime []os.FileInfo

func (s byModTime) Len() int { return len(s) }
func (s byModTime) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}
func (s byModTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byModTimeDesc []os.FileInfo

func (s byModTimeDesc) Len() int { return len(s) }
func (s byModTimeDesc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() > s[j].ModTime().UnixNano()
}
func (s byModTimeDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byNameDesc []os.FileInfo

func (s byNameDesc) Len() int { return len(s) }
func (s byNameDesc) Less(i, j int) bool {
	return s[i].Name() > s[j].Name()
}
func (s byNameDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
