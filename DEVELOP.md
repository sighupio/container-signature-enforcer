# Development notes

In this docuemnt you will discover how to start a local environment to start contributing to OPA Notary Connector.

## TL;DR

```bash
$ make local-stop local-start local-push local-deploy
$ docker pull alpine:3.10
$ docker tag alpine:3.10 localhost:30001/alpine:3.10
$ docker push localhost:30001/alpine:3.10
$ kubectl run debug --image localhost:30001/alpine:3.10  -- sleep 3600
Error from server (Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine): admission webhook "webhook.openpolicyagent.org" denied the request: Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine
```

## Requirements

### Software

- `docker`
- `kind`
  - `kubectl`
- `helm` 3.1.3.
- `notary` cli
- so related software like `make`, `bash`...

### Ports

The setup requires to have ports `30001` and `30003` available in your computer. These ports will be used by
`docker-registry` and `notary server`.

### `/etc/hosts` entries

Append the following entries in your `/etc/hosts` file. Local environment will issue certificates for these local
domains:

```
127.0.0.1 registry.local
127.0.0.1 notary-server.local
```

## Start the environment

First, you need to start the environment. Base environment consist on:

- docker-registry
- notary
- cert-manager

notary is configured with certificates issued by cert-manager.

To create this base environment you just need to run:

```bash
$ make local-start
0. Creating kind cluster
Creating cluster "kind" ...
<TRUNCATED OUTPUT>
1. Deploying docker registry
<TRUNCATED OUTPUT>
2. Deploying cert-manager
<TRUNCATED OUTPUT>
3. Deploying notary
<TRUNCATED OUTPUT>
4. Copying notary-server certificates to webhook namespace
secret/notary-server-crt created
5. Generating delegation key
issuer.cert-manager.io/selfsigned-issuer created
certificate.cert-manager.io/delegation-key created
certificate.cert-manager.io/delegation-key condition met
configmap/opa-notary-connector-config created
  Delegation key available at ./delegation.crt
6. Downloading notary server certificate
  Notary Server certificate available at ./notary-tls.crt
Congratulations!!!
Your local environment has been created.
<TRUNCATED OUTPUT>
```

After a while you will be able to:

- Push images to registry directly from your terminal. *(No need to expose ports)*
- Use notary from your terminal. *(No need to expose ports)*.

Also three new files are now present:

- `delegation.crt`: Delegation certificate to sign images. Run `make local-help` to discover useful commands
- `delegation.key`: Delegation key to sign images. Run `make local-help` to discover useful commands
- `notary-tls.crt`: Notary server certificate. Required because it is a self-signed certificate.
Run `make local-help` to discover useful commands

### Push an example image (test docker-registry)

In order to test `registry` is working, run the following commands in your terminal:

```bash
$ docker pull alpine:3.10
<TRUNCATED OUTPUT>
$ docker tag alpine:3.10 localhost:30001/alpine:3.10
$ docker push localhost:30001/alpine:3.10
The push refers to repository [localhost:30001/alpine]
1b3ee35aacca: Pushed
3.10: digest: sha256:a143f3ba578f79e2c7b3022c488e6e12a35836cd4a6eb9e363d7f3a07d848590 size: 528
```

Then, run an example pod in the cluster:

```bash
$ kubectl run debug --image localhost:30001/alpine:3.10  -- sleep 3600
pod/debug created
$ kubectl get pod debug
NAME    READY   STATUS    RESTARTS   AGE
debug   1/1     Running   0          9s
```

### Push opa-notary-connector to the registry

Before deploying the solution, we need to build and push the opa-notary-connector image to the registry.

Easy as running:

```bash
$ make local-push
<TRUNCATED OUTPUT>
Successfully built 1a9643854451
Successfully tagged opa-notary-connector:latest
The push refers to repository [registry.local:30001/opa-notary-connector]
9b314de2c980: Pushed
a32fd1ea538e: Pushed
3e207b409db3: Pushed
latest: digest: sha256:5796f3c1324b7e17039836ef1693299b3e62a394921d7eb23b5bb74f62eafd85 size: 946
```

The command above runs all tests (tests, gosec), builds the container image then pushes it to the registry.

### Deploy opa-notary-connector into the cluster

Once the opa-notary-connector is available in the cluster, deploy it running:

```bash
$ make local-deploy
configmap/opa-notary-connector-rules created
Release "opa-notary-connector" does not exist. Installing it now.
<TRUNCATED OUTPUT>
deployment.apps/opa-notary-connector condition met
```

### Test opa-notary-connector

At this point, cluster has configured the opa-notary-connector webhook involved on signature verification. If try to
run a container based on an unsigned container image:

