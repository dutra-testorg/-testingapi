# $name;format="lower,hyphen"$

## Objectives

This project is a very basic template for any REST API built with Go providing standards for logging, testing, API documentation (swagger), Datadog tracer APM integration, and based on the following layout guide [https://github.com/golang-standards/project-layout].

_This is not a SDK and the primary objective is to provide guidance and best practices._

## Go environment setup

- **download and install go from _https://golang.org/doc/install_**

- **add to _.bashrc_**:

        export GOROOT=/opt/go
        export PATH=\$PATH:\$HOME/go/bin:\$GOROOT/bin

## Setup dependencies

    go clean -cache -modcache # optional

    git config --global url."git@github.com:".insteadOf "https://github.com/"
    go mod tidy -v
    go mod vendor -v

##### for private repos (optional if https over ssh):

    go env -w GOPRIVATE="github.com/Gympass/*,github.com/gympass/*"
    export GIT_TERMINAL_PROMPT=1

- more info: https://medium.com/@tim_raymond/fetching-private-dependencies-with-go-modules-1d65afe47c62

##### hint:

- clone private repo (e.g.: https://github.com/Gympass/gcore)
- add to your go.mod:

        module $name;format="lower,hyphen"$

        go 1.19

        require (
                github.com/Gympass/gcore v1.0.0
                github.com/kelseyhightower/envconfig v1.4.0
                github.com/sirupsen/logrus v1.6.0
                github.com/swaggo/swag v1.6.7
                ...
                gopkg.in/DataDog/dd-trace-go.v1 v1.23.2
                gopkg.in/yaml.v2 v2.3.0
        )

        replace github.com/Gympass/gcore => <local directory>/gympass/gcore

  - github.com/Gympass/gcore was cloned at _\$HOME/gympass/gcore_
      - this is a good choice when you need to **test/patch** both codes

## Run locally

#### Local environment

        go run cmd/app/main.go

#### Docker

        docker build -t $name;format="lower,hyphen"$:test --build-arg SSH_PRIVATE_KEY="\$(cat \$HOME/.ssh/id_rsa)" .
        docker run --rm -p 8080:8080 --net host --env-file ./docker/env.list $name;format="lower,hyphen"$:test

#### Docker compose

        docker-compose -f docker/docker-compose.dev.yaml up

#### Load test

        k6 run ./scripts/k6/load.js

#### Unit test

        go test -count=10 -race -cover ./...

#### Coverage

        go test -v -count=1 -coverprofile /tmp/cover.out -cover  ./...
        go tool cover -html=/tmp/cover.out

#### Profile

        go run cmd/app/main.go
        http://localhost:8080/debug/pprof/
        # k6 run ./scripts/k6/load.js

        # Online
        go tool pprof http://localhost:8080/debug/pprof/heap
        go tool pprof -png http://localhost:8080/debug/pprof/heap > /tmp/out.png

        # Offline
        curl -sK -v http://localhost:8080/debug/pprof/heap > /tmp/heap.out
        go tool pprof /tmp/heap.out
        curl -sK -v http://localhost:8080/debug/pprof/profile?seconds=10 > /tmp/profile.out
        go tool pprof /tmp/profile.out

## Documentation

#### Swagger

- generate API docs using:

        swag init -o ./api -g ./internal/rest/rest.go ./internal/micro/handler.go
