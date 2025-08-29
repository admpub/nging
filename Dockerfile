FROM alpine
ARG VERSION
ENV VERSION=${VERSION:-5.2.6}
#RUN apk update && apk upgrade

# RUN wget -c https://dl.webx.top/nging/v4.1.5/nging_linux_amd64.tar.gz -O /home/nging_linux_amd64.tar.gz
COPY ./dist/packed/v${VERSION}/nging_linux_amd64.tar.gz /home/nging_linux_amd64.tar.gz
RUN mkdir /home/nging_linux_amd64 && tar -zxvf /home/nging_linux_amd64.tar.gz -C /home/nging_linux_amd64 && rm -rf /home/nging_linux_amd64.tar.gz

WORKDIR /home/nging_linux_amd64

# VOLUME [ "/home/nging_linux_amd64/data/cache", "/home/nging_linux_amd64/data/ftpdir", "/home/nging_linux_amd64/data/logs", "/home/nging_linux_amd64/data/sm2", "/home/nging_linux_amd64/myconfig", "/home/nging_linux_amd64/public" ]

ENTRYPOINT [ "./nging" ]
CMD [ "-p", "9999", "-c", "myconfig/config.yaml" ]

# * build *
# ./build-by-xgo.sh linux_amd64 min
# docker build . -t "admpub/nging-dockermgr:latest" --build-arg VERSION=$(grep NgingVersion ./tool/nging-builder/builder.conf | sed 's/NgingVersion[ ]*:[ ]*//g' | sed 's/"//g' | sed 's/ //g')
# * test * 
# docker run --rm -it -p "7770:9999" admpub/nging-dockermgr:latest