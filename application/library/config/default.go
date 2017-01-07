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
package config

import (
	"database/sql"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/webx-top/com"
)

var (
	Installed             sql.NullBool
	DefaultConfig         = &Config{}
	DefaultCLIConfig      = &CLIConfig{cmds: map[string]*exec.Cmd{}}
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
)

func IsInstalled() bool {
	if !Installed.Valid {
		lockFile := filepath.Join(com.SelfDir(), `installed.lock`)
		if info, err := os.Stat(lockFile); err == nil && info.IsDir() == false {
			Installed.Valid = true
			Installed.Bool = true
		}
	}
	return Installed.Bool
}
