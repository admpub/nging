FROM crazymax/xgo:1.23.0
ARG GO_VERSION
ENV GO_VERSION ${GO_VERSION:-1.24.0}
RUN curl https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -o go.linux-amd64.tar.gz && \
 rm -rf /usr/local/go && \
 tar -zxvf go.linux-amd64.tar.gz -C /usr/local && \
 rm -rf ./go.linux-amd64.tar.gz
 
# docker build -t admpub/xgo:1.24.0 . --build-arg GO_VERSION=1.24.0
