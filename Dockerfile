FROM golang:1.22.3-alpine3.19

ARG COMMAND

ENV binary $COMMAND

WORKDIR /go/bin

COPY ${COMMAND} /go/bin

CMD ["sh", "-c", "${binary}"]