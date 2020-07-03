source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool

source ${PWD}/inc-version.sh

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


# export GOOS=freebsd
# export GOARCH=386
# source ${PWD}/inc-build.sh

# export GOOS=freebsd
# export GOARCH=amd64
# source ${PWD}/inc-build.sh

# export GOOS=linux
# export GOARCH=arm
# source ${PWD}/inc-build.sh

# export GOOS=linux
# export GOARCH=arm64
# source ${PWD}/inc-build.sh

# export GOOS=linux
# export GOARCH=mips
# source ${PWD}/inc-build.sh

# export GOOS=linux
# export GOARCH=mips64
# source ${PWD}/inc-build.sh


# export GOOS=linux
# export GOARCH=mipsle
# source ${PWD}/inc-build.sh

# export GOOS=linux
# export GOARCH=mips64le
# source ${PWD}/inc-build.sh


# windows 放到最后
export NGINGEX=.exe
export BUILDTAGS=" windll"

export GOOS=windows
export GOARCH=386
source ${PWD}/inc-build.sh

export GOOS=windows
export GOARCH=amd64
source ${PWD}/inc-build.sh
