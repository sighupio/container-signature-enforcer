# OPA Notary Connector

[![Build Status](https://ci.sighup.io/api/badges/sighupio/opa-notary-connector/status.svg?ref=refs/tags/v0.1.1)](https://ci.sighup.io/sighupio/opa-notary-connector)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/opa-notary-connector)](https://artifacthub.io/packages/search?repo=opa-notary-connector)

## Idea

The aim of the project was to reimplement [ibm/portieris](https://github.com/IBM/portieris) main features
*(replacing IBM related parts)*.

## What it does

It validate images' signatures against notary instances, rejecting untrusted images and setting the
`sha256` to trusted images in order to pin the exact version.
The enforcing is configurable per namespace and registry.

## Requirements

OPA Notary Connector fits in a running environment with:

- Kubernetes Cluster > 1.16
- Container Registry
- Notary Server

You will need to have:
- Notary delegation key to sign and validate container image.
- Notary server TLS certificate and/or ca certificate that signed the Notary certificate.

## Config

A [Helm Chart is provided](deployments/helm/opa-notary-connector) and a minimal configuration can be found
[here](scripts/opa-notary-connector-values.yaml) and the complete default values
[here](deployments/helm/opa-notary-connector/values.yaml).

The configuration consists of a list of *repositories*, which can be applied via regexes to restrict their scope to
specific matching images.

Each repository can also have a *priority* associated, that in case of multiple matching repositories will be used to
decide the policy to be applied. The *trust* section describes the policy that has to be enforced on the matched
repository: whether signatures have to be *enabled* or not; the *trust server* address; a list of signers that
have to have signed the image in order to be able to accept it.

### `values.yaml` Deep dive

```yaml
# Strict mode
strict: true

# OPA Notary Connector configuration
repositories: []
# - name: "localhost.*"
#   priority: 10
#   trust:
#     enabled: true
#     trustServer: "https://notary-server.notary.svc.cluster.local:4443"
#     signers:
#     - role: "targets/jenkins"
#       publicKey: "H" # base64 encoded public key


opa:
  opa: false

  certManager:
    enabled: true
<<TRUNCATED FILE>>
  admissionControllerNamespaceSelector:
    matchExpressions:
      - {key: sighup.io/webhook, operator: NotIn, values: [ignore]}
<<TRUNCATED FILE>>
  mgmt:
    configmapPolicies:
      namespaces: ["webhook"]
```

First two variables configures the OPA Notary Connector.

- `strict`. Defines what to do in case of invalid SHA. Deploying container with an invalid SHA:
  - `strict: true`. Denied the deployment of the container as it could be an hijacked image.
  - `strict: false`. Replace the SHA with the correct one (if signature is available at Notary).
- `repositories:` Explained before. It indicates how and when and image should be validated against trust servers
*(Notary)*.
- `opa`. Allows you to override any
[OPA helm chart configuration](https://github.com/helm/charts/blob/master/stable/opa/values.yaml). This chart ships a
preconfigured opa subchart configuration while exposing its configuration parameters.

## Intended behavior in specific use cases

Deployed and configured the opa-notary-connector in the target cluster, when a new resource containing one or more images' references is received the webhook will behave as follows:

| Namespace | Image | Trust | Signatures | Accepted |
| :-------: | :---: | :---: | :--------: | :------: |
|     ❌     |   ❓   |   ❓   |     ❓      |    ✅     |
|     ✅     |   ❌   |   ❓   |     ❓      |    ❌     |
|     ✅     |   ✅   |   ❌   |     ❓      |    ✅     |
|     ✅     |   ✅   |   ✅   |     ✅      |    ✅     |
|     ✅     |   ✅   |   ✅   |     ❌      |    ❌     |

In case of more than one repositories matching, the one with highest priority will be used.
In case of more than one image specified in a single resource, all images have to be allowed for the request to be accepted.

### References

**Namespace**: is the namespace matching the admissionControllerNames configured?

**Image**: is the image matching one of the repositories configured?

**Trust**: is trust enabled for the repository matched?

**Signatures**: are configured signers recognized?

**Accepted**: is the request accepted?

❌ : **No**

✅ : **Yes**

❓ : **Whatever**

## Getting started (Local)

The fastest option to try this `opa-notary-connector` in your local machine is to follow [`DEVELOP.md`](DEVELOP.md)
file.

### TL;DR

```bash
$ make local-stop local-start local-push local-deploy
$ docker pull alpine:3.10
$ docker tag alpine:3.10 localhost:30001/alpine:3.10
$ docker login -u admin -p admin localhost:30001
$ docker push localhost:30001/alpine:3.10
$ kubectl run debug --image localhost:30001/alpine:3.10  -- sleep 3600
Error from server (Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine): admission webhook "webhook.openpolicyagent.org" denied the request: Container image localhost:30001/alpine:3.10 invalid: notary-server.notary.svc.cluster.local:4443 does not have trust data for localhost:30001/alpine
```

## LICENSE

For license details please see [LICENSE](LICENSE)
