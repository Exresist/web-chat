export GOFLAGS=-mod=vendor

OUT := bin/web-chat

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

clean:
	rm -rf ./bin/*
.PHONY: clean

build: clean
	go build -ldflags "-w -s" -o $(OUT) ./cmd/web-chat
.PHONY: build

dbuild: vendor
	docker-compose -f ../web-chat/docker-compose.yml build
.PHONY: dbuild

dstart:
	docker-compose -f ../web-chat/docker-compose.yml up
.PHONY: dstart

run: dbuild dstart
.PHONY: run

gen:
	go generate ./...
.PHONY: gen

test:
	go test $(TEST_CONFIG) $(TEST_PACKAGE)
.PHONY: test

prepare: clean install.tools ## performs steps needed before first build
.PHONY: prepare

install.tools: export GOBIN:=$(OUT_PATH)
install.tools: clean.bin ## install all tools mentioned in tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
.PHONY: install.tools