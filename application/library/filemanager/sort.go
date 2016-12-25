/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
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
