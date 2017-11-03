mkdir ..\dist\nging_%GOOS%_%GOARCH%
go build -tags "bindata%BUILDTAGS%" -o ../dist/nging_%GOOS%_%GOARCH%/nging_%GOOS%_%GOARCH%%NGINGEX% ..

xcopy ..\data ..\dist\nging_%GOOS%_%GOARCH%\data /E /Q /H /I /Y
mkdir ..\dist\nging_%GOOS%_%GOARCH%\config
mkdir ..\dist\nging_%GOOS%_%GOARCH%\config\vhosts

xcopy ..\config\config.yaml ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\config.yaml.sample ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\install.sql ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y

if "%GOOS%"=="windows" (xcopy ..\support\sqlite3_%GOARCH%.dll ..\dist\nging_%GOOS%_%GOARCH%\ /E /Q /H /I /Y)

xcopy ..\dist\default ..\dist\nging_%GOOS%_%GOARCH%\ /E /Q /H /I /Y

rem if "%GOOS%"=="windows" (set archiver_extension=zip) else (set archiver_extension=tar.bz2)

rem archiver make ..\dist\nging_%GOOS%_%GOARCH%.%archiver_extension% ..\dist\nging_%GOOS%_%GOARCH%\

rem rd /s /Q ..\dist\nging_%GOOS%_%GOARCH%