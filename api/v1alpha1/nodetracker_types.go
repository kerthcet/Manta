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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeTrackerSpec defines the desired state of NodeTracker
// It acts like a cache.
type NodeTrackerSpec struct {
	// TODO: GC disk rate.
}

// NodeTrackerStatus defines the observed state of NodeTracker
type NodeTrackerStatus struct {
	Files []FileTracker `json:"files,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// NodeTracker is the Schema for the nodetrackers API
type NodeTracker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeTrackerSpec   `json:"spec,omitempty"`
	Status NodeTrackerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeTrackerList contains a list of NodeTracker
type NodeTrackerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeTracker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeTracker{}, &NodeTrackerList{})
}
