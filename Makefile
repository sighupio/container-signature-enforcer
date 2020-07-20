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


.PHONY: deps
deps:
## dep: download deps to populate cache
	@go mod download -x
	@go get github.com/securego/gosec/cmd/gosec
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0

.PHONY: test
## test: Run local golang tests
test:
	@go test -race -v ./... -cover -count 1

.PHONY: gosec
## gosec: Inspects source code for security problems by scanning the Go AST.
gosec:
	@gosec -out gosec.json -fmt json ./...
	
.PHONY: golangci-lint
## golangci-lint: Go linters aggregator.
golangci-lint:
	@golangci-lint run -v --timeout 5m ./...

.PHONY: build
## build: Builds the opa-notary-connector container image. tests and gosec before building
build: test gosec opa-tests
	@docker build -t opa-notary-connector:latest -f build/Dockerfile .

.PHONY: local-start-script
local-start-script:
	@scripts/local-env.sh
	@echo
	@echo "Congratulations!!!"
	@echo "Your local environment has been created."

.PHONY: local-start
## local-start: Starts kind cluster with everything ready to start developing
local-start: local-start-script local-help

.PHONY: local-help
## local-help: Print some useful commands to execute once local environment is ready.
local-help:
	@echo
	@cat scripts/local-help

.PHONY: local-push
## local-push: Pushes to the local registry the opa-notary-connector container image (local rebuild before push)
local-push: build
	@docker tag opa-notary-connector:latest registry.local:30001/opa-notary-connector:latest
	@docker push registry.local:30001/opa-notary-connector:latest

.PHONY: local-deploy
## local-deploy: Deploys opa-notary-connector using helm. Requires to run make local-push before
local-deploy:
	@kubectl create configmap opa-notary-connector-mode -n webhook --from-file mode.json=config/mode.json --dry-run=client -o yaml | kubectl apply -f - -n webhook
	@kubectl label configmap opa-notary-connector-mode openpolicyagent.org/data=opa --overwrite -n webhook
	@kubectl create configmap opa-notary-connector-rules -n webhook --from-file config/config.rego --from-file config/strict.rego --dry-run=client -o yaml | kubectl apply -f - -n webhook
	@kubectl label configmap opa-notary-connector-rules openpolicyagent.org/policy=rego --overwrite -n webhook
	@helm upgrade --install opa-notary-connector stable/opa --namespace webhook --version 1.14.0 --values scripts/opa-notary-connector-values.yaml --set annotations."deploy-date"="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"
	@kubectl wait --for=condition=Available deployment --timeout=3m -n webhook --all

.PHONY: local-stop
## local-stop: Stops local cluster
local-stop:
	@rm -rf ~/.docker/trust/tuf/localhost\:30001/
	@rm -f delegation.key delegation.crt notary-tls.crt
	@kind delete cluster

.PHONY: mock-server-start
## mock-server-start: Starts a local opa-notary-connector server at 8080 port. Useful to run opa-tests
mock-server-start:
	@npm install -g mockserver
	@mockserver -q -p 8080 -m scripts/mocks/ & echo "$$!" > "/tmp/mockserver.pid"
	@sleep 3

.PHONY: mock-server-stop
## mock-server-stop: Stops the local opa-notary-connector server.
mock-server-stop:
	@kill -9 $(shell cat /tmp/mockserver.pid)

.PHONY: rego-tests
rego-tests:
	@opa test -v config/ || true 

.PHONY: opa-tests
## opa-tests: Runs rego code tests against the mock-server
opa-tests: mock-server-start rego-tests mock-server-stop
