export DISTPATH=${PKGPATH}/dist
mkdir ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
xgo -go=latest -goproxy=https://goproxy.cn,direct -image=admpub/xgo -targets=${GOOS}/${GOARCH} -dest=${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH} -tags="bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL}" ./${PKGPATH}
mv ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/${NGING_EXECUTOR}-${GOOS}-* ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/${NGING_EXECUTOR}${NGINGEX}
mkdir ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data
mkdir ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data/logs
cp -R ${PKGPATH}/data/ip2region ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data/ip2region

mkdir ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config
mkdir ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/vhosts

#cp -R ../config/config.yaml ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/config.yaml
cp -R ${PKGPATH}/config/config.yaml.sample ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ${PKGPATH}/config/install.* ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/
cp -R ${PKGPATH}/config/preupgrade.* ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/
cp -R ${PKGPATH}/config/ua.txt ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/ua.txt

export archiver_extension=zip

cp -R ${DISTPATH}/default/* ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/

rm -rf ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}.${archiver_extension}

#${NGING_VERSION}${NGING_LABEL}
arc archive ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}.${archiver_extension} ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/

rm -rf ${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
