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
package forms

import (
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
)

type ElementSetter interface {
	Elements(...config.FormElement)
}

func NewConfig() *config.Config {
	return &config.Config{
		ID:       `Forms`,
		Theme:    `bootstrap3`,
		Template: ``,
		Method:   `POST`,
		Attributes: [][]string{
			[]string{"class", "form-horizontal"},
			[]string{"role", "form"},
		},
		WithButtons: true,
		Buttons:     make([]string, 0),
		Elements:    make([]*config.Element, 0),
	}
}

// GenChoices generate choices
// type Data struct{
// 	ID string
// 	Name string
// }
// data:=[]*Data{
// 	&Data{ID:"a",Name:"One"},
// 	&Data{ID:"b",Name:"Two"},
// }
// GenChoices(len(data), func(index int) (string, string, bool){
// 	return data[index].ID,data[index].Name,false
// })
// or
// GenChoices(map[string]int{
// 	"":len(data),
// }, func(group string,index int) (string, string, bool){
// 	return data[index].ID,data[index].Name,false
// })
func GenChoices(lenType interface{}, fnType interface{}) interface{} {
	switch fn := fnType.(type) {
	case func(int) (string, string, bool):
		length, ok := lenType.(int)
		if !ok {
			return []fields.InputChoice{}
		}
		result := make([]fields.InputChoice, length)
		for key, r := range result {
			r.ID, r.Val, r.Checked = fn(key)
			result[key] = r
		}
		return result
	case func(string, int) (string, string, bool):
		result := make(map[string][]fields.InputChoice)
		values, ok := lenType.(map[string]int)
		if !ok {
			return result
		}
		for group, length := range values {
			if _, ok := result[group]; !ok {
				result[group] = make([]fields.InputChoice, length)
			}
			for key, r := range result[group] {
				r.ID, r.Val, r.Checked = fn(group, key)
				result[group][key] = r
			}
		}
		return result
	}
	return nil
}
