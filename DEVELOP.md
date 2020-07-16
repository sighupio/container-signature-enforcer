# Developer notes

In this notes you will discover how to start a local environment to start contributing to this OPA Notary Connector.

## Requirements

### Software

- docker
- kind

### Ports

This setup requires to have ports 30001 and 30003 available in your computer. These ports will be used by
docker-registry and notary servers.

### `/etc/hosts` entries

Append the following entries in your `/etc/hosts` file. Local environment will issue certificates for these local
domains:

```
127.0.0.1 registry.local
127.0.0.1 notary-server.local
```
