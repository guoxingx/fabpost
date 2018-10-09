.PHONY: network build clean tool lint help

NETWORK=$(n)
ifeq "$(NETWORK)" ""
	NETWORK = default
endif

all: build

net:
	@cd network/$(NETWORK); docker-compose up -d

netdown:
	@cd network/$(NETWORK); docker-compose down

build:
	@go build -v.

clean:
	rm -rf fabtreehole
	go clean -i .

tool:
	go tool vet . |& grep -v vendor; true
	gofmt -w .

lint:
	golint ./...

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
