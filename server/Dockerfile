FROM golang:1.16.2

# envs
ENV GOPROXY https://goproxy.cn,direct
ENV GIN_MODE release

# copy & build
WORKDIR $GOPATH/src/github.com/txfs19260817/scopelens/server
COPY . $GOPATH/src/github.com/txfs19260817/scopelens/server
RUN mv ./config/config.ini ./config/config.ini.dev && mv ./config/config.ini.prod ./config/config.ini
RUN go build -v .

# set path to cert files
RUN \
  mkdir -p /www/server/panel/vhost/cert/scopelens.team && \
  mv fullchain.pem /www/server/panel/vhost/cert/scopelens.team/fullchain.pem && \
  mv privkey.pem /www/server/panel/vhost/cert/scopelens.team/privkey.pem

# set port
ARG PORT
EXPOSE $PORT

# start
ENTRYPOINT ["./server"]