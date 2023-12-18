#go install github.com/admpub/nging-builder@latest
if [ "$1" != "" ]; then
nging-builder --conf ./builder-$1.conf darwin_amd64 min
exit 0
fi
nging-builder --conf ./builder.conf darwin_amd64 min
