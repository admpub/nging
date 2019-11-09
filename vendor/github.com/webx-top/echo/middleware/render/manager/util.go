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

package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func dump(i interface{}) {
	c, _ := json.MarshalIndent(i, "", " ")
	fmt.Println(string(c))
}

func FixDirSeparator(dir string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(dir, "\\", "/", -1)
	}
	return dir
}

func dirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}

	return true
}

func fileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
