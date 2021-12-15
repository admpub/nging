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

package frp

import (
	"log"
	"os/exec"
	"sync"
)

func StartFrp(nodeName string, configFile string, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer func() {
		wg.Done()
		log.Printf("%s is exit!", nodeName)
	}()
	cmd := exec.Command("frps", "-c", configFile)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return err
}

// Run {nodeName:configFile}
func Run(mConfig map[string]string) {
	wg := &sync.WaitGroup{}
	log.Println("frpManger Server was start !")
	for k, v := range mConfig {
		log.Println(k + " was start!")
		go func(k, v string) {
			err := StartFrp(k, v, wg)
			if err != nil {
				log.Println(k+" has an error:", err)
			}
		}(k, v)
	}
	wg.Wait()
}
