package core

//const (
//podSpecPath      = "/spec"
//templateSpecPath = "/spec/template/spec"
//cronJobSpecPath  = "/spec/jobTemplate/spec/template/spec"
//)

//TODO port to rego
//// Extracts the PodSpec field from the accepted resources, returning also the path to the spec, to be used to write patches
//func getPodSpec(ar *vAdmission.AdmissionReview, log *logrus.Entry) (baseSpecPath string, podSpec *corev1.PodSpec, err error) {

//podSpec = &corev1.PodSpec{}

//// 2. extract the pod object from the admission request, only request having a pod as object will arrive
//// Switch between different objects and get pod or podSpec in it
//log.WithField("resource", ar.Request.Resource).Debug("unwrapping resource to correct object")
//switch ar.Request.Resource {
//case metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}:
//pod := corev1.Pod{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &pod); err != nil {
//return "", nil, err
//}
//podSpec = &pod.Spec
//baseSpecPath = podSpecPath

//case metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"}:
//rc := corev1.ReplicationController{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &rc); err != nil {
//return "", nil, err
//}
//podSpec = &rc.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"}:
//deploy := extensionsv1beta1.Deployment{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &deploy); err != nil {
//return "", nil, err
//}
//podSpec = &deploy.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "deployments"}:
//deploy := appsv1beta1.Deployment{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &deploy); err != nil {
//return "", nil, err
//}
//podSpec = &deploy.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "deployments"}:
//deploy := appsv1beta2.Deployment{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &deploy); err != nil {
//return "", nil, err
//}
//podSpec = &deploy.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}:
//deploy := appsv1.Deployment{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &deploy); err != nil {
//return "", nil, err
//}
//podSpec = &deploy.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}:
//rs := appsv1.ReplicaSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &rs); err != nil {
//return "", nil, err
//}
//podSpec = &rs.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"}:
//rs := extensionsv1beta1.ReplicaSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &rs); err != nil {
//return "", nil, err
//}
//podSpec = &rs.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "replicasets"}:
//rs := appsv1beta2.ReplicaSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &rs); err != nil {
//return "", nil, err
//}
//podSpec = &rs.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}:
//ds := appsv1.DaemonSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &ds); err != nil {
//return "", nil, err
//}
//podSpec = &ds.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "daemonsets"}:
//ds := extensionsv1beta1.DaemonSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &ds); err != nil {
//return "", nil, err
//}
//podSpec = &ds.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "daemonsets"}:
//ds := appsv1beta2.DaemonSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &ds); err != nil {
//return "", nil, err
//}
//podSpec = &ds.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}:
//sts := appsv1.StatefulSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &sts); err != nil {
//return "", nil, err
//}
//podSpec = &sts.Spec.Template.Spec
//baseSpecPath = templateSpecPath
//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "statefulsets"}:
//sts := appsv1beta1.StatefulSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &sts); err != nil {
//return "", nil, err
//}
//podSpec = &sts.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "statefulsets"}:
//sts := appsv1beta2.StatefulSet{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &sts); err != nil {
//return "", nil, err
//}
//podSpec = &sts.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}:
//job := batchv1.Job{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &job); err != nil {
//return "", nil, err
//}
//podSpec = &job.Spec.Template.Spec
//baseSpecPath = templateSpecPath

//case metav1.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"}:
//job := batchv1beta1.CronJob{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &job); err != nil {
//return "", nil, err
//}
//podSpec = &job.Spec.JobTemplate.Spec.Template.Spec
//baseSpecPath = cronJobSpecPath

//case metav1.GroupVersionResource{Group: "batch", Version: "v2alpha1", Resource: "cronjobs"}:
//job := batchv2alpha1.CronJob{}
//if err := json.Unmarshal(ar.Request.Object.Raw, &job); err != nil {
//return "", nil, err
//}
//podSpec = &job.Spec.JobTemplate.Spec.Template.Spec
//baseSpecPath = cronJobSpecPath
//default:
//log.WithField("request", ar).Error("Resource not supported")
//return "", nil, fmt.Errorf(`The resource "%s/%s/%s" is not supported. Make sure that you are using a supported kubectl version, and that you are using a supported Kubernetes workload type`, ar.Request.Resource.Group, ar.Request.Resource.Version, ar.Request.Resource.Resource)
//}

//log.WithFields(logrus.Fields{"podSpec": podSpec, "baseSpecPath": baseSpecPath, "request": ar}).Debug("Extracted podspec")
//return baseSpecPath, podSpec, nil

//}
