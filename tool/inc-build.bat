mkdir ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%
go build -tags "bindata sqlite%BUILDTAGS%" -ldflags="-X main.BUILD_TIME=%NGING_BUILD% -X main.COMMIT=%NGING_COMMIT% -X main.VERSION=%NGING_VERSION% -X main.LABEL=%NGING_LABEL%" -o ../dist/%NGING_EXECUTOR%_%GOOS%_%GOARCH%/%NGING_EXECUTOR%%NGINGEX% ..

mkdir ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\data
mkdir ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\data\logs
xcopy ..\data\ip2region ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\data\ip2region /E /Q /H /I /Y


mkdir ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\config
mkdir ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\config\vhosts

xcopy ..\config\config.yaml.sample ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\install.sql ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y
xcopy ..\config\ua.txt ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\config\ /E /Q /H /I /Y

if "%GOOS%"=="windows" (xcopy ..\support\sqlite3_%GOARCH%.dll ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\ /E /Q /H /I /Y)

xcopy ..\dist\default ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%\ /E /Q /H /I /Y

set archiver_extension=zip
rem if "%GOOS%"=="windows" (set archiver_extension=zip) else (set archiver_extension=tar.bz2)

arc archive ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%.%archiver_extension% ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%

rd /s /Q ..\dist\%NGING_EXECUTOR%_%GOOS%_%GOARCH%