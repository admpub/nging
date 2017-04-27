
mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata sqlite" -o ../dist/nging_${GOOS}_${GOARCH}/nging_${GOOS}_${GOARCH}${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data

mkdir ../dist/nging_${GOOS}_${GOARCH}/config
mkdir ../dist/nging_${GOOS}_${GOARCH}/config/vhost

cp -R ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp -R ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp -R ../config/install.sql ../dist/nging_${GOOS}_${GOARCH}/config/install.sql

if [ ${GOOS} = "windows" ] then
    cp -R support/sqlite3_${GOARCH}.dll ../dist/nging_${GOOS}_${GOARCH}/sqlite3_${GOARCH}.dll
fi

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/
cd ../dist/nging_${GOOS}_${GOARCH}
tar -zcvf ../nging_${GOOS}_${GOARCH}.tar.gz ./*
cd ../../tool
rm -rf ../dist/nging_${GOOS}_${GOARCH}
