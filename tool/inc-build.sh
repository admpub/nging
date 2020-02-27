mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata sqlite${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL}" -o ../dist/nging_${GOOS}_${GOARCH}/nging${NGINGEX} ..
mkdir ../dist/nging_${GOOS}_${GOARCH}/data
mkdir ../dist/nging_${GOOS}_${GOARCH}/data/logs
cp -R ../data/ip2region ../dist/nging_${GOOS}_${GOARCH}/data/ip2region


mkdir ../dist/nging_${GOOS}_${GOARCH}/config
mkdir ../dist/nging_${GOOS}_${GOARCH}/config/vhosts

#cp -R ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.* ../dist/nging_${GOOS}_${GOARCH}/config/
cp -R ../config/preupgrade.* ../dist/nging_${GOOS}_${GOARCH}/config/
cp -R ../config/ua.txt ../dist/nging_${GOOS}_${GOARCH}/config/ua.txt

if [ $GOOS = "windows" ]; then
    cp -R ../support/sqlite3_${GOARCH}.dll ../dist/nging_${GOOS}_${GOARCH}/sqlite3_${GOARCH}.dll
	export archiver_extension=zip
else
	export archiver_extension=zip
fi

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/

archiver make ../dist/nging_${GOOS}_${GOARCH}.${archiver_extension} ../dist/nging_${GOOS}_${GOARCH}/

rm -rf ../dist/nging_${GOOS}_${GOARCH}
