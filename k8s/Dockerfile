FROM golang

WORKDIR /usr/app/
COPY Makefile ./
COPY go.mod ./
RUN make deps
COPY ./ ./
RUN make marvel-linux

FROM alpine

WORKDIR /app
COPY --from=0 /usr/app/marvel-linux ./
EXPOSE 8080
ENTRYPOINT ["/app/marvel-linux"]