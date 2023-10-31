FROM ubuntu:latest
ARG GO_VERSION

ENV GO_VERSION ${GO_VERSION:-1.21.3}
ENV GO_INSTALL_FILE go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH "$PATH:/root/go/bin"
ENV GOPATH "/root/go"

RUN apt update 
#RUN apt upgrade -y
RUN apt install wget nginx telnet -y
RUN apt install certbot -y

#RUN wget -c https://studygolang.com/dl/golang/${GO_INSTALL_FILE}
RUN wget -c https://golang.google.cn/dl/${GO_INSTALL_FILE}
RUN tar zxvf ${GO_INSTALL_FILE}
RUN ln -s /go/bin/go /usr/local/bin/go
RUN ln -s /go/bin/gofmt /usr/local/bin/gofmt
RUN rm -rf /${GO_INSTALL_FILE}

VOLUME ["/root/go/src"]

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go install github.com/webx-top/tower@latest

WORKDIR /root/go/src/github.com/nging-plugins/caddymanager
ENTRYPOINT ["/root/go/bin/tower"]

# ENTRYPOINT ["/usr/local/bin/go"]
# CMD ["env"]

# docker build . -t "test-ubuntu-nginx"
