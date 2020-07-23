package kubernetes.admission

import data.kubernetes.admission.mocks

test_strict_mode {
    not strict_mode
}

test_strict_mode {
    strict_mode with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_strict_deny_logic {
    not strict_deny_logic("alpine:3.10")
}

test_strict_deny_logic {
    not strict_deny_logic("alpine:3.10")
}

test_strict_deny_logic {
    msg := strict_deny_logic("alpine@sha256:malware")
    msg == "Container image alpine@sha256:malware digest changed, new digest: alpine@sha256:randomsha"
}

test_strict_deny_logic {
    not strict_deny_logic("alpine@sha256:randomsha")
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_correct_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_incorrect_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 1
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_cronjob_correct_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_cronjob_incorrect_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_cronjob_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_cronjob_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 1
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_deployment_correct_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_deployment_incorrect_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_deployment_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_deployment_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 1
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_correct_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_incorrect_sha
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 1
}