```bash
$ kubectl run debug --image localhost:30001/alpine:3.10  -- sleep 3600
Error from server (Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine): admission webhook "webhook.openpolicyagent.org" denied the request: Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine
```

This error is expected as the image does not have a trusted signature.

#### Add a valid signature

All commands listed bellow are available once you create a cluster using: `make local-start` and also if you run
`make local-help`

```bash
# Init a repository inside notary-server
$ notary -D -p -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert ./notary-tls.crt init localhost:30001/alpine

# Rotate notary repository keys
$ notary -D -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert ./notary-tls.crt key rotate localhost:30001/alpine snapshot -r
$ notary -D -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert ./notary-tls.crt publish localhost:30001/alpine

# Pull an example image, tag them sign and push
$ docker pull alpine:3.10
$ docker tag alpine:3.10 localhost:30001/alpine:3.10
# Set up correct environment variables to enable notary
$ export DOCKER_CONTENT_TRUST=1
$ export DOCKER_CONTENT_TRUST_SERVER=https://notary-server.local:30003
$ docker trust key load ./delegation.key --name jenkins
$ docker trust signer add --key ./delegation.crt jenkins localhost:30001/alpine
$ docker push localhost:30001/alpine:3.10
```

#### Re-Test opa-notary-connector with a valid signature

```bash
$ kubectl run debug --image localhost:30001/alpine:3.10  -- sleep 3600
pod/debug created
```

Now it works. `jenkins` signed the container image.
Try with a different image:

```bash
$ export DOCKER_CONTENT_TRUST=0
$ docker pull alpine:3.11
$ docker tag alpine:3.11 localhost:30001/alpine:3.11
$ docker push localhost:30001/alpine:3.11
$ kubectl run debug --image localhost:30001/alpine:3.11  -- sleep 3600
Error from server (Container image localhost:30001/alpine:3.11 invalid: No valid trust data for 3.11): admission webhook "webhook.openpolicyagent.org" denied the request: Container image localhost:30001/alpine:3.11 invalid: No valid trust data for 3.11
```

Now a new error appeared. Caused by an unsigned image.

## Destroy local environment

To destroy the local environment:

```
$ make local-stop
Deleting cluster "kind" ...
```

It will delete the following files and directories:

- `~/.docker/trust/tuf/localhost\:30001/`
- `delegation.key`
- `delegation.crt`
- `notary-tls.crt`

As they are fully dependant of this local environment.

## Development

### golang

Having the local environment working, if you want to change some golang code and pass some integration tests
(automated or manually) in your local setup you will need to:

```bash
$ make local-push local-deploy
```

After a while you will be able to see your changes in the cluster

### rego

Modify rego code has to be managed in two separate ways

#### Main

Main code is where kubernetes will call once a new requests matches the configuration.
The main code is available [`scripts/opa-notary-connector-values.yaml`](scripts/opa-notary-connector-values.yaml) in
the `bootstrapPolicies` attribute.

This code should not change a lot as it is the interface between Kubernetes and opa-notary-connector.
In case you change it run:

```bash
$ make local-deploy
```

After a while you will be able to see your changes in the cluster

#### Rules

Rego code is available in the [`scripts/opa-notary-connector-config.yaml`](scripts/opa-notary-connector-config.yaml)
file.

This code is in charge of get the image to run, pass it to the golang project and decide what to do next:

- deny: Deny in case opa does not contains a valid signature.
- patch: Patch the request with the image sha to be inmutable.

In case you change it run:

```bash
$ make local-deploy
```

After a while you will be able to see your changes in the cluster

## Special considerations for Minishift

For the admission webhook to work in Minishift it's necessary to enable the `admissions-webhook` addon with the following command:

```shell
minishift addons enable admissions-webhook
```

In the Minishift version we tested (`v1.34.1+c2ff9cb`) the command doesn't work, the changes have to be done manually. From the [addon source code file](
https://github.com/minishift/minishift/blob/master/addons/admissions-webhook/admissions-webhook.addon) we get that the steps to be done are as follows:

```shell
# Login via ssh to the Minishift MV:
minishift ssh

# Make a backup of kubes apiserver master-config.yaml
cd /var/lib/minishift/base/kube-apiserver
cp master-config.yaml master-config-tmp.yaml

# Patch the configuration
/var/lib/minishift/bin/oc ex config patch master-config-tmp.yaml --patch "$(curl https://raw.githubusercontent.com/minishift/minishift/master/addons/admissions-webhook/patch.json)" > master-config.yaml

# Stop the containers of the apiserver and the api, kubernetes should restart them by its own afterwards
docker stop $(docker ps -l -q --filter "label=io.kubernetes.container.name=apiserver")
docker stop $(docker ps -l -q --filter "label=io.kubernetes.container.name=api")
```