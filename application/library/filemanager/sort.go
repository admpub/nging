/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package filemanager

import "os"

type SortByFileType []os.FileInfo

func (s SortByFileType) Len() int { return len(s) }
func (s SortByFileType) Less(i, j int) bool {
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
func (s SortByFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByModTime []os.FileInfo

func (s SortByModTime) Len() int { return len(s) }
func (s SortByModTime) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}
func (s SortByModTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByModTimeDesc []os.FileInfo

func (s SortByModTimeDesc) Len() int { return len(s) }
func (s SortByModTimeDesc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() > s[j].ModTime().UnixNano()
}
func (s SortByModTimeDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByNameDesc []os.FileInfo

func (s SortByNameDesc) Len() int { return len(s) }
func (s SortByNameDesc) Less(i, j int) bool {
	return s[i].Name() > s[j].Name()
}
func (s SortByNameDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
