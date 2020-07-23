package kubernetes.admission

import data.kubernetes.admission.mocks


test_not_denied {
  not denied with deny as set()
}

test_denied {
  denied with deny as {"test"}
}

test_image_pod {
    img := images with input as mocks.alpine_3_11_pod
    count(img) == 1
}

test_image_pod {
    img := images with input as mocks.alpine_3_10_pod
    count(img) == 1
}

test_different_images_pod {
    img := images with input as mocks.alpine_3_11_and_3_10_pod
    count(img) == 2
}

test_same_images_pod {
    img := images with input as mocks.alpine_3_10_and_3_10_pod
    count(img) == 2
}

test_not_deny_pod {
  not denied with input as mocks.alpine_3_10_pod
}

test_not_deny_pod {
  not denied with input as mocks.alpine_3_10_and_3_10_pod
}

test_deny_pod {
  denied with input as mocks.alpine_3_11_pod
}

test_deny_pod {
  denied with input as mocks.alpine_3_11_and_3_10_pod
}


test_deny_deployment {
    denied with input as mocks.alpine_3_11_deployment
}

test_deny_deployment {
    denied with input as mocks.alpine_3_10_and_3_11_deployment
}

test_deny_cronjob {
  denied with input as mocks.alpine_3_11_cronjob
}

test_deny_cronjob {
  denied with input as mocks.alpine_3_10_and_3_11_cronjob
}

test_not_deny_cronjob {
  not denied with input as mocks.alpine_3_10_cronjob
}

test_not_deny_cronjob {
  not denied with input as mocks.alpine_3_10_and_3_10_cronjob
}

test_not_deny_deployment {
  not denied with input as mocks.alpine_3_10_deployment
}

test_not_deny_deployment {
  not denied with input as mocks.alpine_3_10_and_3_10_deployment
}

