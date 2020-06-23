# Notary Admission Webhook

The aim of the project is to reimplement the main features from [ibm/portieris](https://github.com/IBM/portieris), removing the IBM related parts.

## Logic

The main objective of this custom admission webhook is to validate images' signatures against notary instances, rejecting untrusted images and setting the `sha256` to trusted images in order to pin the exact version. The enforcing is configurable per namespace and registry.

## Config

A Helm Chart is provided and a minimal configuration can be found [here](values.yaml) and the complete default values [here](opa-notary-connector/values.yaml) at `notaryWebhook.config`.

The configuration consists of a list of *repositories*, which can be applied via regexes to restrict their scope to specific matching images and/or namespaces.
Each repository can also have a *priority* associated, that in case of multiple matching repositories will be used to decide the policy to be applied.
The *trust* section describes the policy that has to be enforced on the matched repository: whether signatures have to be *enabled* or not; the *trust server* address; a list of signers that have to have signed the image in order to be able to accept it.

## Intended behavior in specific use cases

Deployed and configured the opa-notary-connector in the target cluster, when a new resource containing one or more images' references is received the webhook will behave as follows:

| Namespace   | Image  | Trust  | Signatures | Accepted |
|:---:|:---:|:---:|:---:|:---:|
| ❌  | ❓ | ❓ | ❓ | ✅ |
| ✅  | ❌ | ❓ | ❓ | ❌ |
| ✅  | ✅ | ❌ | ❓ | ✅ |
| ✅  | ✅ | ✅ | ✅ | ✅ |
| ✅  | ✅ | ✅ | ❌ | ❌ |

In case of more than one repositories matching, the one with highest priority will be used.
In case of more than one image specified in a single resource, all images have to be allowed for the request to be accepted.

### References

**Namespace**: is the namespace matching one of the repositories configured?

**Image**: is the image matching one of the repositories configured?

**Trust**: is trust enabled for the repository matched?

**Signatures**: are configured signers recognized?

**Accepted**: is the request accepted?

❌ : **No**

✅ : **Yes**

❓ : **Whatever**

## How to try

Needed components:

- notary server
- opa-notary-connector admission webhook
- docker daemon configured
- running registry

Steps:

1. deploy registry
1. deploy notary server
1. generate notary signing key pair
1. initialize repository on notary
1. setup docker daemon to point to notary server
1. build and push image
1. customize the configuration of the webhook (`values.yaml`):
    - notary TLS certificate
    - repositories' policies (public signatures used during image build / signing)
    - TLS CA for webhook
1. deploy webhook
1. successfully deploy built and signed image

```shell
# deploy registry
docker run -d -p 80:5000 --name registry registry:2

# deploy notary
helm install notary --namespace webhook -n notary
kubectl port-forward -n webhook svc/notary-server 4443:4443 > /dev/null &

# add the following lines  to the /etc/hosts file:
# 127.0.0.1 notary-server # needed by the tls certificate of notary-server
# 127.0.0.1 registry.test

# initialize repository
notary -D -p -v -s https://notary-server:4443 -d ~/.docker/trust --tlscacert certs/notary-tls.crt init registry.test/test/alpine

# generate key pair for image signature (https://docs.docker.com/engine/security/trust/trust_delegation/#manually-generating-keys)
openssl genrsa -out certs/delegation.key 2048
openssl req -new -sha256 -key certs/delegation.key -out certs/delegation.csr
openssl x509 -req -sha256 -days 365 -in certs/delegation.csr -signkey certs/delegation.key -out certs/delegation.crt

# rotate keys (https://docs.docker.com/v17.09/datacenter/ucp/2.0/guides/content-trust/continuous-integration/#enable-content-trust)
notary -D -v -s https://notary-server:4443 -d ~/.docker/trust --tlscacert certs/notary-tls.crt key rotate registry.test/test/alpine snapshot -r
notary -D -v -s https://notary-server:4443 -d ~/.docker/trust --tlscacert certs/notary-tls.crt publish registry.test/test/alpine

# pull a random image and tag it for our registry
docker pull alpine:3.10
docker tag alpine:3.10 registry.test/test/alpine:3.10

# export needed env var for docker
export DOCKER_CONTENT_TRUST=1
export DOCKER_CONTENT_TRUST_SERVER=https://notary-server:4443

docker trust key load certs/delegation.key --name jenkins
docker trust signer add --key certs/delegation.crt jenkins registry.test/test/alpine

# pull, tag and push image to local registry signing image
docker push registry.test/test/alpine:3.10

# build image opa-notary-connector image
make build

# add the public key (./certs/delegation.crt) to the values.yaml file at the root of this project
base64 certs/delegation.crt >> values.yaml # and fix spaces

# install webhook
make install

# test image signed should pass
kubectl run alpine -n webhook --image registry.test/test/alpine:3.10 -- sleep 5000

# test image not signed, should fail
kubectl run alpine-ko -n webhook --image registry.test/test/alpine:3.9 -- sleep 5000

# to cleanup
make cleanup
```

## How to production ready

1. generate dedicated tls certificates for:
    - notary server
    - opa-notary-connector
1. generate keys for each repository in the notary server
1. store securely all generated keys

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

## Demo

Having the example [values.yaml](values.yaml) applied and image `registry.test/test/alpine:3.10` pushed to the registry and signed by the from `targets/jenkins` using the keys in [certs](certs) the following command should successfully create a deployment:

```bash
kubectl run alpine -n webhook --image registry.test/test/alpine:3.10 -- sleep 5000
```

Instead, the following commands should be rejected:

```bash
kubectl run test-ko-1 -n webhook --image registry.test/test/alpine:3.9 -- sleep 5000
kubectl run test-ko-2 -n webhook --image nginx -- sleep 5000
```

Given that we are trying to deploy to a namespace covered by some policy but the webhook is not able to verify the signature of both images.

If we edit the config to allow untrusted images if not matching any other policy, a sort of "catch all" Allow policy with lower priority than all other policies defined, the two previous images should be allowed to be deployed:

```yaml
- name: '.*'
  namespace: "webhook"
  priority: 0
  trust:
    enabled: false
```

If instead we modify the signers public key in the config with a still valid but different base64 encoded public key (an invalid public key would not be loaded):

```yaml
    - name: 'registry.*'
      namespace: "webhook"
      priority: 10
      trust:
        enabled: true
        trustServer: "https://notary-server:4443"
        signers:
        - role: "targets/jenkins"
          publicKey: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURCakNDQWU0Q0NRRHlNeEhRRG5WTWd6QU5CZ2txaGtpRzl3MEJBUXNGQURCRk1Rc3dDUVlEVlFRR0V3SkIKVlRFVE1CRUdBMVVFQ0F3S1UyOXRaUzFUZEdGMFpURWhNQjhHQTFVRUNnd1lTVzUwWlhKdVpYUWdWMmxrWjJsMApjeUJRZEhrZ1RIUmtNQjRYRFRFNU1URXlOVEE1TlRBME9Wb1hEVEl3TVRFeU5EQTVOVEEwT1Zvd1JURUxNQWtHCkExVUVCaE1DUVZVeEV6QVJCZ05WQkFnTUNsTnZiV1V0VTNSaGRHVXhJVEFmQmdOVkJBb01HRWx1ZEdWeWJtVjAKSUZkcFpHZHBkSE1nVUhSNUlFeDBaRENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQgpBTGcrYW1QNTVESjNaYUdLOTdSSXlLaEc5L1I1TEltVWJjaEpraHcrMFdoZUNtUm1HK1M2SlkvbG9DbXhKcTVOClJ3N0lKYmhYQytHRm13MitsTVlZS1I0QjY0UTZVRkdreVN6cndFMWtzVU5JbXZkL3dCRmtmNUJpb3g3eUFWSTAKZUx0T0V0dGdCZUxRb3JaWi8yQWRNYlpSQjFGQ3craXYvaWs3SDJLcGhJdDg0bWNmOXhoUDI5Wmcvcyt5aHVSUQo5bm5yTnNNRUQrNkZYald1QlI4aFZmanhZcHlPUmdWeUVZSDdJZXhLWkR6ckZjMHZvdlNXbVkvTURJZUozN3VHCnpvcU5SMUxGeEtMblduYzNubXZWUXJpajJ1VENjSWVtYW90MW95Z0ZMMXJFRzR2aGxTTjVhT3YyOGFqNjRJNVMKS09mTUU2MU9PODdaZlBYcmoxNFJLdWtDQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFUMEFuek1DNApNdVFoQk5lS1BBV2s2QnllR041ckdHL0hUQStKREhzOHIxb3lRRldpRlJWZ0FmbXNXMlVEV0JvWVN6VGVUdFpUCjMvTDJ1RFg4UmNweEtWb1RoeVRxeDdOY04rZE1lK3BLSHEvbkZqekdrTFR3clI4UkRtQ1RkNXN1SEhlUzZTNm0KOEFFaC9oTVpaaVRqM241czdzRGFuamorYWowcklsNVEyanNOaXFOanUraS9odDcvemJYRnRSV3RFMDJEREpKeQpCL01MODVTcDdEUHZOLzB3SEg0dUxYU0hZRnZ3ODhONEJJMmZjM0FsR1R0QTdDQ0MyZCtiRzNRUDI0UmdDaUpSClhjTit2VWlIc3JncHMxcDMvUkRJQTRQYURGOXdtYXgxM2tqeHZCVTZxSzg4UFdqdzlUbThNNzREWU4wd21FU0kKLytmajJ3aG1oNDg1aGc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
```

And try to deploy the first signed image again:

```bash
kubectl run alpine -n webhook --image registry.test/test/alpine:3.10 -- sleep 5000
```

This will fail because the image has not been signed by the specified key.
