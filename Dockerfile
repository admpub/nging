FROM docker.enote.vip/library/alpine
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION
ENV VERSION=${VERSION:-5.2.7}
#RUN apk update && apk upgrade

# RUN wget -c https://dl.webx.top/nging/v4.1.5/nging_linux_${TARGETARCH}.tar.gz -O /home/nging.tar.gz
# 对应 TARGETARCH 值通常为: amd64, arm64, arm, armv7 等（请确保构建产物与 TARGETARCH 一致）
COPY ./dist/packed/v${VERSION}/nging_linux_${TARGETARCH}.tar.gz /home/nging.tar.gz

# 创建 nging_linux_amd64 文件夹兼容旧版本
RUN mkdir -p /home/nging_linux_amd64 && ln -s /home/nging_linux_amd64 /home/nging \
    && tar -zxvf /home/nging.tar.gz -C /home/nging \
    && rm -f /home/nging.tar.gz

WORKDIR /home/nging

# VOLUME [ "/home/nging/data/cache", "/home/nging/data/ftpdir", "/home/nging/data/logs", "/home/nging/data/sm2", "/home/nging/myconfig", "/home/nging/public" ]

ENTRYPOINT [ "./nging" ]
CMD [ "-p", "9999", "-c", "myconfig/config.yaml" ]

# * build *
# cd ./tool/ngnig-builder && ./run.sh linux_amd64 && cd ../../
# docker build . -t "admpub/nging-dockermgr:latest" --build-arg VERSION=$(grep NgingVersion ./tool/nging-builder/builder.conf | sed 's/NgingVersion[ ]*:[ ]*//g' | sed 's/"//g' | sed 's/ //g')
# * test * 
# docker run --rm -it -p "7770:9999" admpub/nging-dockermgr:latest