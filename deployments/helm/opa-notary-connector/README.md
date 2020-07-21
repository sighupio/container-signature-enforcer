opa-notary-connector
====================
OPA Notary Connector helm chart

Current chart version is `0.1.0`



## Chart Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://kubernetes-charts.storage.googleapis.com | opa | 1.14.0 |

## Chart Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| opa.admissionControllerFailurePolicy | string | `"Fail"` |  |
| opa.admissionControllerKind | string | `"MutatingWebhookConfiguration"` |  |
| opa.admissionControllerNamespaceSelector.matchExpressions[0].key | string | `"sighup.io/webhook"` |  |
| opa.admissionControllerNamespaceSelector.matchExpressions[0].operator | string | `"NotIn"` |  |
| opa.admissionControllerNamespaceSelector.matchExpressions[0].values[0] | string | `"ignore"` |  |
| opa.admissionControllerRules[0].apiGroups[0] | string | `"*"` |  |
| opa.admissionControllerRules[0].apiVersions[0] | string | `"*"` |  |
| opa.admissionControllerRules[0].operations[0] | string | `"CREATE"` |  |
| opa.admissionControllerRules[0].operations[1] | string | `"UPDATE"` |  |
| opa.admissionControllerRules[0].resources[0] | string | `"pods"` |  |
| opa.admissionControllerRules[0].resources[1] | string | `"deployments"` |  |
| opa.admissionControllerRules[0].resources[2] | string | `"replicationcontrollers"` |  |
| opa.admissionControllerRules[0].resources[3] | string | `"replicasets"` |  |
| opa.admissionControllerRules[0].resources[4] | string | `"daemonsets"` |  |
| opa.admissionControllerRules[0].resources[5] | string | `"statefulsets"` |  |
| opa.admissionControllerRules[0].resources[6] | string | `"jobs"` |  |
| opa.admissionControllerRules[0].resources[7] | string | `"cronjobs"` |  |
| opa.bootstrapPolicies.main | string | `"package system\n\nimport data.kubernetes.admission\n\nmain = {\n  \"apiVersion\": \"admission.k8s.io/v1beta1\",\n  \"kind\": \"AdmissionReview\",\n  \"response\": response,\n}\n\ndefault response = {\"allowed\": false, \"status\": {\"reason\": \"Strict mode enabled\"}}\n\nresponse = {\n  \"allowed\": false,\n  \"status\": {\"reason\": reason},\n} {\n  count(admission.deny) > 0\n  reason := concat(\"\\n\", admission.deny)\n}\n\nresponse = {\n  \"allowed\": true,\n  \"patchType\": \"JSONPatch\",\n  \"patch\": patch_bytes,\n} {\n  count(admission.deny) == 0\n  patch := {xw | xw := admission.patches[_][_]}\n  patch_json := json.marshal(patch)\n  patch_bytes := base64.encode(patch_json)\n  patch_bytes != \"W10=\"\n}\n\nresponse = {\n  \"allowed\": false,\n  \"status\": {\"reason\": patch_reason},\n} {\n  count(admission.deny) == 0\n  patch = {xw | xw := admission.patches[_][_]}\n  patch_json := json.marshal(patch)\n  patch_bytes := base64.encode(patch_json)\n  patch_bytes == \"W10=\"\n  patch_reason := \"OPA Notary Connector didn't return a valid value. Look at its logs to debug it\"\n}"` |  |
| opa.certManager.enabled | bool | `true` |  |
| opa.extraContainers[0].args[0] | string | `"--config=/etc/opa-notary-connector/trust.yaml"` |  |
| opa.extraContainers[0].args[1] | string | `"--listen-address=:8080"` |  |
| opa.extraContainers[0].args[2] | string | `"--trust-root-dir=/etc/opa-notary-connector/.trust"` |  |
| opa.extraContainers[0].args[3] | string | `"--verbosity=debug"` |  |
| opa.extraContainers[0].command[0] | string | `"/opa-notary-connector"` |  |
| opa.extraContainers[0].env[0].name | string | `"GIN_MODE"` |  |
| opa.extraContainers[0].env[0].value | string | `"debug"` |  |
| opa.extraContainers[0].image | string | `"localhost:30001/opa-notary-connector:latest"` |  |
| opa.extraContainers[0].imagePullPolicy | string | `"Always"` |  |
| opa.extraContainers[0].livenessProbe.httpGet.path | string | `"/healthz"` |  |
| opa.extraContainers[0].livenessProbe.httpGet.port | string | `"http"` |  |
| opa.extraContainers[0].livenessProbe.httpGet.scheme | string | `"HTTP"` |  |
| opa.extraContainers[0].name | string | `"opa-notary-connector"` |  |
| opa.extraContainers[0].ports[0].containerPort | int | `8080` |  |
| opa.extraContainers[0].ports[0].name | string | `"http"` |  |
| opa.extraContainers[0].ports[0].protocol | string | `"TCP"` |  |
| opa.extraContainers[0].readinessProbe.httpGet.path | string | `"/healthz"` |  |
| opa.extraContainers[0].readinessProbe.httpGet.port | string | `"http"` |  |
| opa.extraContainers[0].readinessProbe.httpGet.scheme | string | `"HTTP"` |  |
| opa.extraContainers[0].securityContext.runAsUser | int | `1001` |  |
| opa.extraContainers[0].volumeMounts[0].mountPath | string | `"/etc/opa-notary-connector/trust.yaml"` |  |
| opa.extraContainers[0].volumeMounts[0].name | string | `"opa-notary-connector-config"` |  |
| opa.extraContainers[0].volumeMounts[0].subPath | string | `"trust.yaml"` |  |
| opa.extraContainers[0].volumeMounts[1].mountPath | string | `"/etc/ssl/certs/ca.crt"` |  |
| opa.extraContainers[0].volumeMounts[1].name | string | `"notary-server-crt"` |  |
| opa.extraContainers[0].volumeMounts[1].subPath | string | `"ca.crt"` |  |
| opa.extraVolumes[0].configMap.name | string | `"opa-notary-connector-config"` |  |
| opa.extraVolumes[0].name | string | `"opa-notary-connector-config"` |  |
| opa.extraVolumes[1].name | string | `"notary-server-crt"` |  |
| opa.extraVolumes[1].secret.secretName | string | `"notary-server-crt"` |  |
| opa.imagePullPolicy | string | `"Always"` |  |
| opa.imageTag | string | `"0.21.1"` |  |
| opa.livenessProbe.httpGet.port | int | `8443` |  |
| opa.mgmt.configmapPolicies.enabled | bool | `true` |  |
| opa.mgmt.configmapPolicies.namespaces[0] | string | `"webhook"` |  |
| opa.mgmt.configmapPolicies.requireLabel | bool | `true` |  |
| opa.mgmt.data.enabled | bool | `true` |  |
| opa.mgmt.imagePullPolicy | string | `"Always"` |  |
| opa.mgmt.imageTag | string | `"0.11"` |  |
| opa.opa | bool | `false` |  |
| opa.port | int | `8443` |  |
| opa.rbac.rules.cluster[0].apiGroups[0] | string | `"*"` |  |
| opa.rbac.rules.cluster[0].resources[0] | string | `"*"` |  |
| opa.rbac.rules.cluster[0].verbs[0] | string | `"get"` |  |
| opa.rbac.rules.cluster[0].verbs[1] | string | `"list"` |  |
| opa.rbac.rules.cluster[0].verbs[2] | string | `"watch"` |  |
| opa.rbac.rules.cluster[1].apiGroups[0] | string | `""` |  |
| opa.rbac.rules.cluster[1].resources[0] | string | `"configmaps"` |  |
| opa.rbac.rules.cluster[1].verbs[0] | string | `"update"` |  |
| opa.rbac.rules.cluster[1].verbs[1] | string | `"patch"` |  |
| opa.readinessProbe.httpGet.port | int | `8443` |  |
| opa.securityContext.enabled | bool | `true` |  |
| opa.securityContext.fsGroup | int | `1001` |  |
| opa.securityContext.runAsNonRoot | bool | `true` |  |
| opa.securityContext.runAsUser | int | `1` |  |
| repositories | list | `[]` |  |
| strict | bool | `true` |  |
