@echo off
go install github.com/akavel/rsrc@latest
rsrc.exe -manifest ico.manifest -o ../nging.syso -ico ../public/assets/backend/images/favicon-lg.ico
