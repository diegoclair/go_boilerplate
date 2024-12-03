.PHONY: tests
tests:
	go test -v -cover ./...

.PHONY: mocks
mocks: install-mockgen domain-mocks infra-mocks

.PHONY: install-mockgen
install-mockgen:
	@echo "=====> Installing mockgen"
	@go install go.uber.org/mock/mockgen@latest

.PHONY: domain-mocks
domain-mocks:
	@echo "=====> Generating domain mocks"

	@rm -rf mocks
	@for file in domain/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=mocks/$$filename; \
	done

	@mockgen -package mocks -source=domain/infrastructure.go -destination=mocks/infrastructure.go

.PHONY: infra-mocks
infra-mocks:
	@echo "=====> Generating infra mocks"

	@rm -rf infra/mocks
	@for file in infra/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=infra/mocks/$$filename; \
	done
	
	@mockgen -package mocks -source=infra/cache/redis.go -destination=infra/mocks/redis.go
	
	@echo "=====> Mocks generated"

# @ to avoid echoing the command
.PHONY: docs
docs:
	@echo "=====> Generating docs"

	@go install github.com/swaggo/swag/cmd/swag@latest
	@cd goswag && \
	go run main.go && \
	cd .. && \
	swag init --pdl=2 --parseInternal -g ./goswag/main.go -o ./docs && \
	swag fmt -d ./goswag/

	@echo "=====> Docs generated"

