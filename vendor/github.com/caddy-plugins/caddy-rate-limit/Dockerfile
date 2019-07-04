FROM golang:1.11-alpine

RUN apk add --update-cache --no-cache git && \
    go get -v github.com/caddyserver/caddy && \
    go get -v github.com/caddyserver/builds && \
    go get -v golang.org/x/time/rate

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

COPY . .

RUN ./insert-plugin.sh && \
    cd /go/src/github.com/caddyserver/caddy/caddy && \
    go run build.go && \
    cp caddy /go/src/github.com/xuqingfeng/caddy-rate-limit/caddy

EXPOSE 2016

CMD ["./caddy"]