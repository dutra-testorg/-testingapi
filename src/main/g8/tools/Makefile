GOCMD=go
CDTMP=cd /tmp
GOINSTALL=\$(GOCMD) install
GOPATHBIN=\$(GOPATH)/bin
GOENVPATH=\$(shell go env GOPATH)
CURL=curl
OS=linux
ARCH=amd64
UNAME_S = \$(shell uname -s)

ifeq (\$(GOENVPATH),)
	GOENVPATH=\$(GOPATHBIN)
else
	GOENVPATH = \$(shell go env GOPATH)/bin
endif

ifeq (\$(UNAME_S),Darwin)
	OS=darwin
endif
UNAME_M = \$(shell uname -m)
ifeq (\$(UNAME_M),arm64)
	ARCH=arm64
endif
OSARCH=\$(OS)-\$(ARCH)

.PHONY: all
all: install

.PHONY: install
install: install-mockgen install-enumer install-swagger install-migrate install-golangci install-gocov install-gocovxml

.PHONY: install-mockgen
install-mockgen:
	\$(GOINSTALL) github.com/derision-test/go-mockgen/cmd/go-mockgen@latest

.PHONY: install-enumer
install-enumer:
	\$(GOINSTALL) github.com/dmarkham/enumer@latest

.PHONY: install-swagger
install-swagger:
	\$(GOINSTALL) github.com/swaggo/swag/cmd/swag@latest

.PHONY: install-migrate
install-migrate:
	\$(CURL) -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.\$(OSARCH).tar.gz | tar -C \$(GOENVPATH) -xvz migrate

.PHONY: install-gocov
install-gocov:
	\$(GOINSTALL) github.com/axw/gocov/gocov@latest

.PHONY: install-gocovxml
install-gocovxml:
	\$(GOINSTALL) github.com/AlekSi/gocov-xml@latest

.PHONY: install-golangci
install-golangci:
	\$(CURL) -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(GOENVPATH) v1.50.0
