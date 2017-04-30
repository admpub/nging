mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata sqlite windll" -o ../dist/nging_${GOOS}_${GOARCH}/nging_${GOOS}_${GOARCH}${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data

mkdir ../dist/nging_${GOOS}_${GOARCH}/config
mkdir ../dist/nging_${GOOS}_${GOARCH}/config/vhosts

cp -R ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.sql ../dist/nging_${GOOS}_${GOARCH}/config/install.sql

if [ $GOOS = "windows" ]; then
    cp -R ../support/sqlite3_${GOARCH}.dll ../dist/nging_${GOOS}_${GOARCH}/sqlite3_${GOARCH}.dll
	export archiver_extension=tar.gz
else
	export archiver_extension=zip
fi

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/

archiver make ../nging_${GOOS}_${GOARCH}.${archiver_extension} ../dist/nging_${GOOS}_${GOARCH}/

rm -rf ../dist/nging_${GOOS}_${GOARCH}
