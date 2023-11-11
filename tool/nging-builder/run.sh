go install github.com/admpub/nging-builder@v0.2.1
platform=""
if [ "$1" != "" ]; then
    platform="$1"
fi
nging-builder --conf ./builder.conf $platform min
