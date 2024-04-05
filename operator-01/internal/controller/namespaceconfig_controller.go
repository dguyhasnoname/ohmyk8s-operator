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
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

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
			log.Info("Namespaceconfig not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		} else {
			log.Error("Failed to get Namespaceconfig: ", err)
		}
		return ctrl.Result{}, err
	}
	namespaceName := o.Spec.Abbreviation + "-" + o.Spec.Environment
	finalizerName := "namespaceconfig.myoperator.io/finalizer"
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
		err = r.Get(ctx, client.ObjectKey{Name: namespaceName}, namespace)
		if err != nil {
			log.Info("Creating Namespace ", namespaceName)
			if err = r.Create(ctx, namespace); err != nil {
				return ctrl.Result{}, err
			} else {
				log.Info("Namespace ", namespaceName, " created")
				o.Status.NamespaceName = namespaceName
				o.Status.Status = "RUNNING"
				o.Status.LastUpdate = metav1.Now().String()
				if err := r.Status().Update(ctx, o); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("Namespaceconfig ", o.GetName(), " status updated")
				nsLimits := r.nsLimits(o, namespaceName)
				if err := r.Create(ctx, nsLimits); err != nil {
					log.Error("Failed to create LimitRangefor namespace", namespaceName, ". Error: ", err)
					return ctrl.Result{}, err
				} else {
					log.Info("LimitRange created")
				}
				nsQuota := r.nsQuota(o, namespaceName)
				if err := r.Create(ctx, nsQuota); err != nil {
					log.Error("Failed to create ResourceQuota for namespace ", namespaceName, ". Error: ", err)
					return ctrl.Result{}, err
				} else {
					log.Info("ResourceQuota created")
				}
			}
		}
	} else {
		log.Info("DeletionTimestamp is not zero for Namespaceconfig ", o.GetName())
		if controllerutil.ContainsFinalizer(o, finalizerName) {
			log.Info("Finalizer found, removing finalizer ", finalizerName, " from Namespaceconfig ", o.GetName())
			// Remove the finalizer from the Namespaceconfig object once the cleanup succeeded
			// This will free up Namespaceconfig resource to be deleted
			controllerutil.RemoveFinalizer(o, finalizerName)
			if err := r.Update(ctx, o); err != nil {
				log.Error("Unable to remove finalizer and update Namespaceconfig ", err)
				return ctrl.Result{}, err
			}
			log.Info("Finalizer removed from Namespaceconfig ", o.GetName())
		}
		// Attempt to delete the namespace
		err = r.Get(ctx, client.ObjectKey{Name: namespaceName}, namespace)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("Namespace ", namespaceName, " not found")
			} else {
				log.Error("Failed to get Namespace ", namespaceName, " ", err)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Deleting Namespace ", namespaceName)
			if err = r.Delete(ctx, namespace); err != nil {
				log.Error("Failed to delete Namespace ", namespaceName, " ", err)
				return ctrl.Result{}, err
			} else {
				log.Info("Namespace ", namespaceName, " deleted")
			}
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *NamespaceconfigReconciler) nsLimits(nc *namespaceconfigv1.Namespaceconfig, namespaceName string) *corev1.LimitRange {
	log := util.Logs
	var (
		namespaceLimits                              corev1.LimitRangeSpec
		limitTypePodMaxCPU                           string
		limitTypePodMaxMemory                        string
		maxStorage                                   string
		maxLimitRequestRatioCPU                      string
		maxLimitRequestRatioMemory                   string
		limitTypeContainerMaxCPU                     string
		limitTypeContainerMaxMemory                  string
		limitTypeContainerDefaultCPU                 string
		limitTypeContainerDefaultMemory              string
		limitTypeContainerDefaultRequestCPU          string
		limitTypeContainerDefaultRequestMemory       string
		limitTypeContainerMaxLimitRequestRatioCPU    string
		limitTypeContainerMaxLimitRequestRatioMemory string
	)

	if nc.Spec.NamespaceSize == "L" {
		limitTypePodMaxCPU = "4"
		limitTypePodMaxMemory = "4Gi"
		maxStorage = "500Gi"
		maxLimitRequestRatioCPU = "5"
		maxLimitRequestRatioMemory = "10Gi"
		limitTypeContainerMaxCPU = "4"
		limitTypeContainerMaxMemory = "4Gi"
		limitTypeContainerDefaultCPU = "0.5"
		limitTypeContainerDefaultMemory = "512Mi"
		limitTypeContainerDefaultRequestCPU = "0.1"
		limitTypeContainerDefaultRequestMemory = "256Mi"
		limitTypeContainerMaxLimitRequestRatioCPU = "5"
		limitTypeContainerMaxLimitRequestRatioMemory = "5Gi"
	} else if nc.Spec.NamespaceSize == "M" {
		limitTypePodMaxCPU = "2"
		limitTypePodMaxMemory = "2Gi"
		maxStorage = "200Gi"
		maxLimitRequestRatioCPU = "5"
		maxLimitRequestRatioMemory = "5Gi"
		limitTypeContainerMaxCPU = "2"
		limitTypeContainerMaxMemory = "2Gi"
		limitTypeContainerDefaultCPU = "0.5"
		limitTypeContainerDefaultMemory = "512Mi"
		limitTypeContainerDefaultRequestCPU = "0.1"
		limitTypeContainerDefaultRequestMemory = "256Mi"
		limitTypeContainerMaxLimitRequestRatioCPU = "5"
		limitTypeContainerMaxLimitRequestRatioMemory = "5Gi"
	} else if nc.Spec.NamespaceSize == "S" {
		limitTypePodMaxCPU = "1"
		limitTypePodMaxMemory = "1Gi"
		maxStorage = "100Gi"
		maxLimitRequestRatioCPU = "5"
		maxLimitRequestRatioMemory = "2Gi"
		limitTypeContainerMaxCPU = "1"
		limitTypeContainerMaxMemory = "1Gi"
		limitTypeContainerDefaultCPU = "0.5"
		limitTypeContainerDefaultMemory = "512Mi"
		limitTypeContainerDefaultRequestCPU = "0.1"
		limitTypeContainerDefaultRequestMemory = "256Mi"
		limitTypeContainerMaxLimitRequestRatioCPU = "5"
		limitTypeContainerMaxLimitRequestRatioMemory = "5Gi"
	}
	namespaceLimits = corev1.LimitRangeSpec{
		Limits: []corev1.LimitRangeItem{
			{
				Type: corev1.LimitTypePod,
				Max: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(limitTypePodMaxCPU),
					corev1.ResourceMemory: resource.MustParse(limitTypePodMaxMemory),
				},
				MaxLimitRequestRatio: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(maxLimitRequestRatioCPU),
					corev1.ResourceMemory: resource.MustParse(maxLimitRequestRatioMemory),
				},
			},
			{
				Type: corev1.LimitTypeContainer,
				Max: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(limitTypeContainerMaxCPU),
					corev1.ResourceMemory: resource.MustParse(limitTypeContainerMaxMemory),
				},
				Default: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(limitTypeContainerDefaultCPU),
					corev1.ResourceMemory: resource.MustParse(limitTypeContainerDefaultMemory)},
				DefaultRequest: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(limitTypeContainerDefaultRequestCPU),
					corev1.ResourceMemory: resource.MustParse(limitTypeContainerDefaultRequestMemory),
				},
				MaxLimitRequestRatio: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(limitTypeContainerMaxLimitRequestRatioCPU),
					corev1.ResourceMemory: resource.MustParse(limitTypeContainerMaxLimitRequestRatioMemory),
				},
			},
			{
				Type: corev1.LimitTypePersistentVolumeClaim,
				Max:  map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resource.MustParse(maxStorage)},
			},
		},
	}
	// if nc.Spec.NamespaceSize == "S" {
	// 	namespaceLimits = corev1.LimitRangeSpec{
	// 		Limits: []corev1.LimitRangeItem{
	// 			{
	// 				Type: corev1.LimitTypePod,
	// 				Max: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("4"),
	// 					corev1.ResourceMemory: resource.MustParse("4Gi"),
	// 				},
	// 				MaxLimitRequestRatio: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("5"),
	// 					corev1.ResourceMemory: resource.MustParse("10Gi"),
	// 				},
	// 			},
	// 			{
	// 				Type: corev1.LimitTypeContainer,
	// 				Max: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("4"),
	// 					corev1.ResourceMemory: resource.MustParse("4Gi"),
	// 				},
	// 				Default: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("0.5"),
	// 					corev1.ResourceMemory: resource.MustParse("512Mi")},
	// 				DefaultRequest: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("0.1"),
	// 					corev1.ResourceMemory: resource.MustParse("256Mi"),
	// 				},
	// 				MaxLimitRequestRatio: map[corev1.ResourceName]resource.Quantity{
	// 					corev1.ResourceCPU:    resource.MustParse("5"),
	// 					corev1.ResourceMemory: resource.MustParse("5Gi"),
	// 				},
	// 			},
	// 			{
	// 				Type: corev1.LimitTypePersistentVolumeClaim,
	// 				Max:  map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resource.MustParse("500Gi")},
	// 			},
	// 		},
	// 	}
	// }
	limits := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-limits-" + strings.ToLower(nc.Spec.NamespaceSize),
			Namespace: namespaceName,
		},
		Spec: namespaceLimits,
	}
	if err := ctrl.SetControllerReference(nc, limits, r.Scheme); err != nil {
		log.Error("Unable to set ownerReference for namespace ", namespaceName, ". Error: ", err)
	}
	return limits
}

