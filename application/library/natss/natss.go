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

package natss

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	natsd "github.com/nats-io/nats-server/v2/server"
	stand "github.com/nats-io/nats-streaming-server/server"
	stores "github.com/nats-io/nats-streaming-server/stores"
	"github.com/webx-top/echo"
)

func Start() {
	// This configure the NATS Server using natsd package
	nOpts := &natsd.Options{
		HTTPPort: 8222,
		Port:     4223,
	}
	/*/ Create the NATS Server
	  ns := natsd.New(nOpts)

	  // Start it as a go routine
	  go ns.Start()

	  // Wait for it to be able to accept connections
	  if !ns.ReadyForConnections(10 * time.Second) {
	    panic("not able to start")
	  }
	  // */

	// Get NATS Streaming Server default options
	sOpts := stand.GetDefaultOptions()
	sOpts.StoreType = stores.TypeFile
	sOpts.FilestoreDir = filepath.Join(echo.Wd(), `data`, `nats`)
	// Force the streaming server to setup its own signal handler
	sOpts.HandleSignals = true
	// override the NoSigs for NATS since Streaming has its own signal handler
	nOpts.NoSigs = true
	// Without this option set to true, the logger is not configured.
	sOpts.EnableLogging = true
	// This will invoke RunServerWithOpts but on Windows, may run it as a service.
	if _, err := stand.Run(sOpts, nOpts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runtime.Goexit()
}
