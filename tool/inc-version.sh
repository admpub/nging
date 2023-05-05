# Env
export PKGPATH=github.com/admpub/nging
export ENTRYDIR=${GOPATH}/src
export WORKDIR=${PWD}

# Go configuration
export GO_VERSION="1.20.4"

# Nging configuration
export NGING_VERSION="5.0.1"
export NGING_BUILD=`date +%Y%m%d%H%M%S`
export NGING_COMMIT=`git rev-parse HEAD`
export NGING_LABEL="beta1"
export NGING_EXECUTOR="nging"