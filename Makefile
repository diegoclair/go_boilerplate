.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mock
mock:
	go install github.com/golang/mock/mockgen@latest
	rm -rf mocks
	for file in domain/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=mocks/$$filename; \
	done