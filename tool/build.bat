go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
%GOPATH%\bin\go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool

set NGINGEX=.exe
set BUILDTAGS= windll

set GOOS=windows
set GOARCH=amd64
call inc-build.bat

pause