package kubernetes.admission

strict_mode = true {
    data.webhook["opa-notary-connector-mode"]["mode.json"].strict
}

deny[msg] {
    strict_mode
    opa_notary_connector_responses[i].response.status_code == 200
    original_img := opa_notary_connector_responses[i].image
    returned_img := opa_notary_connector_responses[i].response.body.image

    contains(original_img, "@sha256:")
    original_img != opa_notary_connector_responses[i].response.body.image

    msg := sprintf("Container image %v digest changed, new digest: %v", [original_img, opa_notary_connector_responses[i].response.body.image])

}
