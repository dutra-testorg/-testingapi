BUILD_TAG = \$(shell git describe --tags \$(git describe --tags --abbrev=0 --always) --always)
BUILD_TIME = \$(shell date +%s)
GOCMD=go
GOVERSION = \$(shell go version | cut -c 14-)
GOBUILD=\$(GOCMD) build
GORUN=\$(GOCMD) run
GOGENERATE=\$(GOCMD) generate
GOBUILD_BASE_PARAMS = -mod vendor -ldflags="-X 'main.buildVersion=\$(BUILD_TAG)' -X 'main.buildTime=\$(BUILD_TIME)'-X 'main.goVersion=\$(GOVERSION)'"
GOBUILD_VARS = CGO_ENABLED=0 GOOS=linux
GOBUILD_PARAMS = \$(GOBUILD_BASE_PARAMS)
GOBUILD_DEBUG_FLAGS = -gcflags="all=-N -l"
GOCLEAN=\$(GOCMD) clean
GOTEST=\$(GOCMD) test
GOTEST_BASE_PARAMS = -p 1 -count=1
GOTEST_BASE_PARAMS_V = -v \$(GOTEST_BASE_PARAMS)
GOTEST_PARAMS = \$(GOTEST_BASE_PARAMS)
GOTEST_PARAMS_V = \$(GOTEST_BASE_PARAMS_V)
PKG_LST = ./internal/... ./pkg/...
COVER_EXCLUDE_PKG_LST = mocks|testutils|testdata
COVER_PKG_LST = \$(shell go list \$(PKG_LST) | grep -Ev "\$(COVER_EXCLUDE_PKG_LST)" | tr '\n' ',' | sed 's/[,]\$\$//')
GOMOD=\$(GOCMD) mod
SWAGCMD=swag
DEBUG_FLAGS=-gcflags="all=-N -l"
SWAG_PARAMS = init --parseInternal --parseDependency --parseVendor --parseDepth 3
SWAG_EXCLUDE = --exclude ./test,./docs,./helm,./scripts,./db,./docker,./bin,./swagger,./tools,./api,vendor/github.com/DataDog
GOLINT_CMD=golangci-lint
GOLINT_BASE_RUN=\$(GOLINT_CMD) run --modules-download-mode vendor --timeout=8m
GOLINT_RUN=\$(GOLINT_BASE_RUN)
GOLINT_EXT_MODULES=-E gosec -E prealloc -E exportloopref -E misspell -E nestif -E unconvert  -E depguard -E goimports -E prealloc -E sqlclosecheck -E bodyclose -E noctx -E govet -E gosimple -E gofmt -E unparam
BINARY_DIR=bin
BINARY_NAME=app
CURRENT_DIR = \$(shell pwd)
DOCKER_IMAGE=$name;format="lower,hyphen"$
OS=linux
ARCH=amd64
UNAME_S = \$(shell uname -s)
UNAME_M = \$(shell uname -m)
MIGRATION=migration
PKG_CONFIG_PATH	:=

ifeq (\$(UNAME_S),Darwin)
	OS=darwin
	PKG_CFG_PATH = PKG_CONFIG_PATH="\$(shell find /opt/homebrew/Cellar/openssl@1.1 -type d -d 1 )/lib/pkgconfig/"
	GOLINT_RUN = \$(GOLINT_BASE_RUN) --build-tags dynamic
	GOBUILD_PARAMS = \$(GOBUILD_BASE_PARAMS) -tags dynamic
	GOTEST_PARAMS = \$(GOTEST_BASE_PARAMS) -tags dynamic
	GOTEST_PARAMS_V = \$(GOTEST_BASE_PARAMS_V) -tags dynamic
endif

ifeq (\$(UNAME_M),arm64)
	ARCH=arm64
endif

.PHONY: all
all: deps lint test build

.PHONY: build
build: bindir
	\$(PKG_CFG_PATH) \$(GOBUILD_VARS) \$(GOBUILD) \$(GOBUILD_PARAMS) -o \$(BINARY_DIR)/\$(BINARY_NAME) ./cmd/app

