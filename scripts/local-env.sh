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
  - containerPort: 30005
    hostPort: 30005
EOF
kubectl create ns notary
kubectl create ns webhook
kubectl create ns cert-manager
kubectl label namespace notary sighup.io/webhook=ignore
kubectl label namespace webhook sighup.io/webhook=ignore
kubectl label namespace cert-manager sighup.io/webhook=ignore

echo "1. Deploying docker registry"
helm upgrade --install registry stable/docker-registry --values scripts/docker-registry-values.yaml -n notary --version 1.9.4

echo "2. Deploying cert-manager"
retry 10 kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.15.2/cert-manager.crds.yaml
helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --version v0.15.2
kubectl wait --for=condition=Available deployment --timeout=3m -n cert-manager --all

echo "3. Deploying notary"
retry 10 kubectl apply -f scripts/notary-pki.yaml
kubectl wait --for=condition=Ready certs --timeout=3m -n notary --all
kubectl apply -f scripts/notary.yaml
kubectl wait --for=condition=Available deployment --timeout=3m -n notary --all

echo "4. Copying notary-server certificates to webhook namespace"
kubectl get secret notary-server-crt -n notary -o yaml | sed s@"namespace: notary"@"namespace: webhook"@ | kubectl apply -n webhook -f -

echo "5. Generating delegation key"
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: delegation-key
spec:
  secretName: delegation-key
  dnsNames:
    - delegation-key
  issuerRef:
    name: selfsigned-issuer
EOF
kubectl wait --for=condition=Ready certs --timeout=3m --all
kubectl get secret delegation-key -o jsonpath='{.data.tls\.crt}' | base64 -d > delegation.crt
kubectl get secret delegation-key -o jsonpath='{.data.tls\.key}' | base64 -d > delegation.key
chmod 744 delegation.crt
chmod 700 delegation.key
echo "  Delegation key available at ./delegation.crt"

echo "6. Downloading notary server certificate"
kubectl get secret -n notary notary-server-crt -o jsonpath='{.data.tls\.crt}' | base64 -d > notary-tls.crt
echo "  Notary Server certificate available at ./notary-tls.crt"
