.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mock
mock:
	rm -rf mocks
	mockgen -package mocks -destination mocks/repo.go github.com/diegoclair/go_boilerplate/domain/contract DataManager,Transaction,AccountRepo,AuthRepo
	mockgen -package mocks -destination mocks/cache.go github.com/diegoclair/go_boilerplate/domain/contract CacheManager
	mockgen -package mocks -destination mocks/service.go github.com/diegoclair/go_boilerplate/domain/service AccountService,AuthService,TransferService
