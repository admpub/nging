source ${PWD}/inc-version.sh

#go install github.com/admpub/xgo
#source ${WORKDIR}/install-archiver.sh

declare -a goBuilderScript=${WORKDIR}/inc-build-x.sh

# 回到入口
cd ${ENTRYDIR}
source ${WORKDIR}/inc-func.sh