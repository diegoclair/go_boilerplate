.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mock
mock:
	rm -rf mock
	mockgen -package mock -destination mock/repo.go github.com/diegoclair/go_boilerplate/domain/contract DataManager,Transaction,AccountRepo,AuthRepo
	mockgen -package mock -destination mock/cache.go github.com/diegoclair/go_boilerplate/domain/contract CacheManager
	mockgen -package mock -destination mock/service.go github.com/diegoclair/go_boilerplate/domain/service AccountService,AuthService,TransferService
