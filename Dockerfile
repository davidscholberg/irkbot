FROM golang:latest
ENV HOME=/root
RUN go get github.com/dvdmuckle/irkbot
ENTRYPOINT ["/go/bin/irkbot"]