.PHONY: debug
debug: bindir
	\$(PKG_CFG_PATH) \$(GOBUILD) \$(GOBUILD_PARAMS) \$(GOBUILD_DEBUG_FLAGS) -o \$(BINARY_DIR)/\$(BINARY_NAME) ./cmd/app

.PHONY: tidy
tidy:
	\$(GOMOD) tidy

.PHONY: vendor
vendor:
	\$(GOMOD) vendor

.PHONY: generate
generate:
	\$(PKG_CFG_PATH) \$(GOGENERATE) \$(PKG_LST)

.PHONY: deps
deps: tidy vendor generate

.PHONY: bindir
bindir:
	mkdir -p \$(BINARY_DIR)

.PHONY: test
test: unit-tests integration-tests

.PHONY: unit-tests
unit-tests:
	\$(PKG_CFG_PATH) \$(GOTEST) \$(GOTEST_PARAMS) \$(PKG_LST)

.PHONY: integration-tests
integration-tests: test-infra-up test-db-migrate
	\$(PKG_CFG_PATH) \$(GOTEST) \$(GOTEST_PARAMS) ./test/integration/... || (make test-infra-down && exit 1)
	make test-infra-down

.PHONY: pipeline-tests
pipeline-tests: test-db-migrate
	\$(PKG_CFG_PATH) \$(GOTEST) \$(GOTEST_PARAMS) -timeout=15m -coverprofile=coverage_all.txt -covermode count -coverpkg="\$(COVER_PKG_LST)" ./...

.PHONY: clean
clean:
	rm -rf \$(BINARY_DIR)
	\$(GOCLEAN) -cache -modcache # optional

test-cov:
	\$(PKG_CFG_PATH) \$(GOTEST) \$(GOTEST_PARAMS) -timeout=15m -coverprofile=coverage_all.txt -covermode count -coverpkg="\$(COVER_PKG_LST)" ./...
	gocov convert coverage_all.txt | gocov report
	gocov convert coverage_all.txt | gocov-html > coverage.html

.PHONY: run
run:
	\$(PKG_CFG_PATH) go run \$(GOBUILD_PARAMS) cmd/app/main.go

docker-build:
	docker build -t \$(DOCKER_IMAGE):test .
	docker run --rm -p 8080:8080 --net host --env-file ./docker/env.list \$(DOCKER_IMAGE):test

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir db/migrations/ \$(MIGRATION)
	ls db/migrations/*.up.sql -r1 | head -n 1 > db/last_migration

.PHONY: swagger
swagger:
	\$(SWAGCMD) \$(SWAG_PARAMS) \$(SWAG_EXCLUDE) -o ./api -g ./cmd/app/main.go

coverage-report: test-infra-up pipeline-coverage test-infra-down
	gocov convert coverage.txt | gocov report
	gocov convert coverage.txt | gocov-html > coverage.html

pipeline-coverage: pipeline-tests
	grep -Ev "enumer.go|fielder.go|factory.go" coverage_all.txt > coverage.txt
	make -f tools/Makefile install-gocov
	make -f tools/Makefile install-gocovxml
	gocov convert coverage.txt | gocov report
	gocov convert coverage.txt | gocov-xml > coverage.xml
 
test-infra-up:
	(docker-compose -f ./docker/docker-compose.test.yaml up --abort-on-container-exit > /dev/null &) && until docker exec -it postgres psql -U postgres -c "select 1" -d testdb > /dev/null 2>&1; do echo waiting db; done

test-infra-down:
	docker-compose -f docker/docker-compose.test.yaml down --volumes > /dev/null

test-db-migrate:
	make -f tools/Makefile install-migrate
	migrate -database 'postgres://postgres:docker@localhost:5432/testdb?connect_timeout=10&sslmode=disable' -path \$(CURRENT_DIR)/db/migrations up

test-db-migrate-down:
	docker run -it --net=host -v \$(CURRENT_DIR)/db/migrations:/migrations migrate/migrate -path=/migrations -database postgres://postgres:docker@localhost:5432/testdb?sslmode=disable down

.PHONY: lint
lint:
	make -f tools/Makefile install-golangci
	\$(GOLINT_RUN) \$(GOLINT_EXT_MODULES)

.PHONY: lint-fix
lint-fix:
	\$(GOLINT_RUN) --fix \$(GOLINT_EXT_MODULES)
