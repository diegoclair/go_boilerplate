.PHONY: tests
tests:
	go test -v -cover ./...

.PHONY: mocks
mocks:
	@echo "=====> Generating mocks"

	@go install go.uber.org/mock/mockgen@latest
	@rm -rf mocks
	@for file in domain/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=mocks/$$filename; \
	done

	@for file in infra/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=infra/mocks/$$filename; \
	done
	
	@mockgen -package mocks -source=infra/cache/redis.go -destination=mocks/redis.go
	
	@echo "=====> Mocks generated"

# @ to avoid echoing the command
.PHONY: docs
docs:
	@echo "=====> Generating docs"

	@go install github.com/swaggo/swag/cmd/swag@latest
	@cd goswag && \
	go run main.go && \
	cd .. && \
	swag init --pd -g ./goswag/main.go && \
	swag fmt -d ./goswag/

	@echo "=====> Docs generated"

