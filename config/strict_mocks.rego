package kubernetes.admission.mocks

alpine_3_10_pod_correct_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "",
      "version": "v1",
      "kind": "Pod"
    },
    "resource": {
      "group": "",
      "version": "v1",
      "resource": "pods"
    },
    "requestKind": {
      "group": "",
      "version": "v1",
      "kind": "Pod"
    },
    "requestResource": {
      "group": "",
      "version": "v1",
      "resource": "pods"
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
      "kind": "Pod",
      "apiVersion": "v1",
      "metadata": {
        "name": "debug",
        "creationTimestamp": null,
        "labels": {
          "run": "debug"
        }
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-99tf5",
            "secret": {
              "secretName": "default-token-99tf5"
            }
          }
        ],
        "containers": [
          {
            "name": "debug",
            "image": "alpine@sha256:randomsha",
            "args": [
              "sleep",
              "3600"
            ],
            "resources": {

            },
            "volumeMounts": [
              {
                "name": "default-token-99tf5",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "securityContext": {

        },
        "schedulerName": "default-scheduler",
        "tolerations": [
          {
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          },
          {
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          }
        ],
        "priority": 0,
        "enableServiceLinks": true
      },
      "status": {

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

alpine_3_10_pod_incorrect_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "39b56036-a356-4014-9117-cf64aadbb3f4",
    "kind": {
      "group": "",
      "version": "v1",
      "kind": "Pod"
    },
    "resource": {
      "group": "",
      "version": "v1",
      "resource": "pods"
    },
    "requestKind": {
      "group": "",
      "version": "v1",
      "kind": "Pod"
    },
    "requestResource": {
      "group": "",
      "version": "v1",
      "resource": "pods"
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
      "kind": "Pod",
      "apiVersion": "v1",
      "metadata": {
        "name": "debug",
        "creationTimestamp": null,
        "labels": {
          "run": "debug"
        }
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-99tf5",
            "secret": {
              "secretName": "default-token-99tf5"
            }
          }
        ],
        "containers": [
          {
            "name": "debug",
            "image": "alpine@sha256:malware",
            "args": [
              "sleep",
              "3600"
            ],
            "resources": {

            },
            "volumeMounts": [
              {
                "name": "default-token-99tf5",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "securityContext": {

        },
        "schedulerName": "default-scheduler",
        "tolerations": [
          {
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          },
          {
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          }
        ],
        "priority": 0,
        "enableServiceLinks": true
      },
      "status": {

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

alpine_3_10_cronjob_correct_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "64e52a35-7fe2-4522-b1c4-db6e25cd7bef",
    "kind": {
      "group": "batch",
      "version": "v1beta1",
      "kind": "CronJob"
    },
    "resource": {
      "group": "batch",
      "version": "v1beta1",
      "resource": "cronjobs"
    },
    "requestKind": {
      "group": "batch",
      "version": "v1beta1",
      "kind": "CronJob"
    },
    "requestResource": {
      "group": "batch",
      "version": "v1beta1",
      "resource": "cronjobs"
    },
    "name": "hello",
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
      "kind": "CronJob",
      "apiVersion": "batch/v1beta1",
      "metadata": {
        "name": "hello",
        "namespace": "default",
        "creationTimestamp": null,
        "annotations": {
        }
      },
      "spec": {
        "schedule": "*/1 * * * *",
        "concurrencyPolicy": "Allow",
        "suspend": false,
        "jobTemplate": {
          "metadata": {
            "creationTimestamp": null
          },
          "spec": {
            "template": {
              "metadata": {
                "creationTimestamp": null
              },
              "spec": {
                "containers": [
                  {
                    "name": "hello",
                    "image": "alpine@sha256:randomsha",
                    "args": [
                      "/bin/sh",
                      "-c",
                      "date; echo Hello from the Kubernetes cluster"
                    ],
                    "resources": {

                    },
                    "terminationMessagePath": "/dev/termination-log",
                    "terminationMessagePolicy": "File",
                    "imagePullPolicy": "Always"
                  }
                ],
                "restartPolicy": "OnFailure",
                "terminationGracePeriodSeconds": 30,
                "dnsPolicy": "ClusterFirst",
                "securityContext": {

                },
                "schedulerName": "default-scheduler"
              }
            }
          }
        },
        "successfulJobsHistoryLimit": 3,
        "failedJobsHistoryLimit": 1
      },
      "status": {

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

alpine_3_10_cronjob_incorrect_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "64e52a35-7fe2-4522-b1c4-db6e25cd7bef",
    "kind": {
      "group": "batch",
      "version": "v1beta1",
      "kind": "CronJob"
    },
    "resource": {
      "group": "batch",
      "version": "v1beta1",
      "resource": "cronjobs"
    },
    "requestKind": {
      "group": "batch",
      "version": "v1beta1",
      "kind": "CronJob"
    },
    "requestResource": {
      "group": "batch",
      "version": "v1beta1",
      "resource": "cronjobs"
    },
    "name": "hello",
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
      "kind": "CronJob",
      "apiVersion": "batch/v1beta1",
      "metadata": {
        "name": "hello",
        "namespace": "default",
        "creationTimestamp": null,
        "annotations": {
        }
      },
      "spec": {
        "schedule": "*/1 * * * *",
        "concurrencyPolicy": "Allow",
        "suspend": false,
        "jobTemplate": {
          "metadata": {
            "creationTimestamp": null
          },
          "spec": {
            "template": {
              "metadata": {
                "creationTimestamp": null
              },
              "spec": {
                "containers": [
                  {
                    "name": "hello",
                    "image": "alpine@sha256:malware",
                    "args": [
                      "/bin/sh",
                      "-c",
                      "date; echo Hello from the Kubernetes cluster"
                    ],
                    "resources": {

                    },
                    "terminationMessagePath": "/dev/termination-log",
                    "terminationMessagePolicy": "File",
                    "imagePullPolicy": "Always"
                  }
                ],
                "restartPolicy": "OnFailure",
                "terminationGracePeriodSeconds": 30,
                "dnsPolicy": "ClusterFirst",
                "securityContext": {

                },
                "schedulerName": "default-scheduler"
              }
            }
          }
        },
        "successfulJobsHistoryLimit": 3,
        "failedJobsHistoryLimit": 1
      },
      "status": {

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

alpine_3_10_deployment_correct_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "80dd2506-70ed-4855-b996-215c8cd678bb",
    "kind": {
      "group": "apps",
      "version": "v1",
      "kind": "Deployment"
    },
    "resource": {
      "group": "apps",
      "version": "v1",
      "resource": "deployments"
    },
    "requestKind": {
      "group": "apps",
      "version": "v1",
      "kind": "Deployment"
    },
    "requestResource": {
      "group": "apps",
      "version": "v1",
      "resource": "deployments"
    },
    "name": "nginx-deployment",
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
      "kind": "Deployment",
      "apiVersion": "apps/v1",
      "metadata": {
        "name": "nginx-deployment",
        "namespace": "default",
        "creationTimestamp": null,
        "labels": {
          "app": "nginx"
        },
        "annotations": {
        }
      },
      "spec": {
        "replicas": 3,
        "selector": {
          "matchLabels": {
            "app": "nginx"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "app": "nginx"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "nginx",
                "image": "alpine@sha256:randomsha",
                "ports": [
                  {
                    "containerPort": 80,
                    "protocol": "TCP"
                  }
                ],
                "resources": {

                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "IfNotPresent"
              }
            ],
            "restartPolicy": "Always",
            "terminationGracePeriodSeconds": 30,
            "dnsPolicy": "ClusterFirst",
            "securityContext": {

            },
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": "25%",
            "maxSurge": "25%"
          }
        },
        "revisionHistoryLimit": 10,
        "progressDeadlineSeconds": 600
      },
      "status": {

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

alpine_3_10_deployment_incorrect_sha = {
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "80dd2506-70ed-4855-b996-215c8cd678bb",
    "kind": {
      "group": "apps",
      "version": "v1",
      "kind": "Deployment"
    },
    "resource": {
      "group": "apps",
      "version": "v1",
      "resource": "deployments"
    },
    "requestKind": {
      "group": "apps",
      "version": "v1",
      "kind": "Deployment"
    },
    "requestResource": {
      "group": "apps",
      "version": "v1",
      "resource": "deployments"
    },
    "name": "nginx-deployment",
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
      "kind": "Deployment",
      "apiVersion": "apps/v1",
      "metadata": {
        "name": "nginx-deployment",
        "namespace": "default",
        "creationTimestamp": null,
        "labels": {
          "app": "nginx"
        },
        "annotations": {
        }
      },
      "spec": {
        "replicas": 3,
        "selector": {
          "matchLabels": {
            "app": "nginx"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "app": "nginx"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "nginx",
                "image": "alpine@sha256:malware",
                "ports": [
                  {
                    "containerPort": 80,
                    "protocol": "TCP"
                  }
                ],
                "resources": {

                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "IfNotPresent"
              }
            ],
            "restartPolicy": "Always",
            "terminationGracePeriodSeconds": 30,
            "dnsPolicy": "ClusterFirst",
            "securityContext": {

            },
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": "25%",
            "maxSurge": "25%"
          }
        },
        "revisionHistoryLimit": 10,
        "progressDeadlineSeconds": 600
      },
      "status": {

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