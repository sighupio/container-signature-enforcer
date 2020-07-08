FROM alpine:3.11

COPY opa-notary-connector /
RUN mkdir /etc/opa-notary-connector && \
    chgrp -R 0 /etc/opa-notary-connector && \
    chmod -R g=u /etc/opa-notary-connector && \
    chmod +x /opa-notary-connector && \
    /opa-notary-connector --help

USER 1001
ENTRYPOINT [ "/opa-notary-connector" ]
