go install github.com/admpub/nging-builder@latest
#go install github.com/admpub/xgo@latest
platform=""
if [ "$1" != "" ]; then
    platform="$1"
fi
nging-builder --conf ./builder.conf $platform min
