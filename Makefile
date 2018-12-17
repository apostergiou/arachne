.PHONY: list build arachned start test install deps lint

LISTEN?="0.0.0.0:53"
PUBLISH?="53"
UPSTREAM?="1.1.1.1:853"
NETWORK?="tcp"
HOST_PORT?=5300

list:
	@echo "Available commands:"
	@echo "  build    - build the image"
	@echo "  arachned - start the daemon"
	@echo "  start    - start the arachne server locally"
	@echo "  test     - run the tests in the container"
	@echo "  install  - install the library as binary"
	@echo "  deps     - check dependencies"
	@echo "  lint     - run the golint tool"

build:
	docker build -t arachne .

arachned: build
	docker run -d -e LISTEN="$(LISTEN)" -e UPSTREAM="$(UPSTREAM)" -e NETWORK=$(NETWORK) -p $(HOST_PORT):$(PUBLISH) arachne

start: build
	docker run -e LISTEN="$(LISTEN)" -e UPSTREAM="$(UPSTREAM)" -e NETWORK=$(NETWORK) -p $(HOST_PORT):$(PUBLISH) arachne

test: build
	docker run arachne \
		go test ./... -v

# Commands for development
install: fmt test
	go install -v ./...

deps:
	dep ensure -v

lint:
	golint `go list ./... | grep -v /vendor/`
