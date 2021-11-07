go get github.com/webx-top/db
go install github.com/webx-top/db/cmd/dbgenerator
dbgenerator -h 127.0.0.1 -d nging -p root -o ../application/dbschema -match "^nging_" -backup "../config/install.sql" -charset utf8mb4
