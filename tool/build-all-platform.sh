go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
cd ..
go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool

export NGINGEX=

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build-tar.sh


export GOOS=linux
export GOARCH=386
source ${PWD}/inc-build-tar.sh

export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build-tar.sh



export NGINGEX=.exe

export GOOS=windows
export GOARCH=386
source ${PWD}/inc-build-tar.sh

export GOOS=windows
export GOARCH=amd64
source ${PWD}/inc-build-tar.sh
