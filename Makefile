.DEFAULT_GOAL: help
SHELL := /bin/bash

PROJECTNAME := $(shell basename "$(PWD)")

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: test
## test: Run local golang tests
test:
	go test -race -v ./... -cover

.PHONY: gosec
## gosec: Inspects source code for security problems by scanning the Go AST.
gosec:
	gosec -out gosec.json -fmt json ./...

.PHONY: build
## build: Builds the opa-notary-connector container image
build: test gosec
	@docker build -t opa-notary-connector:latest -f build/Dockerfile .


.PHONY: local-start
## local-start: Starts kind cluster with everything ready to start developing
local-start:
	@scripts/local-env.sh

.PHONY: local-stop
## local-stop: Stops local cluster
local-stop:
	@rm -f delegation.key delegation.crt notary-tls.crt
	@kind delete cluster

.PHONY: local-push
## local-push: Pushes to the local registry the opa-notary-connector container image (updated)
local-push: build
	@docker tag opa-notary-connector:latest registry.local:30001/opa-notary-connector:latest
	@docker push registry.local:30001/opa-notary-connector:latest

.PHONY: local-deploy
## local-deploy: Deploys opa-notary-connector using helm
local-deploy:
	@kubectl apply -f scripts/opa-notary-connector-config.yaml
	@helm upgrade --install opa-notary-connector stable/opa --namespace webhook --version 1.14.0 --values scripts/opa-notary-connector-values.yaml
	@kubectl wait --for=condition=Available deployment --timeout=3m -n webhook --all

