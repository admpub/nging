
#go install github.com/webx-top/echo
export OSARCH=`go env GOOS`_`go env GOARCH`
go tool compile -I $GOPATH/pkg/$OSARCH ./*.go
mv ${PWD}/main.o ${PWD}/../example.o