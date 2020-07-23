package kubernetes.admission


images[img] {
  input.request.kind.kind == "Pod"
  img := {
    "patch_path": "/spec/containers/%v/image",
    "image": input.request.object.spec.containers[i].image,
    "index": i
  }
}

images[img] {
  input.request.kind.kind == "CronJob"
  img := {
    "patch_path": "/spec/jobTemplate/spec/template/spec/containers/%v/image",
    "image": input.request.object.spec.jobTemplate.spec.template.spec.containers[i].image,
    "index": i
  }
}

images[img] {
  known_resources := {
    "Deployment",
    "Daemonset",
    "Job",
    "ReplicationController",
    "StatefulSet",
    "ReplicaSet"
  }
  known_resources[input.request.kind.kind]
  img := {
    "patch_path": "/spec/template/spec/containers/%v/image",
    "image": input.request.object.spec.template.spec.containers[i].image,
    "index": i
    }
}

opa_notary_connector_responses[resp] {
  images[i].index
  resp := {
    "response": req_opa_notary_connector(images[i].image),
    "image" : images[i].image,
    "patch_path" : images[i].patch_path,
    "index": images[i].index
    }
} {
  resp := {
    "response": req_opa_notary_connector(images[i].image),
    "image" : images[i].image,
    "patch_path" : images[i].patch_path,
    }
}

deny[msg] {
  opa_notary_connector_responses[i].response.status_code != 200
  error_message := opa_notary_connector_responses[i].response.body.error
  msg := sprintf("Container image %v invalid: %v", [opa_notary_connector_responses[i].image, error_message])
}

denied = true {
  count(deny)>0
}


patches[patch] {
  opa_notary_connector_responses[i].index
  opa_notary_connector_responses[i].response.status_code == 200
  patch := {"op": "replace", "path": sprintf(opa_notary_connector_responses[i].patch_path, [opa_notary_connector_responses[i].index]) , "value":opa_notary_connector_responses[i].image }
} {
  opa_notary_connector_responses[i].response.status_code == 200
  patch := {"op": "replace", "path": opa_notary_connector_responses[i].patch_path , "value":opa_notary_connector_responses[i].image }
}

patches[patch]{
  input.request.object.metadata.annotations
  patch := {"op": "add", "path": "/metadata/annotations/opa-notary-connector.sighup.io~1processed", "value": "true"}
} {
  patch := {"op": "add", "path": "/metadata/annotations", "value": {"opa-notary-connector.sighup.io/processed": "true"}}
}

patched = true {
  count(patches) > 1
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

