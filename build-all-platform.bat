go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
go-bindata-assetfs -tags bindata public/... template/...

set GOOS=linux
set GOARCH=amd64
go build -tags "bindata" -o dist/nging_%GOOS%_%GOARCH% ./

set GOOS=linux
set GOARCH=386
go build -tags "bindata" -o dist/nging_%GOOS%_%GOARCH% ./

set GOOS=windows
set GOARCH=386
go build -tags "bindata" -o dist/nging_%GOOS%_%GOARCH%.exe ./ 

set GOOS=windows
set GOARCH=amd64
go build -tags "bindata" -o dist/nging_%GOOS%_%GOARCH%.exe ./

set GOOS=darwin
set GOARCH=amd64
go build -tags "bindata" -o dist/nging_%GOOS%_%GOARCH% ./

xcopy data "dist/data" /E /Q /H /I /Y
xcopy config "dist/config" /E /Q /H /I /Y

pause