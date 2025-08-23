# nging-builder
if which nging-builder &> /dev/null; then
    echo "found nging-builder"
else
    echo "installing nging-builder"
    go install github.com/admpub/nging-builder@latest
fi

# xgo
if which xgo &> /dev/null; then
    echo "found xgo"
else
    echo "installing xgo"
    go install github.com/admpub/xgo@latest
fi

# go-bindata
if which go-bindata &> /dev/null; then
    echo "found go-bindata"
else
    echo "installing go-bindata"
    go install github.com/admpub/bindata/v3/go-bindata@latest
fi

platform=""
if [ "$1" != "" ]; then
    platform="$1"
fi
nging-builder --conf ./builder.conf $platform min
