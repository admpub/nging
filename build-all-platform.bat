set GOOS=linux
set GOARCH=amd64
go build -o dist/nging_%GOOS%_%GOARCH% ./

set GOOS=linux
set GOARCH=386
go build -o dist/nging_%GOOS%_%GOARCH% ./

set GOOS=windows
set GOARCH=386
go build -o dist/nging_%GOOS%_%GOARCH%.exe ./ 

set GOOS=windows
set GOARCH=amd64
go build -o dist/nging_%GOOS%_%GOARCH%.exe ./

set GOOS=darwin
set GOARCH=386
go build -o dist/nging_%GOOS%_%GOARCH% ./

set GOOS=darwin
set GOARCH=amd64
go build -o dist/nging_%GOOS%_%GOARCH% ./
pause