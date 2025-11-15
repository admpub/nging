FROM --platform=$TARGETPLATFORM alpine
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION
ENV VERSION=${VERSION:-5.2.7}
#RUN apk update && apk upgrade

# RUN wget -c https://dl.webx.top/nging/v4.1.5/nging_linux_${TARGETARCH}.tar.gz -O /home/nging.tar.gz
# 对应 TARGETARCH 值通常为: amd64, arm64, arm, armv7 等（请确保构建产物与 TARGETARCH 一致）
COPY ./dist/packed/v${VERSION}/nging_linux_${TARGETARCH}.tar.gz /home/nging.tar.gz

RUN mkdir -p /home/nging_linux_${TARGETARCH} \
    && tar -zxvf /home/nging.tar.gz -C /home/nging_linux_${TARGETARCH} \
    && rm -f /home/nging.tar.gz

WORKDIR /home/nging_linux_${TARGETARCH}

# VOLUME [ "/home/nging_linux_amd64/data/cache", "/home/nging_linux_amd64/data/ftpdir", "/home/nging_linux_amd64/data/logs", "/home/nging_linux_amd64/data/sm2", "/home/nging_linux_amd64/myconfig", "/home/nging_linux_amd64/public" ]

ENTRYPOINT [ "./nging" ]
CMD [ "-p", "9999", "-c", "myconfig/config.yaml" ]

# * build *
# ./build-by-xgo.sh linux_amd64 min
# docker build . -t "admpub/nging:latest"
# * test * 
# docker run --rm -it -p "7770:9999" admpub/nging:latest