go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd ..
%GOPATH%\bin\go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
cd tool

set NGINGEX=.exe
set GOOS=windows
set GOARCH=386
call inc-build-zip.bat

pause