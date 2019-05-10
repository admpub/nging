go get github.com/karalabe/xgo
go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
$GOPATH/bin/go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool
export NGING_VERSION="1.5.1"
export NGING_BUILD=`date +%Y%m%d%H%M%S%Z`
export NGING_COMMIT=`git rev-parse HEAD`
export NGING_LABEL=`beta`

export NGINGEX=
export BUILDTAGS=

export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build-x.sh
