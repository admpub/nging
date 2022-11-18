// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		if home, err := HomeDir(); err == nil {
			userGopath := filepath.Join(home, `go`)
			if fi, err := os.Stat(userGopath); err == nil && fi.IsDir() {
				gopath = userGopath
			}
		}
	}
	var paths []string
	if IsWindows {
		gopath = strings.Replace(gopath, "\\", "/", -1)
		paths = strings.Split(gopath, ";")
	} else {
		paths = strings.Split(gopath, ":")
	}
	return paths
}

func FixDirSeparator(fpath string) string {
	if IsWindows {
		fpath = strings.Replace(fpath, "\\", "/", -1)
	}
	return fpath
}

var ErrFolderNotFound = errors.New("unable to locate source folder path")

// GetSrcPath returns app. source code path.
// It only works when you have src. folder in GOPATH,
// it returns error not able to locate source folder path.
func GetSrcPath(importPath string) (appPath string, err error) {
	paths := GetGOPATHs()
	for _, p := range paths {
		if IsExist(p + "/src/" + importPath + "/") {
			appPath = p + "/src/" + importPath + "/"
			break
		}
	}

	if len(appPath) == 0 {
		return "", ErrFolderNotFound
	}

	appPath = filepath.Dir(appPath) + "/"
	if IsWindows {
		// Replace all '\' to '/'.
		appPath = strings.Replace(appPath, "\\", "/", -1)
	}

	return appPath, nil
}

// HomeDir returns path of '~'(in Linux) on Windows,
// it returns error when the variable does not exist.
func HomeDir() (string, error) {
	return os.UserHomeDir()
}
