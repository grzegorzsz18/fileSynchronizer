FROM golang:1.11

COPY server ./src/fileSender/server
COPY pkg ./src/fileSender/pkg

RUN export GOROOT=/go/src/bin
RUN export GOPATH=/go/src/bin

RUN go get fileSender/pkg
RUN go get fileSender/server

RUN go get github.com/mattn/go-sqlite3
RUN go install github.com/mattn/go-sqlite3

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ./src/fileSender/server/main.go



FROM alpine:latest

RUN apk --update upgrade
RUN apk add sqlite
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN rm -rf /var/cache/apk/*


WORKDIR /root/

COPY --from=0 /go/main .
COPY --from=0 /go/src/fileSender/server/users.db .

EXPOSE 18080
EXPOSE 22222

CMD ["./main"]