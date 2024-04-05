/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	namespaceconfigv1 "github.com/dguyhasnoname/ohmyk8s-operator/api/v1"
	"github.com/dguyhasnoname/ohmyk8s-operator/pkg/util"
)

// NamespaceconfigReconciler reconciles a Namespaceconfig object
type NamespaceconfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=namespaceconfig.myoperator.io,resources=namespaceconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=namespaceconfig.myoperator.io,resources=namespaceconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=namespaceconfig.myoperator.io,resources=namespaceconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespaceconfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NamespaceconfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := util.Logs
	log.Info("Starting namespace config reconcilliation")

	o := &namespaceconfigv1.Namespaceconfig{}
	err := r.Get(ctx, req.NamespacedName, o)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error("Namespaceconfig ", o.GetName(), " not found.")
		}
	}
	namespaceName := o.Spec.Abbreviation + "-" + o.Spec.Environment
	finalizerName := "namespaceconfig.myoperator.io/finalizer"
	// check if the object is being deleted
	if o.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Debug("DeletionTimestamp is zero for Namespaceconfig ", o.GetName())
		if !controllerutil.ContainsFinalizer(o, finalizerName) {
			log.Info("Finalizer not found, adding finalizer ", finalizerName, " to Namespaceconfig ", o.GetName())
			controllerutil.AddFinalizer(o, finalizerName)
			if err := r.Update(ctx, o); err != nil {
				return ctrl.Result{}, err
			} else {
				log.Info("Finalizer added to Namespaceconfig ", o.GetName())
			}
		}
		// Attempt to create the namespace
		if o.Status.Status == "" {
			log.Info("Creating Namespace ", namespaceName)
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
					Annotations: map[string]string{
						"managed-by": "namespaceconfig.myoperator.io",
					},
					Labels: map[string]string{
						"owner": o.Spec.NamespaceOwner,
						"env":   o.Spec.Environment,
					},
				},
			}
			if err = r.Create(ctx, namespace); err != nil {
				return ctrl.Result{}, err
			} else {
				log.Info("Namespace ", namespaceName, " created.")
				log.Info("Setting ownerReference for namespace object ", namespaceName)
				if err := ctrl.SetControllerReference(o, namespace, r.Scheme); err != nil {
					log.Error(err, "Unable to set ownerReference for namespace object", namespaceName)
				} else {
					log.Info("OwnerReference set for namespace", namespaceName)
				}
				o.Status.NamespaceName = namespaceName
				o.Status.Status = "RUNNING"
				o.Status.LastUpdate = metav1.Now().String()
				if err := r.Status().Update(ctx, o); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("Namespaceconfig ", o.GetName(), " status updated.")
				return ctrl.Result{}, nil
			}
		} else {
			log.Info("Namespace ", namespaceName, " already exists.")
			return ctrl.Result{}, nil
		}
	} else {
		log.Info("DeletionTimestamp is not zero for Namespaceconfig ", o.GetName())
		if controllerutil.ContainsFinalizer(o, finalizerName) {
			log.Info("Finalizer found, removing finalizer ", finalizerName, " from Namespaceconfig ", o.GetName())
			// Remove the finalizer from the Namespaceconfig object once the cleanup succeeded
			// This will free up Namespaceconfig resource to be deleted
			controllerutil.RemoveFinalizer(o, finalizerName)
			if err := r.Update(ctx, o); err != nil {
				log.Error(err, "Unable to remove finalizer and update Namespaceconfig")
				return ctrl.Result{}, err
			}
		} else {
			log.Warn("Finalizer not found for Namespaceconfig ", o.GetName(), " marked for deletion.")
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceconfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&namespaceconfigv1.Namespaceconfig{}).
		Complete(r)
}
