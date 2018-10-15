# note: call scripts from /scripts

.PHONY: default build builder-image binary-image test stop clean-images clean push apply deploy

BUILDER ?= proxyinjector-builder
BINARY ?= ProxyInjector
DOCKER_IMAGE ?= stakater/proxyinjector
# Default value "dev"
DOCKER_TAG ?= 0.0.4
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

VERSION=$(shell cat .version)
BUILD=

GOCMD = go
GLIDECMD = glide
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build test

install:
	"$(GLIDECMD)" install

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

builder-image:
	@docker build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -

test:
	"$(GOCMD)" test -timeout 1800s -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i

push: ## push the latest Docker image to DockerHub
	docker push $(REPOSITORY)

apply:
	kubectl apply -f deployments/manifests/ -n temp-proxyinjector

deploy: binary-image push apply
