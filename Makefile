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

build:
	docker build -t opa-notary-connector:latest -f build/Dockerfile .

sign:
	@echo todo

gosec:
	gosec -out gosec.json -fmt json ./...

all:
	$(MAKE) cleanup
	$(MAKE) build
	$(MAKE) install
