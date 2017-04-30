set GOOS=windows
set GOARCH=386
go test
set GOOS=windows
set GOARCH=amd64
go test
pause