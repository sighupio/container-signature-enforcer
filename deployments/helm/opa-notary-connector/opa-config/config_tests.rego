package kubernetes.admission

import data.kubernetes.admission.mocks

test_is_pod {
    is_pod with input as {"request": {"kind": {"kind": "Pod"}}}
}

test_is_cronjob {
    is_cronjob with input as {"request": {"kind": {"kind": "CronJob"}}}
}

test_not_pod_and_not_cronjob {
    not_pod_and_not_cronjob with input as {"request": {"kind": {"kind": "Deployment"}}}
}

test_not_pod_and_not_cronjob {
    not not_pod_and_not_cronjob with input as {"request": {"kind": {"kind": "Pod"}}}
}

test_not_pod_and_not_cronjob {
    not not_pod_and_not_cronjob with input as {"request": {"kind": {"kind": "CronJob"}}}
}

test_gen_patch {
    patch := gen_patch("Pod", 1, "alpine:3.10")
    patch == [{"op": "replace", "path": "/spec/containers/1/image", "value": "alpine:3.10"}]
}

test_gen_patch {
    patch := gen_patch("Deployment", 1, "alpine:3.10")
    patch == [{"op": "replace", "path": "/spec/template/spec/containers/1/image", "value": "alpine:3.10"}]
}

test_gen_patch {
    patch := gen_patch("CronJob", 2, "alpine:3.11")
    patch == [{"op": "replace", "path": "/spec/jobTemplate/spec/template/spec/containers/2/image", "value": "alpine:3.11"}]
}

test_annotation_patch {
    metadata := {
        "name": "hello",
        "labels": {
            "run": "hello"
        }
    }
    patch := annotation_patch(metadata)
    patch == [{"op": "add", "path": "/metadata/annotations", "value": {"opa-notary-connector.sighup.io/processed": "true"}}]
}

test_annotation_patch {
    metadata := {
        "name": "hello",
        "labels": {
            "run": "hello"
        },
        "annotations": {
            "run": "hello"
        }
    }
    patch := annotation_patch(metadata)
    patch == [{"op": "add", "path": "/metadata/annotations/opa-notary-connector.sighup.io~1processed", "value": "true"}]
}

test_req_opa_notary_connector {
    r := req_opa_notary_connector("alpine:3.10")
    new_container_image := r.body.image
    new_container_image == "alpine@sha256:randomsha"
}

test_req_opa_notary_connector {
    r := req_opa_notary_connector("alpine:3.11")
    error := r.body.error
    error == "image not found"
}

test_deny_logic {
    not deny_logic("alpine:3.10")
}

test_deny_logic {
    msg := deny_logic("alpine:3.11")
    msg == "Container image alpine:3.11 invalid: image not found"
}

test_patch_logic {
    not patch_logic("Pod", "1", "alpine:3.11")
}

test_patch_logic {
    p := patch_logic("Pod", "1", "alpine:3.10")
    p == [{"op": "replace", "path": "/spec/containers/1/image", "value": "alpine@sha256:randomsha"}]
}

test_prepare_patch {
    p := prepare_patch(input.request, "1", "alpine:3.10") with input as mocks.alpine_3_10_pod
    p == [{"op": "replace", "path": "/spec/containers/1/image", "value": "alpine@sha256:randomsha"}, {"op": "add", "path": "/metadata/annotations", "value": {"opa-notary-connector.sighup.io/processed": "true"}}]
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_pod
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_and_3_10_pod
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_11_pod
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}

test_deny {
    msg := deny with input as mocks.alpine_3_11_and_3_10_pod
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}

test_patch {
    patch := patches with input as mocks.alpine_3_11_and_3_10_pod
    count(patch) == 1
}

test_patch {
    patch := patches with input as mocks.alpine_3_10_and_3_10_pod
    count(patch) == 2
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_cronjob
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_and_3_10_cronjob
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_11_cronjob
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_and_3_11_cronjob
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_deployment
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_and_3_10_deployment
    count(msg) == 0
}

test_deny {
    msg := deny with input as mocks.alpine_3_11_deployment
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}

test_deny {
    msg := deny with input as mocks.alpine_3_10_and_3_11_deployment
    contains(msg[_], "Container image alpine:3.11 invalid: image not found")
}
