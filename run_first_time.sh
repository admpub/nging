#!/bin/sh
sudo launchctl limit maxfiles 65535
ulimit -n 65535
go install github.com/webx-top/tower@latest
tower