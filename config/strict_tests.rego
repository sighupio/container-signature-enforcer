package kubernetes.admission

import data.kubernetes.admission.mocks

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as false
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as false
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
