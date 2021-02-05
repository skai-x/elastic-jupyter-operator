/*
Copyright 2021.

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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JupyterGatewaySpec defines the desired state of JupyterGateway
type JupyterGatewaySpec struct {
	// We will add kernels at runtime, thus we do not make it a type.
	Kernels       []string `json:"kernels,omitempty"`
	DefaultKernel *string  `json:"defaultKernel,omitempty"`
}

type KernelType string

// JupyterGatewayStatus defines the observed state of JupyterGateway
type JupyterGatewayStatus struct {
	appsv1.DeploymentStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JupyterGateway is the Schema for the jupytergateways API
type JupyterGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterGatewaySpec   `json:"spec,omitempty"`
	Status JupyterGatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JupyterGatewayList contains a list of JupyterGateway
type JupyterGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterGateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterGateway{}, &JupyterGatewayList{})
}
