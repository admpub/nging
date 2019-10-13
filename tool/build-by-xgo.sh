go get github.com/karalabe/xgo
source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool
export NGING_VERSION="2.0.6"
export NGING_BUILD=`date +%Y%m%d%H%M%S`
export NGING_COMMIT=`git rev-parse HEAD`
export NGING_LABEL="stable"

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build-x.sh


export GOOS=linux
export GOARCH=386
source ${PWD}/inc-build-x.sh

export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build-x.sh



export NGINGEX=.exe

export GOOS=windows
export GOARCH=386
source ${PWD}/inc-build-x.sh

export GOOS=windows
export GOARCH=amd64
source ${PWD}/inc-build-x.sh
