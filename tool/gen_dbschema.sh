go get github.com/webx-top/db
go install github.com/webx-top/db/_tools/generator
generator -h 127.0.0.1 -d nging -p root -o ../application/dbschema -ignore "^official_" -backup "../config/install.sql" -charset utf8mb4
