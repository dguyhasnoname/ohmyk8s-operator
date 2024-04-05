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
if err = (&controller.NSAssetReconciler{
    Client: mgr.GetClient(),
    Scheme: mgr.GetScheme(),
}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "unable to create controller", "controller", "NSAsset")
}
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
