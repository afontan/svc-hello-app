SERVICE := svc-hello
OWNER   ?= afontan
IMAGE   := ghcr.io/$(OWNER)/$(SERVICE)
VERSION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

.PHONY: tidy test build run docker push clean

tidy:
	go mod tidy

test:
	go test ./... -v

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(SERVICE) .

run: build
	DD_SERVICE=$(SERVICE) DD_ENV=dev DD_VERSION=$(VERSION) ./bin/$(SERVICE)

docker:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

push: docker
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

clean:
	rm -rf bin/
