# Short Drama Recommendation - go-zero generation/build configuration.

SERVICE=ShortDrama
SERVICE_STYLE=short_drama
SERVICE_LOWER=short_drama
SERVICE_SNAKE=short_drama
SERVICE_DASH=short-drama

PROJECT_STYLE=go_zero
PROJECT_I18N=true

ENT_FEATURE=sql/execquery,intercept,sql/modifier

.PHONY: infra-up infra-down infra-logs gen gen-api gen-rpc gen-ent tidy fmt test

infra-up:
	docker compose up -d

infra-down:
	docker compose down

infra-logs:
	docker compose logs -f

gen: gen-api gen-rpc gen-ent
	@echo "go-zero API/RPC and Ent generation completed"

gen-api:
	goctl api go -api api/auth.api -dir ./api/auth-api
	goctl api go -api api/drama.api -dir ./api/drama-api
	goctl api go -api api/drama-admin.api -dir ./api/drama-admin-api
	@echo "go-zero API generation completed"

ifeq ($(OS),Windows_NT)
gen-rpc:
	powershell -NoProfile -ExecutionPolicy Bypass -File scripts/gen-rpc.ps1
else
gen-rpc:
	bash scripts/gen.sh
endif
	@echo "go-zero RPC generation completed"

gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./rpc/ent/template/*.tmpl" ./rpc/ent/schema --feature $(ENT_FEATURE)
	@echo "Ent ORM generation completed"

tidy:
	go mod tidy

fmt:
	gofmt -w $$(find . -name '*.go' -type f -not -path './vendor/*')

test:
	go test ./...
