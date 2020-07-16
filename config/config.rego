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

req_opa_notary_connector(s) = x {
    request := {
        "url": "http://localhost:8080/checkImage",
        "method": "POST",
        "body": {
            "image": s,
        }
    }
    x := http.send(request)
}

is_pod(k) = x {
    k != "CronJob"
    k != "Deployment"
    k != "Daemonset"
    k != "Job"
    k != "ReplicationController"
    k != "StatefulSet"
    k != "ReplicaSet"
    x := true
}

is_cronjob(k) = x {
    k != "Pod"
    k != "Deployment"
    k != "Daemonset"
    k != "Job"
    k != "ReplicationController"
    k != "StatefulSet"
    k != "ReplicaSet"
    x := true
}

not_pod_and_not_cronjob(k) = x {
    k != "Pod"
    k != "CronJob"
    x := true
}

deny[msg] {
    is_pod(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code != 200
    error_message := response.body.error

    msg := sprintf("Container image %v invalid: %v", [container_image, error_message])
}

deny[msg] {
    is_cronjob(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.jobTemplate.spec.template.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code != 200
    error_message := response.body.error

    msg := sprintf("Container image %v invalid: %v", [container_image, error_message])
}

deny[msg] {
    not_pod_and_not_cronjob(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.template.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code != 200
    error_message := response.body.error

    msg := sprintf("Container image %v invalid: %v", [container_image, error_message])
}

patches["pod_sha"] = patch {
    is_pod(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code == 200
    new_container_image := response.body.image

    patch := gen_patch(input.request.kind.kind, j, new_container_image)
}

patches["cronjob_sha"] = patch {
    is_cronjob(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.jobTemplate.spec.template.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code == 200
    new_container_image := response.body.image

    patch := gen_patch(input.request.kind.kind, j, new_container_image)
}

patches["others_sha"] = patch {
    not_pod_and_not_cronjob(input.request.kind.kind)

    some j;
    container_image = input.request.object.spec.template.spec.containers[j].image

    response := req_opa_notary_connector(container_image)
    response.status_code == 200
    new_container_image := response.body.image

    patch := gen_patch(input.request.kind.kind, j, new_container_image)
}
