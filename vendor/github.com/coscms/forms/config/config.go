/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

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

type Config struct {
	ID           string      `json:"id"`
	Theme        string      `json:"theme"`
	Template     string      `json:"template"`
	Method       string      `json:"method"`
	Action       string      `json:"action"`
	Attributes   [][]string  `json:"attributes"`
	WithButtons  bool        `json:"withButtons"`
	Buttons      []string    `json:"buttons"`
	BtnsTemplate string      `json:"btnsTemplate"`
	Elements     []*Element  `json:"elements"`
	Languages    []*Language `json:"languages"`
}

func (c *Config) Clone() *Config {
	r := *c
	return &r
}

func (c *Config) HasName(name string) bool {
	return c.hasName(name, c.Elements, c.Languages)
}

func (c *Config) hasName(name string, elements []*Element, languages []*Language) bool {
	for _, elem := range elements {
		if elem.Name == name {
			return elem.Type != `langset` && elem.Type != `fieldset`
		}
		if elem.Type == `langset` {
			if c.hasName(name, elem.Elements, elem.Languages) {
				return true
			}
			continue
		}
		if elem.Type == `fieldset` {
			if c.hasName(name, elem.Elements, languages) {
				return true
			}
			continue
		}
		if len(languages) == 0 {
			continue
		}
		for _, lang := range languages {
			if lang.HasName(name) || name == lang.Name(elem.Name) {
				return true
			}
		}
	}
	return false
}

func (c *Config) GetNames() []string {
	return getNames(c.Elements, c.Languages)
}

type Element struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Label      string      `json:"label"`
	Value      string      `json:"value"`
	HelpText   string      `json:"helpText"`
	Template   string      `json:"template"`
	Valid      string      `json:"valid"`
	Attributes [][]string  `json:"attributes"`
	Choices    []*Choice   `json:"choices"`
	Elements   []*Element  `json:"elements"`
	Format     string      `json:"format"`
	Languages  []*Language `json:"languages"`
}

func (e *Element) Clone() *Element {
	r := *e
	return &r
}

type Choice struct {
	Group   string   `json:"group"`
	Option  []string `json:"option"` //["value","text"]
	Checked bool     `json:"checked"`
}
