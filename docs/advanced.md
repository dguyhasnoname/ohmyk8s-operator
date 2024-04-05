# Advanced configs

## Remote watching of resources

Go library: https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/manager#Options

```
mgr, err := ctrl.NewManager(<remote-k8s-client>, ctrl.Options{
    Scheme:                 scheme,
    Metrics:                metricsserver.Options{BindAddress: metricsAddr},
    HealthProbeBindAddress: probeAddr,
    LeaderElection:         enableLeaderElection,
    LeaderElectionID:       leaderElectionID,
    LeaderElectionConfig:   <another-k8s-client can be used>,
})
```

### From where the reconcile func is called?

```
if err = (&controller.NSAssetReconciler{                                         # creates an instance of the controller.NSAssetReconciler struct. The & operator is used to take the address of the struct, which is often required when passing the struct to a function that expects a pointer.
    Client: mgr.GetClient(),                                                     # retrieves a Kubernetes client from the controller manager (mgr) and assigns it to the Client field of the NSAssetReconciler struct. This client will be used by the controller to interact with the Kubernetes API server.
    Scheme: mgr.GetScheme(),                                                     # retrieves the Kubernetes runtime scheme from the controller manager (mgr) and assigns it to the Scheme field of the NSAssetReconciler struct. The scheme is used by the controller to understand the types of Kubernetes objects it manages.
}).SetupWithManager(mgr); err != nil {                                           # calls the SetupWithManager function on the newly created NSAssetReconciler instance. This function is responsible for registering the controller with the controller manager (mgr). 
    setupLog.Error(err, "unable to create controller", "controller", "NSAsset")
}
```

```
return ctrl.NewControllerManagedBy(mgr).                # creates a new controller object. The mgr argument  reference to the controller manager instance in the project, which is responsible for running all the controllers.
    For(&namespaceconfigv1.Namespaceconfig{}).          # specifies the Kubernetes resource type that the controller will watch for changes.
    Owns(&corev1.Namespace{}, predicateNamespace).      # defines the relationship b/w NamespaceConfig and ns resources. It means that  controller will manage ns resources that are "owned" by the NamespaceConfig resources it watches.
    Complete(r)                                         # completes the controller definition by passing the controller object (r) to be managed by the controller manager.
```

### What are predicates?

Predicates are used by Controllers to filter Events before they are provided to EventHandlers. Works with [builder](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/builder) pkg.

```
predicateNamespace := builder.WithPredicates(predicate.Funcs{
    CreateFunc: func(createEvent event.CreateEvent) bool {
        return true
    },
    UpdateFunc: func(updateEvent event.UpdateEvent) bool {
        oldNS := updateEvent.ObjectOld.(*corev1.Namespace)
        newNS := updateEvent.ObjectNew.(*corev1.Namespace)
        same := reflect.DeepEqual(oldNS.ObjectMeta.Labels, newNS.ObjectMeta.Labels)
        same = same && oldNS.ObjectMeta.Name == newNS.ObjectMeta.Name
        return !same
    },
    DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
        return true
    },
    GenericFunc: func(genericEvent event.GenericEvent) bool {
        return true
    },
	})
```

### How logging is controlled?

[logging](../operator-01/pkg/logs/log.go)

### How testing can be done?

[testing](../operator-01/internal/controller/suite_test.go)


Sample test suite execution:

```
[09:59 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  make test
test -s /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen && /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen --version | grep -q v0.13.0 || \
	GOBIN=/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
test -s /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/setup-envtest || GOBIN=/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
go: downloading sigs.k8s.io/controller-runtime v0.17.2
go: downloading sigs.k8s.io/controller-runtime/tools/setup-envtest v0.0.0-20240405052210-76d3d0826fa9
go: sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20240405052210-76d3d0826fa9 requires go >= 1.22.0; switching to go1.22.2
go: downloading go1.22.2 (darwin/arm64)
KUBEBUILDER_ASSETS="/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/k8s/1.28.0-darwin-arm64" go test ./... -coverprofile cover.out
?   	github.com/dguyhasnoname/ohmyk8s-operator/api/v1	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/cmd	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/pkg/logs	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/pkg/util	[no test files]
ok  	github.com/dguyhasnoname/ohmyk8s-operator/internal/controller	0.040s	coverage: 0.0% of statements
[10:01 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 

```

Testing with a sample namespaceconfig:

```
[10:23 PM IST 05.04.2024 ☸ minikube 📁 ~/git/dguyhasnoname/ohmyk8s-operator/operator-01 ❱ main ▲] 
 ┗━ ॐ  make test
test -s /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen && /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen --version | grep -q v0.13.0 || \
	GOBIN=/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
KUBEBUILDER_ASSETS="/Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/bin/k8s/1.28.0-darwin-arm64" go test -v ./... -coverprofile cover.out
?   	github.com/dguyhasnoname/ohmyk8s-operator/api/v1	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/cmd	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/pkg/logs	[no test files]
?   	github.com/dguyhasnoname/ohmyk8s-operator/pkg/util	[no test files]
{"level":"INFO","timestamp":"2024-04-05T22:24:22.364+0530","caller":"util/util.go:19","log":"Initializing env..."}
=== RUN   TestControllers
Running Suite: Controller Suite - /Users/mk/git/dguyhasnoname/ohmyk8s-operator/operator-01/internal/controller
=========================================================================================================================
Random Seed: 1712336062

Will run 1 of 1 specs
Running test on minikube  cluster.
•

Ran 1 of 1 Specs in 3.063 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestControllers (3.06s)
PASS
coverage: 0.0% of statements
ok  	github.com/dguyhasnoname/ohmyk8s-operator/internal/controller	3.095s	coverage: 0.0% of statements
```