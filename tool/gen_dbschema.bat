go get github.com/webx-top/db
go install github.com/webx-top/db/_tools/generator
generator -d caddyui -p root -o ../application/dbschema
pause