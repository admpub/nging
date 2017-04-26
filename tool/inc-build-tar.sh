
mkdir ../dist/nging_$GOOS_$GOARCH
go build -tags "bindata" -o ../dist/nging_$GOOS_$GOARCH/nging_$GOOS_$GOARCH%NGINGEX% ..
cp -R ../data ../dist/nging_$GOOS_$GOARCH/data
cp -R ../config ../dist/nging_$GOOS_$GOARCH/config
cp -R ../dist/default ../dist/nging_$GOOS_$GOARCH/

tar -zcvf ../dist/nging_$GOOS_$GOARCH.tar.gz ../dist/nging_$GOOS_$GOARCH/*
rm -rf ../dist/nging_$GOOS_$GOARCH
