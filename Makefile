
clean:
	[ -f marvel-linux ] && rm marvel-linux || true
	[ -f marvel-osx ] && rm marvel-osx || true
	[ -f marvel-worker-linux ] && rm marvel-worker-linux || true
	[ -f marvel-worker-osx ] && rm marvel-worker-osx || true

marvel-linux: main.go
	GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "static"' -o $@

marvel-osx: main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@

marvel-worker-linux: ./worker/main.go
	GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "static"' -o $@

marvel-worker-osx: ./worker/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@

deps:
	go mod download

test:
	go test ./... -v

characters:
	go run ./cmd/main.go characters

migration:
	go run ./cmd/main.go migration

delete:
	go run ./cmd/main.go delete

mock:
	mockgen -source=client/client.go -destination=client/mock/client_mock.go -package=mock
	mockgen -source=business/usecase/usecase.go -destination=business/usecase/mock/usecase_mock.go -package=mock
	mockgen -source=business/data/character/store.go -destination=business/data/character/mock/store_mock.go -package=mock
	mockgen -source=business/data/character/cache.go -destination=business/data/character/mock/cache_mock.go -package=mock