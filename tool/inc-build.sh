export DISTPATH=../dist
export RELEASEDIR=${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
if [ "$GOARM" != "" ]; then
	export RELEASEDIR=${RELEASEDIR}v${GOARM}
fi
mkdir ${RELEASEDIR}
go build -tags "bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL}" -o ${RELEASEDIR}/${NGING_EXECUTOR}${NGINGEX} ..
mkdir ${RELEASEDIR}/data
mkdir ${RELEASEDIR}/data/logs
cp -R ../data/ip2region ${RELEASEDIR}/data/ip2region


mkdir ${RELEASEDIR}/config
mkdir ${RELEASEDIR}/config/vhosts

#cp -R ../config/config.yaml ${RELEASEDIR}/config/config.yaml
cp -R ../config/config.yaml.sample ${RELEASEDIR}/config/config.yaml.sample
cp -R ../config/install.* ${RELEASEDIR}/config/
cp -R ../config/preupgrade.* ${RELEASEDIR}/config/
cp -R ../config/ua.txt ${RELEASEDIR}/config/ua.txt

if [ "$GOOS" = "windows" ]; then
    cp -R ../support/sqlite3_${GOARCH}.dll ${RELEASEDIR}/sqlite3_${GOARCH}.dll
	export archiver_extension=zip
else
	export archiver_extension=zip
fi

cp -R ../dist/default/* ${RELEASEDIR}/

arc archive ${RELEASEDIR}.${archiver_extension} ${RELEASEDIR}

rm -rf ${RELEASEDIR}
