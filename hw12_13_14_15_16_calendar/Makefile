BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
DOCKER_IMG="calendar:develop"

PROTO_DIR="proto"
PB_DIR="pb"

GIT_HASH := $(shell git log --format="%h" -n 1)
BUILD_DATE := ""
LDFLAGS := -X main.release="develop" -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)
PROTOC_INCLUDES := "$(shell go env GOPATH)\pkg\mod\github.com\protocolbuffers\protobuf@v3\src"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

run: build
	$(BIN) -config ./configs/config.yaml

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... #./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

docker-run:
	cd deployments && docker-compose up -d

migrate:
	goose up -dir migrations

generate-proto:
	protoc \
		--proto_path=. \
		--proto_path=$(PROTOC_INCLUDES) \
		--go_out=internal/server/grpc/pb \
    	--go-grpc_out=internal/server/grpc/pb \
    	api/*.proto