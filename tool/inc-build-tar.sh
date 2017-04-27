
mkdir ../dist/nging_${GOOS}_${GOARCH}
go build -tags "bindata" -o ../dist/nging_${GOOS}_${GOARCH}/nging_${GOOS}_${GOARCH}${NGINGEX} ..
cp -R ../data ../dist/nging_${GOOS}_${GOARCH}/data
cp -R ../config ../dist/nging_${GOOS}_${GOARCH}/config
cp -R ../dist/default/* ../dist/nging_${GOOS}_${GOARCH}/
cd ../dist/nging_${GOOS}_${GOARCH}
tar -zcvf ../nging_${GOOS}_${GOARCH}.tar.gz ./*
cd ../../tool
rm -rf ../dist/nging_${GOOS}_${GOARCH}
