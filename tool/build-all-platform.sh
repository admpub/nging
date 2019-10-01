source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool
export NGING_VERSION="2.0.5"
export NGING_BUILD=`date +%Y%m%d%H%M%S`
export NGING_COMMIT=`git rev-parse HEAD`
export NGING_LABEL="stable"

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build.sh


export GOOS=linux
export GOARCH=386
source ${PWD}/inc-build.sh

export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build.sh


export GOOS=freebsd
export GOARCH=386
source ${PWD}/inc-build.sh

export GOOS=freebsd
export GOARCH=amd64
source ${PWD}/inc-build.sh

export GOOS=linux
export GOARCH=arm
source ${PWD}/inc-build.sh

export GOOS=linux
export GOARCH=arm64
source ${PWD}/inc-build.sh

export GOOS=linux
export GOARCH=mips
source ${PWD}/inc-build.sh

export GOOS=linux
export GOARCH=mips64
source ${PWD}/inc-build.sh


export GOOS=linux
export GOARCH=mipsle
source ${PWD}/inc-build.sh

export GOOS=linux
export GOARCH=mips64le
source ${PWD}/inc-build.sh


# windows 放到最后
export NGINGEX=.exe
export BUILDTAGS=" windll"

export GOOS=windows
export GOARCH=386
source ${PWD}/inc-build.sh

export GOOS=windows
export GOARCH=amd64
source ${PWD}/inc-build.sh
