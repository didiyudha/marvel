FROM golang

WORKDIR /usr/app/
COPY Makefile ./
COPY go.mod ./
RUN make deps
COPY ./ ./
RUN make marvel-worker-linux

FROM alpine

WORKDIR /app
COPY --from=0 /usr/app/marvel-worker-linux ./
ENTRYPOINT ["/app/marvel-linux"]