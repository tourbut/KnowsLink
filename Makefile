# Reproducible development checks; every failure stops its target.
.PHONY: install lint lint-go lint-adapters lint-config test build verify verify-runtime

install:
	go mod download
	npm ci --prefix adapters

lint: lint-go lint-adapters lint-config

lint-go:
	@set -eu; files="$$( $$(go env GOROOT)/bin/gofmt -l cmd internal)"; \
	if [ -n "$$files" ]; then printf 'Go formatting violations:\n%s\n' "$$files"; exit 1; fi
	go vet ./cmd/... ./internal/...
	go mod verify

lint-adapters:
	npm run check --prefix adapters

lint-config:
	npm exec --prefix adapters -- prettier --check compose.yaml sqlc.yaml
	python3 scripts/check_compose.py
	python3 -m py_compile scripts/check_compose.py scripts/verify_setup.py scripts/verify_runtime.py

test:
	go test -race ./cmd/... ./internal/...

build:
	mkdir -p build
	go build -trimpath -o build/relay ./cmd/relay
	go build -trimpath -o build/migrate ./cmd/migrate
	npm run build --prefix adapters

verify:
	python3 scripts/verify_setup.py

verify-runtime: build
	python3 scripts/verify_runtime.py
