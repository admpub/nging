go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
cd ..
go-bindata-assetfs -tags bindata public/... template/...
cd tool

set NGINGEX=

set GOOS=linux
set GOARCH=amd64
source ${PWD}/inc-build-tar.sh


set GOOS=linux
set GOARCH=386
source ${PWD}/inc-build-tar.sh

set GOOS=darwin
set GOARCH=amd64
source ${PWD}/inc-build-tar.sh



set NGINGEX=.exe

set GOOS=windows
set GOARCH=386
source ${PWD}/inc-build-tar.sh

set GOOS=windows
set GOARCH=amd64
source ${PWD}/inc-build-tar.sh
