#echo $PWD && exit
reset() {
    export TMPDIR=
    export NGINGEX=
    export BUILDTAGS=
    export MINIFYFLAG=
}

open_minify(){
    export MINIFYFLAG="-s -w"
}

close_minify(){
    export MINIFYFLAG=""
}

all() {
    linux_amd64
    linux_arm5
    linux_arm6
    linux_arm7
    linux_arm64
    linux_386
    darwin_amd64
    windows_386
    windows_amd64
}

linux_amd64() {
    export NGINGEX=
    export GOOS=linux
    export GOARCH=amd64
    source $goBuilderScript
}

linux_arm5() {
    export NGINGEX=

    export GOOS=linux
    export GOARCH=arm-5
    source $goBuilderScript
}

linux_arm6() {
    export NGINGEX=

    export GOOS=linux
    export GOARCH=arm-6
    source $goBuilderScript
}

linux_arm7() {
    export NGINGEX=

    export GOOS=linux
    export GOARCH=arm-7
    source $goBuilderScript
}

linux_arm64() {
    export NGINGEX=

    export GOOS=linux
    export GOARCH=arm64
    source $goBuilderScript
}

linux_386() {
    export NGINGEX=

    export GOOS=linux
    export GOARCH=386
    source $goBuilderScript
}

darwin_amd64() {
    export NGINGEX=

    export GOOS=darwin
    export GOARCH=amd64
    source $goBuilderScript
}

windows_386() {
    export NGINGEX=.exe

    export GOOS=windows
    export GOARCH=386
    source $goBuilderScript
}

windows_amd64() {
    export NGINGEX=.exe

    export GOOS=windows
    export GOARCH=amd64
    source $goBuilderScript
}

windows_arm64() {
    export NGINGEX=.exe

    export GOOS=windows
    export GOARCH=arm64
    source $goBuilderScript
}

reset

case "$2" in
    "min"|"m")
    open_minify
    ;;
    *)
    close_minify
esac

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
    "darwin_amd64")
        darwin_amd64
        ;;
    "windows*")
        windows_386
        windows_amd64
#        windows_arm64
        ;;
    "windows_386")
        windows_386
        ;;
    "windows_amd64")
        windows_amd64
        ;;
#    "windows_arm64")
#        windows_arm64
#        ;;
    "min"|"m")
        open_minify
        all
        ;;
    "")
        all
        ;;
    *)
        echo "Unknown option $1"
        exit 1
        ;;
esac