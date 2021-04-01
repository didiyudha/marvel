mock:
	mockgen -source=client/client.go -destination=client/mock/client_mock.go -package=mock
	mockgen -source=business/usecase/usecase.go -destination=business/usecase/mock/usecase_mock.go -package=mock
	mockgen -source=business/data/character/store.go -destination=business/data/character/mock/store_mock.go -package=mock
	mockgen -source=business/data/character/cache.go -destination=business/data/character/mock/cache_mock.go -package=mock
test:
	go test ./... -v
characters:
	go run ./cmd/main.go characters
migration:
	go run ./cmd/main.go migration
delete:
	go run ./cmd/main.go delete