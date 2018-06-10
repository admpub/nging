mkdir -p $GOPATH/src/golang.org/x
go get github.com/golang/crypto
go get github.com/golang/sys
mv $GOPATH/src/github.com/golang/crypto $GOPATH/src/golang.org/x/crypto
mv $GOPATH/src/github.com/golang/sys $GOPATH/src/golang.org/x/sys
go get github.com/webx-top/tower
export PATH="$PATH:$GOPATH/bin"
go get github.com/karalabe/xgo
go get github.com/jteeuwen/go-bindata/...
go get github.com/admpub/go-bindata-assetfs/...
cd $GOPATH/src/github.com/nging
$GOPATH/bin/go-bindata-assetfs -tags bindata public/... template/... config/i18n/...
