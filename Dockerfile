FROM golang:1.19-alpine3.17

ENV GOLANG_VERSION 1.19.6
ENV GOOSE_VERSION 3.6.1

RUN apk add --no-cache bash \
    && apk add --no-cache tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
    && echo "Europe/Moscow" >  /etc/timezone \
    # see more https://github.com/golang/go/issues/22846
    && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# install goose
RUN wget -q -O goose "https://github.com/pressly/goose/releases/download/v${GOOSE_VERSION}/goose_linux_x86_64" \
    && mv goose ${GOPATH}/bin/goose && chmod +x ${GOPATH}/bin/goose


RUN mkdir /app

COPY config.yaml /app/
COPY migrations /app/migrations
COPY binfile /app/psn_discounter
COPY entrypoint.sh /app/

WORKDIR /app

ENTRYPOINT ["bash","entrypoint.sh"]