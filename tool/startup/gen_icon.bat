@echo off
go install github.com/akavel/rsrc@latest
rsrc.exe -arch amd64 -manifest ico.manifest -o ./ico/windows/amd64/startup.syso -ico ../../public/assets/backend/images/favicon-lg.ico
rsrc.exe -arch arm64 -manifest ico.manifest -o ./ico/windows/arm64/startup.syso -ico ../../public/assets/backend/images/favicon-lg.ico
rsrc.exe -arch 386 -manifest ico.manifest -o ./ico/windows/i386/startup.syso -ico ../../public/assets/backend/images/favicon-lg.ico
