mkdir ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
xgo -go=1.14.4 -image=crazymax/xgo -targets=${GOOS}/${GOARCH} -dest=../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH} -tags="bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL}" ../
mv ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/${NGING_EXECUTOR}-${GOOS}-* ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/${NGING_EXECUTOR}${NGINGEX}
mkdir ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data
mkdir ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data/logs
cp -R ../data/ip2region ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/data/ip2region

mkdir ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config
mkdir ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/vhosts

#cp -R ../config/config.yaml ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.* ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/
cp -R ../config/preupgrade.* ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/
cp -R ../config/ua.txt ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/config/ua.txt

export archiver_extension=zip

cp -R ../dist/default/* ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/

rm -rf ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}.${archiver_extension}

#${NGING_VERSION}${NGING_LABEL}
arc archive ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}.${archiver_extension} ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}/

rm -rf ../dist/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
