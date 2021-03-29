FROM golang:1.16.2

# envs
ENV GOPROXY https://goproxy.cn,direct
ENV GIN_MODE release

# copy & build
WORKDIR $GOPATH/src/github.com/txfs19260817/scopelens/server
COPY . $GOPATH/src/github.com/txfs19260817/scopelens/server
RUN mv ./config/config.ini ./config/config.ini.dev && mv ./config/config.ini.prod ./config/config.ini
RUN go build -v .

# HTTPs cert&key (-v $CERT_PATH:$CERT_PATH)
ARG CERT_PATH
RUN mkdir -p $CERT_PATH

# set port (-p $PORT:$PORT)
ARG PORT
EXPOSE $PORT

# run
ENTRYPOINT ["./server"]