/*
Copyright 2020 wuming.lgy@alibaba-inc.com.

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
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ManorSpec defines the desired state of Manor
type ManorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Size      *int32                  `json:"size"`
	Image     string                  `json:"image"`
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	Envs      []v1.EnvVar             `json:"envs,omitempty"`
	Ports     []v1.ServicePort        `json:"ports,omitempty"`
}

// ManorStatus defines the observed state of Manor
type ManorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	appv1.DeploymentStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Manor is the Schema for the manors API
type Manor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManorSpec   `json:"spec,omitempty"`
	Status ManorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ManorList contains a list of Manor
type ManorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Manor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Manor{}, &ManorList{})
}
