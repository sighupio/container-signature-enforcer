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

.PHONY: local-start
## local-start: Starts kind cluster with everything ready to start developing
local-start:
	@scripts/local-env.sh

.PHONY: local-stop
## local-stop: Stops local cluster
local-stop:
	@kind delete cluster

.PHONY: local-push
## local-push: Pushes to the local registry the opa-notary-connector container image (updated)
local-push: build
	@docker tag opa-notary-connector:latest registry.local:30001/opa-notary-connector:latest
	@docker push registry.local:30001/opa-notary-connector:latest

cleanup:
	helm delete --purge --no-hooks opa-notary-connector
	kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io notary-admission-config
	kubectl delete deployments alpine

install:
	helm install ./opa-notary-connector --name opa-notary-connector --namespace "webhook" --atomic -f values.yaml

run-ok:
	kubectl run alpine --image registry.test/test/alpine:3.10 -- sleep 5000

run-ko-not-matching:
	kubectl run nginx --image nginx:latest
run-ko-not-signed:
	kubectl run nginx --image registry.test/test/alpine:3.9

test:
	go test -race -v ./... -cover

.PHONY: build
build:
	@docker build -t opa-notary-connector:latest -f build/Dockerfile .

sign:
	@echo todo

gosec:
	gosec -out gosec.json -fmt json ./...

# all:
# 	$(MAKE) cleanup
# 	$(MAKE) build
# 	$(MAKE) install
