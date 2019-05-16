mkdir ../dist/nging_${GOOS}_${GOARCH}
xgo -go=1.12 -image=admpub/xgo:master -targets=${GOOS}/${GOARCH} -dest=../dist/nging_${GOOS}_${GOARCH} -tags="bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL}" ../
mv ../dist/nging_${GOOS}_${GOARCH}/nging-${GOOS}-${GOARCH}${NGINGEX} ../dist/nging_${GOOS}_${GOARCH}/nging${NGINGEX}
mkdir ../dist/nging_${GOOS}_${GOARCH}/data
mkdir ../dist/nging_${GOOS}_${GOARCH}/data/logs
cp -R ../data/ip2region ../dist/nging_${GOOS}_${GOARCH}/data/ip2region

mkdir ../dist/nging_${GOOS}_${GOARCH}/config
mkdir ../dist/nging_${GOOS}_${GOARCH}/config/vhosts

#cp -R ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.sql ../dist/nging_${GOOS}_${GOARCH}/config/install.sql
cp -R ../config/ua.txt ../dist/nging_${GOOS}_${GOARCH}/config/ua.txt

export archiver_extension=zip

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/
#${NGING_VERSION}${NGING_LABEL}
archiver make ../dist/nging_${GOOS}_${GOARCH}.${archiver_extension} ../dist/nging_${GOOS}_${GOARCH}/

rm -rf ../dist/nging_${GOOS}_${GOARCH}
