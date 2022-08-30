go install github.com/webx-top/db/cmd/dbgenerator@latest
dbgenerator -h 127.0.0.1 -d nging -p root -o ./application/dbschema -match "^nging_ftp_" -backup "./pkg/library/setup/install.sql" -charset utf8mb4
