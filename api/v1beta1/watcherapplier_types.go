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

package v1beta1

import (
	topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/common/tls"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WatcherApplierSpec defines the desired state of WatcherApplier
type WatcherApplierSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	WatcherCommon `json:",inline"`

	// +kubebuilder:validation:Required
	// Secret containing all passwords / keys needed
	Secret string `json:"secret"`

	WatcherSubCrsCommon `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Minimum=0
	// Replicas of Watcher service to run
	Replicas *int32 `json:"replicas"`

	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// TLS - Parameters related to the TLS
	TLS tls.Ca `json:"tls,omitempty"`
}

// WatcherApplierStatus defines the observed state of WatcherApplier
type WatcherApplierStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// ObservedGeneration - the most recent generation observed for this
	// service. If the observed generation is less than the spec generation,
	// then the controller has not processed the latest changes injected by
	// the openstack-operator in the top-level CR (e.g. the ContainerImage)
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ReadyCount of watcher Applier instances
	ReadyCount int32 `json:"readyCount,omitempty"`

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// LastAppliedTopology - the last applied Topology
	LastAppliedTopology *topologyv1.TopoRef `json:"lastAppliedTopology,omitempty"`
}

// WatcherApplierTemplatce defines the input parameters specified by the user to
// create a WatcherApplier via higher level CRDs.
type WatcherApplierTemplate struct {
	WatcherSubCrsTemplate `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Minimum=0
	// Replicas of WatcherApplier service to run
	Replicas *int32 `json:"replicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// WatcherApplier is the Schema for the watcherappliers API
type WatcherApplier struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WatcherApplierSpec   `json:"spec,omitempty"`
	Status WatcherApplierStatus `json:"status,omitempty"`
}

// IsReady returns true if the ReadyCondition is true
func (r *WatcherApplier) IsReady() bool {
	return r.Status.Conditions.IsTrue(condition.ReadyCondition)
}

//+kubebuilder:object:root=true

// WatcherApplierList contains a list of WatcherApplier
type WatcherApplierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WatcherApplier `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WatcherApplier{}, &WatcherApplierList{})
}

// GetSpecTopologyRef - Returns the TopologyRef defined in the Spec
func (instance *WatcherApplier) GetSpecTopologyRef() *topologyv1.TopoRef {
	return instance.Spec.TopologyRef
}

// GetLastAppliedTopology - Returns the LastAppliedTopology Set in the Status
func (instance *WatcherApplier) GetLastAppliedTopology() *topologyv1.TopoRef {
	return instance.Status.LastAppliedTopology
}

// SetLastAppliedTopology - Sets the LastAppliedTopology value in the Status
func (instance *WatcherApplier) SetLastAppliedTopology(topologyRef *topologyv1.TopoRef) {
	instance.Status.LastAppliedTopology = topologyRef
}
