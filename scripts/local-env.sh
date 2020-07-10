#!/bin/bash

function retry {
  local retries=$1
  shift

  local count=0
  until "$@"; do
    exit=$?
    wait=$((2 ** $count))
    count=$(($count + 1))
    if [ $count -lt $retries ]; then
      echo "Retry $count/$retries exited $exit, retrying in $wait seconds..."
      sleep $wait
    else
      echo "Retry $count/$retries exited $exit, no more retries left."
      return $exit
    fi
  done
  return 0
}

echo "0. Creating kind cluster"

unset DOCKER_CONTENT_TRUST

cat <<EOF | kind create cluster --image=docker.io/kindest/node:v1.18.4 --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:30001"]
    endpoint = ["http://localhost:30001"]

nodes:
- role: control-plane
- role: worker
  extraPortMappings:
  - containerPort: 30001
    hostPort: 30001
  - containerPort: 30003
    hostPort: 30003
EOF
kubectl create ns notary
kubectl create ns webhook
kubectl create ns cert-manager
kubectl label namespace notary sighup.io/webhook=ignore
kubectl label namespace webhook sighup.io/webhook=ignore
kubectl label namespace cert-manager sighup.io/webhook=ignore

echo "1. Deploying docker registry"
helm upgrade --install registry stable/docker-registry --set service.type=NodePort,service.nodePort=30001 -n notary

echo "2. Deploying cert-manager"
retry 10 kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.15.2/cert-manager.crds.yaml
helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --version v0.15.2
kubectl wait --for=condition=Available deployment --timeout=3m -n cert-manager --all

echo "3. Deploy notary"
retry 10 kubectl apply -f scripts/notary-pki.yaml
kubectl wait --for=condition=Ready certs --timeout=3m -n notary --all
kubectl apply -f scripts/notary.yaml
kubectl wait --for=condition=Available deployment --timeout=3m -n notary --all

echo "4. Finish!"
cat << EOF
Congratulations!!!
Your local environment has been created.

Don't forget to add the following entries in your /etc/hosts file:

  127.0.0.1 registry.local
  127.0.0.1 notary-server.local

registry.local uses port 30001 in your local computer
notary-server.local uses port 30003 in your local computer

Follow the commands bellow to test your setup:

# Clean before tests
$ rm -rf ~/.docker/trust/tuf/registry.local\:30001/

# Get the crt from the notary-server
$ echo | openssl s_client -servername notary-server.local -connect notary-server.local:30003 | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > /tmp/notary-tls.crt

# Init a repository inside notary-server
$ notary -D -p -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert /tmp/notary-tls.crt init registry.local:30001/alpine

# Gen a private key, csr and crt
$ openssl genrsa -out /tmp/delegation.key 2048
$ chmod 700 /tmp/delegation.key
$ openssl req -new -sha256 -key /tmp/delegation.key -out /tmp/delegation.csr
$ openssl x509 -req -sha256 -days 365 -in /tmp/delegation.csr -signkey /tmp/delegation.key -out /tmp/delegation.crt

# Rotate notary repository keys
$ notary -D -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert /tmp/notary-tls.crt key rotate registry.local:30001/alpine snapshot -r
$ notary -D -v -s https://notary-server.local:30003 -d ~/.docker/trust --tlscacert /tmp/notary-tls.crt publish registry.local:30001/alpine

# Pull an example image, tag them sign and push
$ docker pull alpine:3.10
$ docker tag alpine:3.10 registry.local:30001/alpine:3.10
# Set up correct environment variables to enable notary
$ export DOCKER_CONTENT_TRUST=1
$ export DOCKER_CONTENT_TRUST_SERVER=https://notary-server.local:30003
$ docker trust key load /tmp/delegation.key --name jenkins
$ docker trust signer add --key /tmp/delegation.crt jenkins registry.local:30001/alpine
$ docker push registry.local:30001/alpine:3.10
EOF
