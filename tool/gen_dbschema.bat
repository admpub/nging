go get github.com/webx-top/db
go install github.com/webx-top/db/cmd/dbgenerator
dbgenerator -d nging -p root -o ../application/dbschema -ignore "^official_" -backup "../config/install.sql" -charset utf8mb4
pause