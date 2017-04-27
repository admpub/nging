
mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata" -o ../dist/nging_${GOOS}_${GOARCH}/nging_${GOOS}_${GOARCH}${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data

cp ../config/config.yaml ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml
cp ../config/config.yaml.sample ../dist/nging_${GOOS}_${GOARCH}/config/config.yaml.sample
cp ../config/install.sql ../dist/nging_${GOOS}_${GOARCH}/config/install.sql
cp -R ../config/vhost ../dist/nging_${GOOS}_${GOARCH}/config/vhost

cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/
cd ../dist/nging_${GOOS}_${GOARCH}
tar -zcvf ../nging_${GOOS}_${GOARCH}.tar.gz ./*
cd ../../tool
rm -rf ../dist/nging_${GOOS}_${GOARCH}
