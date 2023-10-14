FROM crazymax/xgo:1.19.2
RUN curl https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz -o go.linux-amd64.tar.gz
RUN rm -rf /usr/local/go || true
RUN tar -zxvf go.linux-amd64.tar.gz -C /usr/local
RUN rm -rf ./go.linux-amd64.tar.gz

# docker build -t localhost/crazymax/xgo:1.21.3 .
