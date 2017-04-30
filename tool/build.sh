
# set GOPATH GOROOT
export GOPATH="/home/admpub/go"
export GOROOT="/usr/local/go"

# set PATH so it includes user's private bin directories
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
$GOPATH/bin/go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build.sh

