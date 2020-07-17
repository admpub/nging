#!/bin/sh

#
# this script is only used when building caddy-rate-limit docker image
#

ed $GOPATH/src/github.com/caddyserver/caddy/caddy/caddymain/run.go << EOF
41i
    _ "github.com/xuqingfeng/caddy-rate-limit"
.
w
q
EOF