.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mock
mock:
# examples before:
# mockgen -package mocks -source=domain/contract/repo.go -destination=mocks/repo.go
# mockgen -package mocks -destination mocks/repo.go github.com/diegoclair/go_boilerplate/domain/contract DataManager,AccountRepo,AuthRepo
	go install github.com/golang/mock/mockgen@latest
	rm -rf mocks
	for file in domain/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=mocks/$$filename; \
	done
