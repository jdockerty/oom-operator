/*
Copyright 2023 Jack Dockerty.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OomerSpec defines the desired state of Oomer
type OomerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is the container image to use for the oomer application, if unspecified will default
	// to the latest version.
	Image *string `json:"image,omitempty"`

	// Replicas is the number of desired OOMKilled pods to deploy.
	Replicas *int32 `json:"replicas"`

	// Labels are passed directly to the oomer application.
	Labels map[string]string `json:"labels,omitempty"`
}

// OomerStatus defines the observed state of Oomer
type OomerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ObservedReplicas are number of observed OOMKilled pods, this should
	// match the number of configured replicas.
	ObservedReplicas *int32 `json:"observedReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Oomer is the Schema for the oomers API
type Oomer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OomerSpec   `json:"spec,omitempty"`
	Status OomerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OomerList contains a list of Oomer
type OomerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Oomer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Oomer{}, &OomerList{})
}
