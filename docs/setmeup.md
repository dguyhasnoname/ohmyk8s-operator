# Bootstrap development env

## Pre-requisites

- install Go
  - version 1.20 or above
- container runtime:
  - docker[preferred]
  - podman
- kubernetes cluster
  - kind
  - minikube
  - k3s
- kubectl binary

### Installing Go

Please follow the steps mention in official Go documentation: [Download and install Go quickly](https://go.dev/doc/install)

### Installing container runtime

Please follow below steps:

- docker: [get-docker](https://docs.docker.com/get-docker/)
- podman: [Podman Installation Instructions](https://podman.io/docs/installation)

### Installing Kubernetes cluster

Please follow below steps:

- kind: [quick-start](https://kind.sigs.k8s.io/docs/user/quick-start/)
- minikube: [minikube-start](https://minikube.sigs.k8s.io/docs/start/)
- k3s: [installation](https://docs.k3s.io/installation)

```
[12:19 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  minikube start
😄  minikube v1.32.0 on Darwin 14.4.1 (arm64)
✨  Using the docker driver based on existing profile
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
🔄  Restarting existing docker container for "minikube" ...
🐳  Preparing Kubernetes v1.28.3 on Docker 24.0.7 ...
🔗  Configuring bridge CNI (Container Networking Interface) ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
[12:19 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  kg no
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   76d   v1.28.3
```

### Installing Kubebuilder

Please execute below commands:

```bash
$ curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
$ chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```
