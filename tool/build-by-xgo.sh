go get github.com/karalabe/xgo
go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
$GOPATH/bin/go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool
export NGING_VERSION="1.3.3"
export NGING_BUILD=`date +%Y%m%d%H%M%S`
export NGING_COMMIT=`git rev-parse HEAD`
export NGING_LABEL="beta"

export NGINGEX=
export BUILDTAGS=" official zbar"

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build-x.sh


export GOOS=linux
export GOARCH=386
source ${PWD}/inc-build-x.sh

export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build-x.sh



export NGINGEX=.exe

export GOOS=windows
export GOARCH=386
source ${PWD}/inc-build-x.sh

export GOOS=windows
export GOARCH=amd64
source ${PWD}/inc-build-x.sh
