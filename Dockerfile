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

FROM alpine
WORKDIR /usr/app/
COPY --from=0 /usr/app/marvel-linux /usr/app/
COPY --from=0 /usr/app/config.yaml /usr/app/config.yaml
ENV MARVEL_CONFIG=./config.yaml
EXPOSE 8080
ENTRYPOINT ["/usr/app/marvel-linux"]
