FROM golang:1.11

COPY server ./src/fileSender/server
COPY pkg ./src/fileSender/pkg

RUN export GOROOT=/go/src/bin
RUN export GOPATH=/go/src/bin

RUN go get fileSender/pkg
RUN go get fileSender/server

RUN CGO_ENABLED=0 GOOS=linux go build ./src/fileSender/server/main.go



FROM alpine:latest

WORKDIR /root/

COPY --from=0 /go/main .

EXPOSE 18080

CMD ["./main"]