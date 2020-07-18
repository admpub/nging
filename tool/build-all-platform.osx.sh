#source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool

source ${PWD}/inc-version.sh

export NGINGEX=
export BUILDTAGS=

export CGO_ENABLED=1
export CC_32=/usr/local/gcc-4.8.1-for-linux32/bin/i586-pc-linux-gcc
export CC_64=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc
export CC_DEFAULT=clang

export CC=$CC_32
export GOOS=linux
export GOARCH=386
source ${PWD}/inc-build.sh

exit 0;

export CC=$CC_64
export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build.sh

#exit 0;

export CC=$CC_DEFAULT
export GOOS=darwin
export GOARCH=amd64
source ${PWD}/inc-build.sh

export CGO_ENABLED=

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
