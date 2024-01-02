FROM golang:1.21-alpine 

ARG COMMAND

ENV binary $COMMAND

WORKDIR /go/bin

COPY ${COMMAND} /go/bin

CMD ["sh", "-c", "${binary}"]