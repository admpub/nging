GO_VERSION=1.24.3
if [ "$1" != "" ];then
    GO_VERSION="$1"
fi
docker build -t admpub/xgo:${GO_VERSION} . --build-arg GO_VERSION=${GO_VERSION}
if [ "$2" = "push" ];then
    docker push admpub/xgo:${GO_VERSION}
    #docker tag admpub/xgo:${GO_VERSION} admpub/xgo:latest
    #docker push admpub/xgo:latest
fi
