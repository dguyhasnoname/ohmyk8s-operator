# Bootstrap operator

## Setup Go project for operator

Please use below commands to initialize a new project, setting up all the necessary base files and configuring the domain for the CRDs. Please make sure the directory is empty.

```bash
- mkdir ~/your-dir-where-you-wish-to-keep-operator-code
- cd ~/your-dir-where-you-wish-to-keep-operator-code/operator-01
- kubebuilder init --domain myoperator.io --repo github.com/dguyhasnoname/ohmyk8s-operator
```

Output

```
[10:46 AM IST 02.04.2024 ☸ 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 𖦥 main] 
 ┗━ ॐ  kubebuilder init --domain myoperator.io --repo github.com/dguyhasnoname/ohmyk8s-operator
INFO Writing kustomize manifests for you to edit... 
INFO Writing scaffold for you to edit...          
INFO Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.16.3 
INFO Update dependencies:
$ go mod tidy           
Next: define a resource with:
$ kubebuilder create api
```

Directory structure created:

```
[10:50 AM IST 02.04.2024 ☸ 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 𖦥 main] 
 ┗━ ॐ  tree ./
./
├── Dockerfile                                          # used for building a containerized version of the operator.
├── Makefile                                            # containing targets for building, testing, and deploying the operator
├── PROJECT                                             # YAML file containing project metadata
├── README.md
├── cmd
│   └── main.go                                         # entrypoint of the operator
├── config                                              # contains different config on how operator will be deployed
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   └── manager_config_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   └── rbac
│       ├── auth_proxy_client_clusterrole.yaml
│       ├── auth_proxy_role.yaml
│       ├── auth_proxy_role_binding.yaml
│       ├── auth_proxy_service.yaml
│       ├── kustomization.yaml
│       ├── leader_election_role.yaml
│       ├── leader_election_role_binding.yaml
│       ├── role.yaml
│       ├── role_binding.yaml
│       └── service_account.yaml
├── go.mod
├── go.sum
└── hack
    └── boilerplate.go.txt

8 directories, 25 files
```

## API creation

Please use below command to scaffold the necessary files for our operator API under the api and conrollers directories.

```bash
kubebuilder create api --group namespaceconfig --version v1 --kind Namespaceconfig --namespaced false
```

Output:

```
[10:50 AM IST 02.04.2024 ☸ 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 𖦥 main] 
 ┗━ ॐ  kubebuilder create api --group namespaceconfig --version v1 --kind Namespaceconfig
INFO Create Resource [y/n]                        
y
INFO Create Controller [y/n]                      
y
INFO Writing kustomize manifests for you to edit... 
INFO Writing scaffold for you to edit...          
INFO api/v1/namespaceconfig_types.go              
INFO api/v1/groupversion_info.go                  
INFO internal/controller/suite_test.go            
INFO internal/controller/namespaceconfig_controller.go 
INFO Update dependencies:
$ go mod tidy           
INFO Running make:
$ make generate                
mkdir -p /Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin
test -s /Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen && /Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen --version | grep -q v0.13.0 || \
	GOBIN=/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

Scaffolding created:

```
[10:51 AM IST 02.04.2024 ☸ 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 𖦥 main] 
 ┗━ ॐ  tree .
.
├── Dockerfile
├── Makefile
├── PROJECT
├── README.md
├── api                                                     # contain the API definitions for CRDs
│   └── v1
│       ├── groupversion_info.go                            # GVK of the operator
│       ├── namespaceconfig_types.go                        # spec definition for the CRD object             
│       └── zz_generated.deepcopy.go                        # deepcopy functions
├── bin
│   └── controller-gen                                      # binary used for controller build
├── cmd
│   └── main.go
├── config
│   ├── crd
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   └── manager_config_patch.yaml                       # patch to make changes in controller deployment
│   ├── manager                                             # controller deployment spec
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac                                                 # RBAC definitions to be used by controller
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   ├── namespaceconfig_editor_role.yaml
│   │   ├── namespaceconfig_viewer_role.yaml
│   │   ├── role.yaml
│   │   ├── role_binding.yaml
│   │   └── service_account.yaml
│   └── samples                                                # sample namespaceconfigs 
│       ├── kustomization.yaml
│       └── namespaceconfig_v1_namespaceconfig.yaml            
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── internal                                                   # contain the controller implementations. Logic of the operator
    └── controller
        ├── namespaceconfig_controller.go
        └── suite_test.go                                      # testing setup of controller logic

15 directories, 37 files
```

## Define API struct

API Struct is defined in file namespaceconfig_types.go

## Deploy the CRD

Run `make manifests` command to generate CRD manifests, RBAC manifests, and webhook manifests etc.:

```
[12:19 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  make manifests
/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
github.com/dguyhasnoname/ohmyk8s-operator/api/v1:-: use of unimported package "v1"
github.com/dguyhasnoname/ohmyk8s-operator/api/v1:-: use of unimported package "v1"
/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

Run `make install` to apply the generated manifest to test Kubernetes cluster:

```
[12:21 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  make install
/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
test -s /Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/kustomize || GOBIN=/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@v5.2.1
/Users/Mukund_Bihari/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/namespaceconfigs.namespaceconfig.myoperator.io created
[12:21 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  kg crd
NAME                                             CREATED AT
namespaceconfigs.namespaceconfig.myoperator.io   2024-04-05T06:51:43Z
[12:21 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
```