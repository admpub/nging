go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
%GOPATH%\bin\go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool

set NGINGEX=
set BUILDTAGS=

set GOOS=linux
set GOARCH=amd64
call inc-build_linux.bat

pause