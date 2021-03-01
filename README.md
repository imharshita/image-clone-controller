# `Image-Clone-Controller`

The `Image Clone Controller` watches for Deployments/DaemonSets and checks 
Here Registry name is "backupregistry"

It leans heavily on the lower level
[`controller-runtime`](https://github.com/kubernetes-sigs/controller-runtime) package 
and [`remote`](https://github.com/google/go-containerregistry/tree/main/pkg/v1/remote) package

* Watch the Kubernetes Deployment and DaemonSet objects
* Check if any of them provision pods with images that are not from the backup
registry
* If yes, copy the image over to a corresponding repository and tag in the backup
registry
* Modify the Deployment/DaemonSet to use the image from the backup registry
* IMPORTANT: The Deployments and DaemonSets in the kube-system namespace
is ignored!


Registry secret yaml should be present

```
$ kubectl create -f conifg/secrets/secret.yaml
$ kubectl create -f conifg/manager/service_account.yaml
$ kubectl create -f conifg/manager/role.yaml
$ kubectl create -f conifg/manager/role_binding.yaml
```
### Goal
Goal here is to be safe against the risk of public container images disappearing from the registry while
we use them, breaking our deployments.

### 
Controller is written in Go





