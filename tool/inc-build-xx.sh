declare -a targetsArgValue=""
declare -a targetsArgSep=""
for target in ${NGING_TARGETS[@]}
do
    targetsArgValue="${targetsArgValue}${targetsArgSep}${target}"
    targetsArgSep=","
done
echo "⛵️ building: $targetsArgValue"
if [[ "$targetsArgValue" = "" ]]; then
    exit 1;
fi

export DISTPATH=${PKGPATH}/dist
export RELEASE_TEMPDIR=${DISTPATH}/${NGING_EXECUTOR}
export LDFLAGS="-extldflags '-static'"
mkdir ${RELEASE_TEMPDIR}

xgo -go=${GO_VERSION} -goproxy="https://goproxy.cn,direct" -image="crazymax/xgo:${GO_VERSION}" -targets="${targetsArgValue}" -dest="${RELEASE_TEMPDIR}" -out=${NGING_EXECUTOR} -tags="bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL} -X main.BUILD_OS=${GOOS} -X main.BUILD_ARCH=${GOARCH} ${MINIFYFLAG} ${LDFLAGS}" ./${PKGPATH}

pack(){
    export nging_filename="${1%.*}"
    export nging_filename="${nging_filename//-4.0/}"
    export nging_filename="${nging_filename//-/_}"
    export RELEASEDIR="${DISTPATH}/${nging_filename}"
    mkdir ${RELEASEDIR}
    export nging_extension="${1##*.}"
    if [ "$nging_extension" != "" ] && [ "$nging_extension" != "$1" ]; then
        export nging_extension=".${nging_extension}"
    else
        export nging_extension=""
    fi
    mv "${RELEASE_TEMPDIR}/$1" "${RELEASEDIR}/${NGING_EXECUTOR}${nging_extension}"
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

    mkdir ${RELEASEDIR}/public
    mkdir ${RELEASEDIR}/public/upload
    cp -R ${PKGPATH}/public/upload/.gitkeep ${RELEASEDIR}/public/upload/.gitkeep

    export archiver_extension="tar.gz"

    rm -rf ${RELEASEDIR}.${archiver_extension}

    #${NGING_VERSION}${NGING_LABEL}

    tar -zcvf ${RELEASEDIR}.${archiver_extension} -C ${DISTPATH} ${nging_filename}
    # 解压: tar -zxvf nging_linux_amd64.tar.gz -C ./nging_linux_amd64

    rm -rf ${RELEASEDIR}
}

binFiles=`ls ${RELEASE_TEMPDIR}`
for file in $binFiles
do 
 $(pack $file)
done

rm -rf ${RELEASE_TEMPDIR}
