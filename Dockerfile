FROM golang:1.13.7-alpine AS builder

COPY . /src/

WORKDIR /src

RUN go build -mod vendor -o /opa-notary-connector


FROM alpine:3.11

COPY --from=builder /opa-notary-connector /

RUN mkdir /etc/opa-notary-connector && \
    chgrp -R 0 /etc/opa-notary-connector && \
    chmod -R g=u /etc/opa-notary-connector

USER 1001
CMD ["/opa-notary-connector", "--help"]
