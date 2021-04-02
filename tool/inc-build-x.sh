export DISTPATH=${PKGPATH}/dist
export OSVERSIONDIR=${NGING_EXECUTOR}_${GOOS}_${GOARCH}
export RELEASEDIR=${DISTPATH}/${OSVERSIONDIR}
export LDFLAGS="-extldflags '-static'"
mkdir ${RELEASEDIR}

case "$GOARCH" in
    "arm"|"arm64"|"arm-7"|"arm-6"|"arm-5")
        export LDFLAGS="-extldflags '-static'"
        ;;
    *)
        export LDFLAGS=""
esac

xgo -go=1.16.2 -goproxy=https://goproxy.cn,direct -image=crazymax/xgo -targets=${GOOS}/${GOARCH} -dest=${RELEASEDIR} -out=${NGING_EXECUTOR} -tags="bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL} ${LDFLAGS}" ./${PKGPATH}

mv ${RELEASEDIR}/${NGING_EXECUTOR}-${GOOS}-* ${RELEASEDIR}/${NGING_EXECUTOR}${NGINGEX}
mkdir ${RELEASEDIR}/data
mkdir ${RELEASEDIR}/data/logs
cp -R ${PKGPATH}/data/ip2region ${RELEASEDIR}/data/ip2region

mkdir ${RELEASEDIR}/config
mkdir ${RELEASEDIR}/config/vhosts

#cp -R ../config/config.yaml ${RELEASEDIR}/config/config.yaml
cp -R ${PKGPATH}/config/config.yaml.sample ${RELEASEDIR}/config/config.yaml.sample
cp -R ${PKGPATH}/config/install.* ${RELEASEDIR}/config/
cp -R ${PKGPATH}/config/preupgrade.* ${RELEASEDIR}/config/
cp -R ${PKGPATH}/config/ua.txt ${RELEASEDIR}/config/ua.txt


cp -R ${DISTPATH}/default/* ${RELEASEDIR}/

export archiver_extension="tar.gz"

rm -rf ${RELEASEDIR}.${archiver_extension}

#${NGING_VERSION}${NGING_LABEL}

tar -zcvf ${RELEASEDIR}.${archiver_extension} -C ${DISTPATH} ${OSVERSIONDIR}
# 解压: tar -zxvf nging_linux_amd64.tar.gz -C ./nging_linux_amd64

rm -rf ${RELEASEDIR}
