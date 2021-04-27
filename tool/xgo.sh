source ${PWD}/inc-version.sh

#go install github.com/admpub/xgo
#source ${WORKDIR}/install-archiver.sh

cd ../
go generate

# 回到入口
cd ${ENTRYDIR}
export NGING_BUILDER="inc-build-xx.sh"
declare -a NGING_TARGETS=()
#echo $PWD && exit
reset() {
    export TMPDIR=
    export BUILDTAGS=
    export MINIFYFLAG=
}

open_minify(){
    export MINIFYFLAG="-s -w"
}

close_minify(){
    export MINIFYFLAG=""
}

setTarget(){
    NGING_TARGETS=(${NGING_TARGETS[*]} "${GOOS}/${GOARCH}")
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
    export GOOS=linux
    export GOARCH=amd64
    setTarget
}

linux_arm5() {
    export GOOS=linux
    export GOARCH=arm-5
    setTarget
}

linux_arm6() {
    export GOOS=linux
    export GOARCH=arm-6
    setTarget
}

linux_arm7() {
    export GOOS=linux
    export GOARCH=arm-7
    setTarget
}

linux_arm64() {
    export GOOS=linux
    export GOARCH=arm64
    setTarget
}

linux_386() {
    export GOOS=linux
    export GOARCH=386
    setTarget
}

darwin_amd64() {
    export GOOS=darwin
    export GOARCH=amd64
    setTarget
}

windows_386() {
    export GOOS=windows
    export GOARCH=386
    setTarget
}

windows_amd64() {
    export GOOS=windows
    export GOARCH=amd64
    setTarget
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
    "linux_386")
        linux_386
        ;;
    "darwin_amd64")
        darwin_amd64
        ;;
    "windows_386")
        windows_386
        ;;
    "windows_amd64")
        windows_amd64
        ;;
    "min"|"m")
        open_minify
        all
        ;;
    *)
        all
esac

source ${WORKDIR}/${NGING_BUILDER}