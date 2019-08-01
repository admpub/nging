rem 本脚本在编译非windows下的程序时，将忽略sqlite
rem 如要支持sqlite，请在linux系统内执行build-all-platform.sh
call install-archiver.bat
cd ..
go generate
cd tool


set time_hh=%time:~0,2%
if /i %time_hh% LSS 10 (set time_hh=0%time:~1,1%)

set NGING_VERSION="2.0.3"
set NGING_BUILD=%date:~,4%%date:~5,2%%date:~8,2%%time_hh%%time:~3,2%%time:~6,2%
set NGING_COMMIT=
for /F %%i in ('git rev-parse HEAD') do ( set NGING_COMMIT=%%i)
set NGING_LABEL="stable"

set NGINGEX=
set BUILDTAGS= windll

set GOOS=linux
set GOARCH=amd64
call inc-build.bat

set GOOS=linux
set GOARCH=386
call inc-build.bat

set GOOS=darwin
set GOARCH=amd64
call inc-build.bat


set GOOS=windows
set GOARCH=386
call inc-build.bat


set NGINGEX=.exe
set BUILDTAGS=

set GOOS=windows
set GOARCH=amd64
call inc-build.bat

pause