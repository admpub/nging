// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package caddyhttp

import (
	// plug in the server
	_ "github.com/admpub/caddy/caddyhttp/httpserver"

	// plug in the standard directives
	_ "github.com/admpub/caddy/caddyhttp/basicauth"
	_ "github.com/admpub/caddy/caddyhttp/bind"
	_ "github.com/admpub/caddy/caddyhttp/browse"
	_ "github.com/admpub/caddy/caddyhttp/errors"
	_ "github.com/admpub/caddy/caddyhttp/expvar"
	_ "github.com/admpub/caddy/caddyhttp/extensions"
	_ "github.com/admpub/caddy/caddyhttp/fastcgi"
	_ "github.com/admpub/caddy/caddyhttp/gzip"
	_ "github.com/admpub/caddy/caddyhttp/header"
	_ "github.com/admpub/caddy/caddyhttp/index"
	_ "github.com/admpub/caddy/caddyhttp/internalsrv"
	_ "github.com/admpub/caddy/caddyhttp/limits"
	_ "github.com/admpub/caddy/caddyhttp/log"
	_ "github.com/admpub/caddy/caddyhttp/markdown"
	_ "github.com/admpub/caddy/caddyhttp/mime"
	_ "github.com/admpub/caddy/caddyhttp/pprof"
	_ "github.com/admpub/caddy/caddyhttp/proxy"
	_ "github.com/admpub/caddy/caddyhttp/push"
	_ "github.com/admpub/caddy/caddyhttp/redirect"
	_ "github.com/admpub/caddy/caddyhttp/requestid"
	_ "github.com/admpub/caddy/caddyhttp/rewrite"
	_ "github.com/admpub/caddy/caddyhttp/root"
	_ "github.com/admpub/caddy/caddyhttp/status"
	_ "github.com/admpub/caddy/caddyhttp/templates"
	_ "github.com/admpub/caddy/caddyhttp/timeouts"
	_ "github.com/admpub/caddy/caddyhttp/websocket"
	_ "github.com/admpub/caddy/onevent"
)
