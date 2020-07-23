package kubernetes.admission

image = {
    "Pod": "/spec/containers/%v/image",

    "CronJob": "/spec/jobTemplate/spec/template/spec/containers/%v/image",

    "Deployment": "/spec/template/spec/containers/%v/image",
    "Daemonset": "/spec/template/spec/containers/%v/image",
    "Job": "/spec/template/spec/containers/%v/image",
    "ReplicationController": "/spec/template/spec/containers/%v/image",
    "StatefulSet": "/spec/template/spec/containers/%v/image",
    "ReplicaSet": "/spec/template/spec/containers/%v/image"
}

gen_patch(k, i, c) = p {
    p := [{"op": "replace", "path": sprintf(image[k], [i]), "value": c}]
}

annotation_patch(metadata) = p {
    not metadata.annotations
    p := [{"op": "add", "path": "/metadata/annotations", "value": {"opa-notary-connector.sighup.io/processed": "true"}}]
}

annotation_patch(metadata) = p {
    metadata.annotations
    p := [{"op": "add", "path": "/metadata/annotations/opa-notary-connector.sighup.io~1processed", "value": "true"}]
}

req_opa_notary_connector(s) = x {
    request := {
        "url": "http://localhost:8080/checkImage",
        "method": "POST",
        "headers": {
            "X-Request-ID": uuid.rfc4122("")
        },
        "body": {
            "image": s,
        }
    }
    x := http.send(request)
}

is_pod {
    input.request.kind.kind == "Pod"
}

is_cronjob {
    input.request.kind.kind == "CronJob"
}

not_pod_and_not_cronjob {
    input.request.kind.kind != "Pod"
    input.request.kind.kind != "CronJob"
}

deny_logic(c) = m {
    response := req_opa_notary_connector(c)
    response.status_code != 200
    error_message := response.body.error
    m := sprintf("Container image %v invalid: %v", [c, error_message])
}

patch_logic(k, i, c) = p {
    response := req_opa_notary_connector(c)
    response.status_code == 200
    new_container_image := response.body.image
    p := gen_patch(k, i, new_container_image)
}

prepare_patch(request, index, container_image) = p {
    c_patch := patch_logic(request.kind.kind, index, container_image)
    a_patch := annotation_patch(request.object.metadata)
    p = array.concat(c_patch, a_patch)
}

deny[msg] {
    is_pod

    some j;
    container_image := input.request.object.spec.containers[j].image

    msg := deny_logic(container_image)
}

deny[msg] {
    is_cronjob

    some j;
    container_image := input.request.object.spec.jobTemplate.spec.template.spec.containers[j].image

    msg := deny_logic(container_image)
}

deny[msg] {
    not_pod_and_not_cronjob

    some j;
    container_image := input.request.object.spec.template.spec.containers[j].image

    msg := deny_logic(container_image)
}

patches[patch] {
    is_pod

    some j;
    container_image := input.request.object.spec.containers[j].image

    patch := prepare_patch(input.request, j, container_image)
}

patches[patch] {
    is_cronjob

    some j;
    container_image := input.request.object.spec.jobTemplate.spec.template.spec.containers[j].image

    patch := prepare_patch(input.request, j, container_image)
}

patches[patch] {
    not_pod_and_not_cronjob

    some j;
    container_image := input.request.object.spec.template.spec.containers[j].image

    patch := prepare_patch(input.request, j, container_image)
}
