package kubernetes.admission

images[img] {
  input.request.kind.kind == "Tenant"
  img := {
    "patch_path": "/spec/image",
    "image": input.request.object.spec.image,
  }
}
