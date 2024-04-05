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

package v1

import (
	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NamespaceconfigSpec defines the desired state of Namespaceconfig
type NamespaceconfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MaxLength=5
	Environment string `json:"Environment"`
	//+kubebuilder:validation:MaxLength=8
	Abbreviation string `json:"Abbreviation"`
	// NamespaceLimits v1.LimitRangeSpec    `json:"NamespaceLimits,omitempty"`
	// NamespaceQuota  v1.ResourceQuotaSpec `json:"NamespaceQuota,omitempty"`
	NamespaceOwner string `json:"NamespaceOwner,omitempty"`
	//+kubebuilder:validation:Enum=S;M;L
	NamespaceSize string `json:"NamespaceSize,omitempty"`
}

// NamespaceconfigStatus defines the observed state of Namespaceconfig
type NamespaceconfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NamespaceName string `json:"NamespaceName,omitempty"`
	Status        string `json:"Status,omitempty"`
	LastUpdate    string `json:"LastUpdate,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName={"nsc","nc","nsconfig"}

// Namespaceconfig is the Schema for the namespaceconfigs API
type Namespaceconfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceconfigSpec   `json:"spec,omitempty"`
	Status NamespaceconfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NamespaceconfigList contains a list of Namespaceconfig
type NamespaceconfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Namespaceconfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Namespaceconfig{}, &NamespaceconfigList{})
}
