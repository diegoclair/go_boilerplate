.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mock
mock:
	go install github.com/golang/mock/mockgen@latest
	rm -rf mocks
	mockgen -package mocks -source=domain/contract/repo.go -destination=mocks/repo.go
	mockgen -package mocks -source=domain/contract/cache.go -destination=mocks/cache.go
	mockgen -package mocks -source=domain/contract/service.go -destination=mocks/service.go
