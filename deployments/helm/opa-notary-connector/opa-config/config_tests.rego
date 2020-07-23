package kubernetes.admission

import data.kubernetes.admission.mocks


test_images {
    img := images with input as mocks.alpine_3_11_pod
    count(img) == 1
}

test_images {
    img := images with input as mocks.alpine_3_11_and_3_10_pod
    count(img) == 2
}

test_images {
    img := images with input as mocks.alpine_3_10_and_3_10_pod
    count(img) == 2
}

test_not_deny_pod {
  msg := deny with input as mocks.alpine_3_10_pod
  count(msg) == 0
}

test_not_deny_pod {
  msg := deny with input as mocks.alpine_3_10_and_3_10_pod
  count(msg) == 0
}

test_deny_pod {
  deny with input as mocks.alpine_3_11_pod
}

test_deny_pod {
  deny with input as mocks.alpine_3_11_and_3_10_pod
}


test_deny_deployment {
    deny with input as mocks.alpine_3_11_deployment
}

test_deny_deployment {
    deny with input as mocks.alpine_3_10_and_3_11_deployment
}

test_deny_cronjob {
  deny with input as mocks.alpine_3_11_cronjob
}

test_deny_cronjob {
  deny with input as mocks.alpine_3_10_and_3_11_cronjob
}

test_not_deny_cronjob {
  msg := deny with input as mocks.alpine_3_10_cronjob
  count(msg) == 0
}

test_not_deny_cronjob {
  msg := deny with input as mocks.alpine_3_10_and_3_10_cronjob
  count(msg) == 0
}

test_not_deny_deployment {
  msg := deny with input as mocks.alpine_3_10_deployment
  count(msg) == 0
}

test_not_deny_deployment {
  msg := deny with input as mocks.alpine_3_10_and_3_10_deployment
  count(msg) == 0
}

