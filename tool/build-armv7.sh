#source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool

source ${PWD}/inc-version.sh

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARM=7
export GOARCH=arm

# export LDFLAGS="-extldflags '-static'"
# export CGO_ENABLED=1

source ${PWD}/inc-build.sh
