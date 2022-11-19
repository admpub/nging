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
package language

import "net/http"

type Config struct {
	Project      string
	Default      string
	Fallback     string
	AllList      []string
	RulesPath    []string
	MessagesPath []string
	Reload       bool
	fsFunc       func(string) http.FileSystem
}

func (c *Config) SetFSFunc(fsFunc func(string) http.FileSystem) *Config {
	c.fsFunc = fsFunc
	return c
}

func (c *Config) FSFunc() func(string) http.FileSystem {
	return c.fsFunc
}
