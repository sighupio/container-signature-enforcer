package kubernetes.admission

strict_mode = true {
    data.webhook["opa-notary-connector-mode"]["mode.json"].strict
}

strict_deny_logic(c) = m {
    contains(c, "@sha256:")

    response := req_opa_notary_connector(c)
    response.status_code == 200
    response.body.image != c

    m := sprintf("Container image %v digest changed, new digest: %v", [c, response.body.image])
}

deny[msg] {
    strict_mode
    is_pod

    some j;
    container_image := input.request.object.spec.containers[j].image

    msg := strict_deny_logic(container_image)
}

deny[msg] {
    strict_mode
    is_cronjob

    some j;
    container_image := input.request.object.spec.jobTemplate.spec.template.spec.containers[j].image

    msg := strict_deny_logic(container_image)
}

deny[msg] {
    strict_mode
    not_pod_and_not_cronjob

    some j;
    container_image := input.request.object.spec.template.spec.containers[j].image

    msg := strict_deny_logic(container_image)
}
