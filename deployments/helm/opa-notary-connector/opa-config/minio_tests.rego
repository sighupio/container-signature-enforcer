package kubernetes.admission

import data.kubernetes.admission.mocks


test_not_deny_minio {
    msg := deny with input as mocks.alpine_3_10_minio_tenant
    count(msg) == 0
}

test_deny_minio {
    msg := deny with input as mocks.alpine_3_11_minio_tenant
    count(msg) > 0
}

test_patch_minio {
    patch := patches with input as mocks.alpine_3_10_minio_tenant
    count(patch) == 2
}

test_patch_minio {
    patch := patches with input as mocks.alpine_3_11_minio_tenant
    count(patch) == 1
}

test_not_deny_minio_correct_sha {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_correct_sha
    count(msg) == 0
}

test_not_deny_minio_incorrect_sha {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_incorrect_sha
    count(msg) == 0
}

test_not_deny_minio_correct_sha_strict {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 0
}

test_deny_minio_incorrect_sha_strict {
    msg := deny with input as mocks.alpine_3_10_minio_tenant_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
    count(msg) == 1
}
