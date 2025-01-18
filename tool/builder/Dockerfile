FROM crazymax/xgo:1.23.0

RUN curl https://dl.google.com/go/go1.23.5.linux-amd64.tar.gz -o go.linux-amd64.tar.gz
RUN rm -rf /usr/local/go || true
RUN tar -zxvf go.linux-amd64.tar.gz -C /usr/local
RUN rm -rf ./go.linux-amd64.tar.gz
 
# docker build -t admpub/xgo:1.23.5 .
