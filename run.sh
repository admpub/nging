#!/bin/sh
sudo launchctl limit maxfiles 655350
ulimit -n 655350
tower