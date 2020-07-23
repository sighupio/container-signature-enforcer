package kubernetes.admission.mocks

alpine_3_10_minio_tenant = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "resource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "requestKind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "requestResource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "name": "debug",
    "namespace": "default",
    "operation": "CREATE",
    "userInfo": {
      "username": "kubernetes-admin",
      "groups": [
        "system:masters",
        "system:authenticated"
      ]
    },
    "object": {
      "apiVersion": "minio.min.io/v1",
      "kind": "Tenant",
      "metadata": {
        "name": "minio"
      },
      "spec": {
        "metadata": {
          "labels": {
            "app": "minio"
          },
          "annotations": {
            "prometheus.io/path": "/minio/prometheus/metrics",
            "prometheus.io/port": "9000",
            "prometheus.io/scrape": "true"
          }
        },
        "image": "alpine:3.10",
        "serviceName": "minio-internal-service",
        "zones": [
          {
            "servers": 4,
            "volumesPerServer": 4,
            "volumeClaimTemplate": {
              "metadata": {
                "name": "data"
              },
              "spec": {
                "accessModes": [
                  "ReadWriteOnce"
                ],
                "resources": {
                  "requests": {
                    "storage": "1Ti"
                  }
                }
              }
            }
          }
        ],
        "mountPath": "/export",
        "credsSecret": {
          "name": "minio-creds-secret"
        },
        "podManagementPolicy": "Parallel",
        "requestAutoCert": false,
        "certConfig": {
          "commonName": "",
          "organizationName": [
            
          ],
          "dnsNames": [
            
          ]
        },
        "liveness": {
          "initialDelaySeconds": 10,
          "periodSeconds": 1,
          "timeoutSeconds": 1
        }
      }
    },
    "oldObject": null,
    "dryRun": false,
    "options": {
      "kind": "CreateOptions",
      "apiVersion": "meta.k8s.io/v1"
    }
  }
}

alpine_3_11_minio_tenant = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "resource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "requestKind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "requestResource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "name": "debug",
    "namespace": "default",
    "operation": "CREATE",
    "userInfo": {
      "username": "kubernetes-admin",
      "groups": [
        "system:masters",
        "system:authenticated"
      ]
    },
    "object": {
      "apiVersion": "minio.min.io/v1",
      "kind": "Tenant",
      "metadata": {
        "name": "minio"
      },
      "spec": {
        "metadata": {
          "labels": {
            "app": "minio"
          },
          "annotations": {
            "prometheus.io/path": "/minio/prometheus/metrics",
            "prometheus.io/port": "9000",
            "prometheus.io/scrape": "true"
          }
        },
        "image": "alpine:3.11",
        "serviceName": "minio-internal-service",
        "zones": [
          {
            "servers": 4,
            "volumesPerServer": 4,
            "volumeClaimTemplate": {
              "metadata": {
                "name": "data"
              },
              "spec": {
                "accessModes": [
                  "ReadWriteOnce"
                ],
                "resources": {
                  "requests": {
                    "storage": "1Ti"
                  }
                }
              }
            }
          }
        ],
        "mountPath": "/export",
        "credsSecret": {
          "name": "minio-creds-secret"
        },
        "podManagementPolicy": "Parallel",
        "requestAutoCert": false,
        "certConfig": {
          "commonName": "",
          "organizationName": [
            
          ],
          "dnsNames": [
            
          ]
        },
        "liveness": {
          "initialDelaySeconds": 10,
          "periodSeconds": 1,
          "timeoutSeconds": 1
        }
      }
    },
    "oldObject": null,
    "dryRun": false,
    "options": {
      "kind": "CreateOptions",
      "apiVersion": "meta.k8s.io/v1"
    }
  }
}

alpine_3_10_minio_tenant_correct_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "resource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "requestKind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "requestResource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "name": "debug",
    "namespace": "default",
    "operation": "CREATE",
    "userInfo": {
      "username": "kubernetes-admin",
      "groups": [
        "system:masters",
        "system:authenticated"
      ]
    },
    "object": {
      "apiVersion": "minio.min.io/v1",
      "kind": "Tenant",
      "metadata": {
        "name": "minio"
      },
      "spec": {
        "metadata": {
          "labels": {
            "app": "minio"
          },
          "annotations": {
            "prometheus.io/path": "/minio/prometheus/metrics",
            "prometheus.io/port": "9000",
            "prometheus.io/scrape": "true"
          }
        },
        "image": "alpine@sha256:randomsha",
        "serviceName": "minio-internal-service",
        "zones": [
          {
            "servers": 4,
            "volumesPerServer": 4,
            "volumeClaimTemplate": {
              "metadata": {
                "name": "data"
              },
              "spec": {
                "accessModes": [
                  "ReadWriteOnce"
                ],
                "resources": {
                  "requests": {
                    "storage": "1Ti"
                  }
                }
              }
            }
          }
        ],
        "mountPath": "/export",
        "credsSecret": {
          "name": "minio-creds-secret"
        },
        "podManagementPolicy": "Parallel",
        "requestAutoCert": false,
        "certConfig": {
          "commonName": "",
          "organizationName": [
            
          ],
          "dnsNames": [
            
          ]
        },
        "liveness": {
          "initialDelaySeconds": 10,
          "periodSeconds": 1,
          "timeoutSeconds": 1
        }
      }
    },
    "oldObject": null,
    "dryRun": false,
    "options": {
      "kind": "CreateOptions",
      "apiVersion": "meta.k8s.io/v1"
    }
  }
}

alpine_3_10_minio_tenant_incorrect_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "resource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "requestKind": {
      "group": "minio.min.io",
      "version": "v1",
      "kind": "Tenant"
    },
    "requestResource": {
      "group": "minio.min.io",
      "version": "v1",
      "resource": "tenants"
    },
    "name": "debug",
    "namespace": "default",
    "operation": "CREATE",
    "userInfo": {
      "username": "kubernetes-admin",
      "groups": [
        "system:masters",
        "system:authenticated"
      ]
    },
    "object": {
      "apiVersion": "minio.min.io/v1",
      "kind": "Tenant",
      "metadata": {
        "name": "minio"
      },
      "spec": {
        "metadata": {
          "labels": {
            "app": "minio"
          },
          "annotations": {
            "prometheus.io/path": "/minio/prometheus/metrics",
            "prometheus.io/port": "9000",
            "prometheus.io/scrape": "true"
          }
        },
        "image": "alpine@sha256:malware",
        "serviceName": "minio-internal-service",
        "zones": [
          {
            "servers": 4,
            "volumesPerServer": 4,
            "volumeClaimTemplate": {
              "metadata": {
                "name": "data"
              },
              "spec": {
                "accessModes": [
                  "ReadWriteOnce"
                ],
                "resources": {
                  "requests": {
                    "storage": "1Ti"
                  }
                }
              }
            }
          }
        ],
        "mountPath": "/export",
        "credsSecret": {
          "name": "minio-creds-secret"
        },
        "podManagementPolicy": "Parallel",
        "requestAutoCert": false,
        "certConfig": {
          "commonName": "",
          "organizationName": [
            
          ],
          "dnsNames": [
            
          ]
        },
        "liveness": {
          "initialDelaySeconds": 10,
          "periodSeconds": 1,
          "timeoutSeconds": 1
        }
      }
    },
    "oldObject": null,
    "dryRun": false,
    "options": {
      "kind": "CreateOptions",
      "apiVersion": "meta.k8s.io/v1"
    }
  }
}
