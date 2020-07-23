package kubernetes.admission

import data.kubernetes.admission.mocks


test_not_deny_minio {
    not denied with input as mocks.alpine_3_10_minio_tenant
}

test_deny_minio {
    denied with input as mocks.alpine_3_11_minio_tenant
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
    not denied with input as mocks.alpine_3_10_minio_tenant_correct_sha
}

test_not_deny_minio_incorrect_sha {
    not denied with input as mocks.alpine_3_10_minio_tenant_incorrect_sha
}

test_not_deny_minio_correct_sha_strict {
    not denied with input as mocks.alpine_3_10_minio_tenant_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_deny_minio_incorrect_sha_strict {
    denied with input as mocks.alpine_3_10_minio_tenant_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}
