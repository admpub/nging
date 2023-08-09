cd C:\Users\PC\go\src\github.com\admpub\goforever
go test -v -count=1 -run "TestProcessStartByUser" --user=hank-minipc\test
@REM go test -v -count=1 -run "TestWindowsToken"
@REM go test -v -count=1 -run "TestGetTokenByPid"
@REM runas /user:hank-minipc\test "C:\Users\test\example.exe"
pause