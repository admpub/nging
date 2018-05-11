mkdir ../dist/nging_${GOOS}_${GOARCH}
gox -tags "bindata sqlite" -osarch ${GOOS}/${GOARCH} -output ../dist/nging_${GOOS}_${GOARCH}/nging${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data

mkdir ../dist/nging_${GOOS}_${GOARCH}/config
mkdir ../dist/nging_${GOOS}_${GOARCH}/config/vhosts

cp -R ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.sql ../dist/nging_${GOOS}_${GOARCH}/config/install.sql

if [ $GOOS = "windows" ]; then
	export archiver_extension=zip
else
	export archiver_extension=tar.bz2
fi

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/

archiver make ../dist/nging_${GOOS}_${GOARCH}.${archiver_extension} ../dist/nging_${GOOS}_${GOARCH}/

rm -rf ../dist/nging_${GOOS}_${GOARCH}
