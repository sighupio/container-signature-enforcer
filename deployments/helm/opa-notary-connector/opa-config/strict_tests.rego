package kubernetes.admission

import data.kubernetes.admission.mocks

test_strict_mode {
  not strict_mode
}

test_strict_mode {
  strict_mode with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_deny_pod_strict {
  denied with input as mocks.alpine_3_10_pod_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_deny_cronjob_strict {
  denied with input as mocks.alpine_3_10_cronjob_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_deny_deployment_strict {
  denied with input as mocks.alpine_3_10_deployment_incorrect_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_not_deny_pod_strict {
  not denied with input as mocks.alpine_3_10_pod_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_not_deny_deployment_strict {
  not denied with input as mocks.alpine_3_10_deployment_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_not_deny_cronjob_strict {
  not denied with input as mocks.alpine_3_10_cronjob_correct_sha with data.webhook["opa-notary-connector-mode"]["mode.json"].strict as true
}

test_not_deny_cronjob {
  not denied with input as mocks.alpine_3_10_cronjob_correct_sha
}

test_not_deny_cronjob {
  not denied with input as mocks.alpine_3_10_cronjob_incorrect_sha
}

test_not_deny_deployment {
  not denied with input as mocks.alpine_3_10_deployment_correct_sha
}

test_not_deny_deployment {
  not denied with input as mocks.alpine_3_10_deployment_incorrect_sha
}


test_not_deny_pod {
  not denied with input as mocks.alpine_3_10_pod_correct_sha
}

test_not_deny_pod {
  not denied with input as mocks.alpine_3_10_pod_incorrect_sha
}

