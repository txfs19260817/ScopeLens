# Build stage
FROM golang:1.16-alpine3.13 AS builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go build -v .

# Run stage (set port `-p $PORT:$PORT`)
FROM alpine:3.13
WORKDIR $GOPATH/src/github.com/txfs19260817/scopelens/server
COPY --from=builder /app/server .
RUN mkdir config log assets
COPY assets ./assets
ARG PORT=8888
EXPOSE $PORT

# run
ENTRYPOINT ["./server"]