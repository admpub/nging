#source ${PWD}/install-archiver.sh
cd ..
go generate
cd tool

source ${PWD}/inc-version.sh

declare -a goBuilderScript=${PWD}/inc-build.sh

reset() {
    export NGINGEX=
    export BUILDTAGS=
    export MINIFYFLAG="-s -w" 
    # export CGO_ENABLED=1
}

linux_amd64() {
    reset
    export GOOS=linux
    export GOARCH=amd64
    source $goBuilderScript
}


linux_386() {
    reset
    export GOOS=linux
    export GOARCH=386
    source $goBuilderScript
}

darwin_amd64() {
    reset
    export GOOS=darwin
    export GOARCH=amd64
    source $goBuilderScript
}


freebsd_386() {
    reset
    export GOOS=freebsd
    export GOARCH=386
    source $goBuilderScript
}

freebsd_amd64() {
    reset
    export GOOS=freebsd
    export GOARCH=amd64
    source $goBuilderScript
}

linux_arm5() {
    reset
    export GOOS=linux
    export GOARM=5
    export GOARCH=arm
    source $goBuilderScript
}

linux_arm6() {
    reset
    export GOOS=linux
    export GOARM=6
    export GOARCH=arm
    source $goBuilderScript
}

linux_arm7() {
    reset
    export GOOS=linux
    export GOARM=7
    export GOARCH=arm
}

linux_arm64() {
    reset
    export GOOS=linux
    export GOARM=
    export GOARCH=arm64
    source $goBuilderScript
}

linux_mips() {
    reset
    export GOOS=linux
    export GOARCH=mips
    source $goBuilderScript
}

linux_mips64() {
    reset
    export GOOS=linux
    export GOARCH=mips64
    source $goBuilderScript
}


linux_mipsle() {
    reset
    export GOOS=linux
    export GOARCH=mipsle
    source $goBuilderScript
}

linux_mips64le() {
    reset
    export GOOS=linux
    export GOARCH=mips64le
    source $goBuilderScript
}


windows_set() {
    reset
    export NGINGEX=.exe
    export BUILDTAGS=" windll"
}

windows_386() {
    windows_set
    export GOOS=windows
    export GOARCH=386
    source $goBuilderScript
}

windows_amd64() {
    windows_set
    export GOOS=windows
    export GOARCH=amd64
    source $goBuilderScript
}

windows_arm64() {
    windows_set
    export GOOS=windows
    export GOARM=
    export GOARCH=amd64
    source $goBuilderScript
}


all() {
    linux_amd64
    linux_arm5
    linux_arm6
    linux_arm7
    linux_arm64
    linux_386

    freebsd_amd64
    freebsd_386
    linux_mips
    linux_mips64
    linux_mipsle
    linux_mips64le

    darwin_amd64
    windows_386
    windows_amd64
    windows_arm64
}

case "$1" in
    "linux_amd64")
        linux_amd64
        ;;
    "linux_arm5")
        linux_arm5
        ;;
    "linux_arm6")
        linux_arm6
        ;;
    "linux_arm7")
        linux_arm7
        ;;
    "linux_arm64")
        linux_arm64
        ;;
    "linux_arm*")
        linux_arm5
        linux_arm6
        linux_arm7
        linux_arm64
        ;;
    "linux_386")
        linux_386
        ;;
    "freebsd_amd64")
        freebsd_amd64
        ;;
    "freebsd_386")
        freebsd_386
        ;;
    "linux_mips")
        linux_mips
        ;;
    "linux_mips64")
        linux_mips64
        ;;
    "linux_mipsle")
        linux_mipsle
        ;;
    "linux_mips64le")
        linux_mips64le
        ;;

    "darwin_amd64")
        darwin_amd64
        ;;
    "windows*")
        windows_386
        windows_amd64
        ;;
    "windows_386")
        windows_386
        ;;
    "windows_amd64")
        windows_amd64
        ;;
    "windows_arm64")
        windows_arm64
        ;;
    "")
        all
        ;;
    *)
        echo "Unknown option $1"
        exit 1
        ;;
esac