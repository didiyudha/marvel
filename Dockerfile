FROM golang

ARG APP_CONF
WORKDIR /usr/app/
COPY Makefile ./
COPY go.mod ./
COPY go.sum ./
RUN make deps
COPY ./ ./
RUN make marvel-linux
RUN echo "$APP_CONF" > /usr/app/config.yaml
RUN cat /usr/app/config.yaml

FROM alpine
RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=0 /usr/app/marvel-linux ./
COPY --from=0 /usr/app/config.yaml ./
RUN ls -al
ENV MARVEL_CONFIG=./config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/marvel-linux"]
