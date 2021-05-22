# Copyright IBM Corp All Rights Reserved.
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - all (default) - compiles protos, downloads vendor libs, and build all executables
#   - proto-compiler - downloads protobuf compiler. this is a pre-requisite in order to run the project
#   - protos - generate all protobuf artifacts based on .proto files
#   - clean-protos - removes the generated proto go files
#   - vendor - download all third party libraries and puts them inside vendor directory
#   - clean-vendor - removes third party libraries from vendor directory
#   - ping-vnf-manager - builds ping-vnfm as an executable and puts it under build/bin
#   - ping-vnf - builds ping-vnf as an executable and puts it under build/bin
#   - pong-vnf-manager - builds pong-vnfm as an executable and puts it under build/bin
#   - pong-vnf - builds pong-vnf as an executable and puts it under build/bin
#   - docker-build - builds docker image locally for running the components using docker
#   - docker-push - pushes the local docker image to docker registry. in order to run it need to to add REGISTRY and TAG in arguments.
#   - clean - cleans the build area (all executables under build/bin)
#   - clean-all - superset of 'clean' that also removes vendor dir and compiled proto go files

.PHONY: all				##compiles protos, downloads vendor libs, and build all executables
all: proto-compiler protos vendor ping-vnf-manager ping-vnf pong-vnf-manager pong-vnf

.PHONY: proto-compiler		##downloads protobuf compiler. this is a pre-requisite in order to run the project
proto-compiler:
	@build/scripts/install-proto-compiler.sh

.PHONY: protos			##generate all protobuf artifacts based on .proto files
protos: clean-protos
	@build/scripts/compile-protos.sh

.PHONY: clean-protos			##removes the generated proto go files
clean-protos:
	-@rm -rf ./protos/*.pb.go

.PHONY: vendor			##download all third party libraries and puts them inside vendor directory
vendor:
	@go mod vendor

.PHONY: clean-vendor			##removes third party libraries from vendor directory
clean-vendor:
	-@rm -rf vendor

.PHONY: ping-vnf-manager		##builds ping-vnf-manager as an executable and puts it under build/bin
ping-vnf-manager:
	@go build -o build/bin/ping-vnf-manager cmd/ping-vnf-manager/main.go

.PHONY: ping-vnf			##builds ping-vnf as an executable and puts it under build/bin
ping-vnf:
	@go build -o build/bin/ping-vnf cmd/ping-vnf/main.go

.PHONY: pong-vnf-manager		##builds pong-vnf-manager as an executable and puts it under build/bin
pong-vnf-manager:
	@go build -o build/bin/pong-vnf-manager cmd/pong-vnf-manager/main.go

.PHONY: pong-vnf			##builds pong-vnf as an executable and puts it under build/bin
pong-vnf:
	@go build -o build/bin/pong-vnf cmd/pong-vnf/main.go

.PHONY: docker-build			##builds docker images locally for running the components using docker
docker-build: all
	@docker build -t ping-pong-sample/ping-vnf-manager -f build/docker/Dockerfile.ping_vnf_manager .
	@docker build -t ping-pong-sample/ping-vnf -f build/docker/Dockerfile.ping_vnf .
	@docker build -t ping-pong-sample/pong-vnf-manager -f build/docker/Dockerfile.pong_vnf_manager .
	@docker build -t ping-pong-sample/pong-vnf -f build/docker/Dockerfile.pong_vnf .

.PHONY: docker-push			##pushes the local docker images to docker registry. in order to run it need to add REGISTRY and TAG as an arguments.
docker-push: docker-build
	@docker tag ping-pong-sample/ping-vnf-manager ${REGISTRY}/ping-pong-sample/ping-vnf-manager:${TAG};
	@docker tag ping-pong-sample/ping-vnf ${REGISTRY}/ping-pong-sample/ping-vnf:${TAG};
	@docker tag ping-pong-sample/pong-vnf-manager ${REGISTRY}/ping-pong-sample/pong-vnf-manager:${TAG};
	@docker tag ping-pong-sample/pong-vnf ${REGISTRY}/ping-pong-sample/pong-vnf:${TAG};
	@docker push ${REGISTRY}/ping-pong-sample/ping-vnf-manager:${TAG};
	@docker push ${REGISTRY}/ping-pong-sample/ping-vnf:${TAG};
	@docker push ${REGISTRY}/ping-pong-sample/pong-vnf-manager:${TAG};
	@docker push ${REGISTRY}/ping-pong-sample/pong-vnf:${TAG};

.PHONY: clean			##cleans the build area (all executables under build/bin)
clean:
	@rm -rf build/bin

.PHONY: clean-all			##superset of 'clean' that also removes vendor dir and compiled proto go files
clean-all: clean-protos clean-vendor clean

.PHONY: help				##show this help message
help:
	@echo "usage: make [target]\n"; echo "options:"; \fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | sed 's/.PHONY:*//' | sed -e 's/^/  /'; echo "";
