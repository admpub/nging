#!/bin/sh
sudo launchctl limit maxfiles 65535
ulimit -n 65535
go get github.com/webx-top/tower
go install github.com/webx-top/tower
tower