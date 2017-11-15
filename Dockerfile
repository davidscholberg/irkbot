FROM golang:latest
ENV HOME=/root
RUN go get github.com/jholtom/irkbot
ENTRYPOINT ["/go/bin/irkbot"]
