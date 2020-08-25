source ${PWD}/inc-version.sh

#go install github.com/admpub/xgo
#source ${WORKDIR}/install-archiver.sh

cd ../
go generate

# 回到入口
cd ${ENTRYDIR}

#echo $PWD && exit

export TMPDIR=

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${WORKDIR}/inc-build-x.sh


export GOOS=linux
export GOARCH=386
source ${WORKDIR}/inc-build-x.sh

export GOOS=darwin
export GOARCH=amd64
source ${WORKDIR}/inc-build-x.sh



export NGINGEX=.exe

export GOOS=windows
export GOARCH=386
source ${WORKDIR}/inc-build-x.sh

export GOOS=windows
export GOARCH=amd64
source ${WORKDIR}/inc-build-x.sh
