#source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool

source ${PWD}/inc-version.sh

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=arm

source ${PWD}/inc-build.sh
