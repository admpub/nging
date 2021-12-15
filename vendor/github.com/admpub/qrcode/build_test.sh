export GOOS=linux
export GOARCH=amd64
# export CGO_ENABLED=1
go build -tags zbar -ldflags "-linkmode external -extldflags -static"