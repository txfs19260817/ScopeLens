FROM golang:1.16.2

# envs
ENV GOPROXY https://goproxy.cn,direct
ENV GIN_MODE release

# copy & build
WORKDIR $GOPATH/src/github.com/txfs19260817/scopelens/server
COPY . $GOPATH/src/github.com/txfs19260817/scopelens/server
RUN go build -v .

# set port
ARG PORT
EXPOSE $PORT

# start
ENTRYPOINT ["./server"]