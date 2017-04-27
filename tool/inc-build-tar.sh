
mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata" -o ../dist/nging_${GOOS}_${GOARCH}/nging_${GOOS}_${GOARCH}${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data
cp -R ../config ../dist/nging_${GOOS}_${GOARCH}/config
cp -R ../dist/default ../dist/nging_${GOOS}_${GOARCH}/

tar -zcvf ../dist/nging_${GOOS}_${GOARCH}.tar.gz ../dist/nging_${GOOS}_${GOARCH}/*
rm -rf ../dist/nging_${GOOS}_${GOARCH}
