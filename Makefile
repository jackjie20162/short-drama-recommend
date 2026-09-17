.PHONY: infra-up infra-down infra-logs gen fmt test

infra-up:
	docker compose up -d

infra-down:
	docker compose down

infra-logs:
	docker compose logs -f

gen:
	@echo "Install goctl and protoc locally, then run scripts/gen.sh"
	bash scripts/gen.sh

fmt:
	gofmt -w $$(find . -name '*.go' -type f -not -path './vendor/*')

test:
	go test ./...
