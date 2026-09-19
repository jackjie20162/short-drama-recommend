# Simple Admin compatible build/generation configuration.

SERVICE=ShortDrama
SERVICE_STYLE=short_drama
SERVICE_LOWER=short_drama
SERVICE_SNAKE=short_drama
SERVICE_DASH=short-drama

PROJECT_STYLE=go_zero
PROJECT_I18N=true

# Keep this aligned with the Simple Admin Ent workflow.
ENT_FEATURE=sql/execquery,intercept,sql/modifier

.PHONY: infra-up infra-down infra-logs gen gen-ent fmt test

infra-up:
	docker compose up -d

infra-down:
	docker compose down

infra-logs:
	docker compose logs -f

gen:
	@echo "Install goctl/goctls and protoc locally, then run scripts/gen.sh"

gen-ent:
	goctls run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./rpc/ent/template/*.tmpl" ./rpc/ent/schema --feature $(ENT_FEATURE)
	@echo "Generate Ent files successfully"

fmt:
	gofmt -w $$(find . -name '*.go' -type f -not -path './vendor/*')

test:
	go test ./...
