mkdir ..\dist\nging_%GOOS%_%GOARCH%
go build -tags "bindata sqlite%BUILDTAGS%" -o ../dist/nging_%GOOS%_%GOARCH%/nging ..

mkdir ..\dist\nging_%GOOS%_%GOARCH%\data
mkdir ..\dist\nging_%GOOS%_%GOARCH%\data\logs
xcopy ..\data\ip2region ..\dist\nging_%GOOS%_%GOARCH%\data\ip2region /E /Q /H /I /Y


mkdir ..\dist\nging_%GOOS%_%GOARCH%\config
mkdir ..\dist\nging_%GOOS%_%GOARCH%\config\vhosts

xcopy ..\config\config.yaml ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\config.yaml.sample ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\install.sql ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\ua.txt ..\dist\nging_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y

if "%GOOS%"=="windows" (xcopy ..\support\sqlite3_%GOARCH%.dll ..\dist\nging_%GOOS%_%GOARCH%\ /E /Q /H /I /Y)

xcopy ..\dist\default ..\dist\nging_%GOOS%_%GOARCH%\ /E /Q /H /I /Y

if "%GOOS%"=="windows" (set archiver_extension=zip) else (set archiver_extension=tar.bz2)

archiver make ..\dist\nging_%GOOS%_%GOARCH%.%archiver_extension% ..\dist\nging_%GOOS%_%GOARCH%\

rd /s /Q ..\dist\nging_%GOOS%_%GOARCH%