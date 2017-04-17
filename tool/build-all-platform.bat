go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
cd ..
go-bindata-assetfs -tags bindata public/... template/...
cd tool

set NGINGEX=

set GOOS=linux
set GOARCH=amd64
call inc-build-zip.bat

set GOOS=linux
set GOARCH=386
call inc-build-zip.bat

set GOOS=darwin
set GOARCH=amd64
call inc-build-zip.bat



set NGINGEX=.exe

set GOOS=windows
set GOARCH=386
call inc-build-zip.bat

set GOOS=windows
set GOARCH=amd64
call inc-build-zip.bat

pause