func (r *NamespaceconfigReconciler) nsQuota(nc *namespaceconfigv1.Namespaceconfig, namespaceName string) *corev1.ResourceQuota {
	log := util.Logs
	namespaceQuota := &corev1.ResourceQuotaSpec{}
	if nc.Spec.NamespaceSize == "S" {
		namespaceQuota = &corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceName("cpu"):                    resource.MustParse("8"),
				corev1.ResourceName("memory"):                 resource.MustParse("8Gi"),
				corev1.ResourceName("persistentvolumeclaims"): resource.MustParse("200"),
				corev1.ResourceName("pods"):                   resource.MustParse("200"),
				corev1.ResourceName("replicationcontrollers"): resource.MustParse("200"),
				corev1.ResourceName("services"):               resource.MustParse("200"),
			},
		}
	}
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-limits-" + strings.ToLower(nc.Spec.NamespaceSize),
			Namespace: namespaceName,
		},
		Spec: *namespaceQuota,
	}
	if err := ctrl.SetControllerReference(nc, quota, r.Scheme); err != nil {
		log.Error("Unable to set ownerReference for namespace ", namespaceName, ". Error: ", err)
	}
	return quota
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceconfigReconciler) SetupWithManager(mgr ctrl.Manager) error {

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

	return ctrl.NewControllerManagedBy(mgr).
		For(&namespaceconfigv1.Namespaceconfig{}).
		Owns(&corev1.Namespace{}, predicateNamespace).
		Complete(r)
}
