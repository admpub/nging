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

package seaweedfs

import (
	"net/url"
	"time"

	"github.com/admpub/goseaweedfs"
	"github.com/admpub/goseaweedfs/libs"
)

var DefaultConfig = &Config{}

type FilerURL struct {
	Public  string //Readonly URL
	Private string //Manage URL
}

type Config struct {
	Scheme    string
	Master    string
	Filers    []*FilerURL
	ChunkSize int64
	Timeout   time.Duration
	// TTL Time to live.
	// 3m: 3 minutes
	// 4h: 4 hours
	// 5d: 5 days
	// 6w: 6 weeks
	// 7M: 7 months
	// 8y: 8 years
	TTL string
}

func (c *Config) New() *goseaweedfs.Seaweed {
	if len(c.Scheme) == 0 {
		c.Scheme = "http"
	}
	if c.ChunkSize <= 0 {
		c.ChunkSize = 2 * 1024 * 1024
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Minute
	}
	if len(c.Master) == 0 {
		c.Master = `localhost:9333`
	}
	if c.Filers == nil || len(c.Filers) == 0 {
		c.Filers = []*FilerURL{
			{
				Public:  `http://localhost:8989`,
				Private: `http://localhost:8888`,
			},
		}
	}
	filers := make([]string, len(c.Filers))
	for index, filerURL := range c.Filers {
		filers[index] = filerURL.Private
	}
	return goseaweedfs.NewSeaweed(c.Scheme, c.Master, filers, c.ChunkSize, c.Timeout)
}

func (c *Config) MakeURL(path string, args url.Values) string {
	return libs.MakeURL(c.Scheme, c.Master, path, args)
}
