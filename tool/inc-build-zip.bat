mkdir ..\dist\nging_%GOOS%_%GOARCH%
go build -tags "bindata" -o ../dist/nging_%GOOS%_%GOARCH%/nging_%GOOS%_%GOARCH%%NGINGEX% ..

xcopy ..\data ..\dist\nging_%GOOS%_%GOARCH%\data /E /Q /H /I /Y
xcopy ..\config ..\dist\nging_%GOOS%_%GOARCH%\config /E /Q /H /I /Y
xcopy ..\dist\default ..\dist\nging_%GOOS%_%GOARCH%\ /E /Q /H /I /Y

7zr.exe a ..\dist\nging_%GOOS%_%GOARCH%.zip ..\dist\nging_%GOOS%_%GOARCH%\* -r
rd /s /Q ..\dist\nging_%GOOS%_%GOARCH%