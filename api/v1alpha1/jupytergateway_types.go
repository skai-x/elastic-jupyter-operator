// Tencent is pleased to support the open source community by making TKEStack
// available.

// Copyright (C) 2012-2020 Tencent. All Rights Reserved.

// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License. You may obtain a copy of the
// License at

// https://opensource.org/licenses/Apache-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JupyterGatewaySpec defines the desired state of JupyterGateway
type JupyterGatewaySpec struct {
	// Knernels defines the kernels in the gateway.
	// We will add kernels at runtime, thus we do not make it a type.
	Kernels []string `json:"kernels,omitempty"`
	// DefaultKernel defines the default kernel in the gateway.
	DefaultKernel *string `json:"defaultKernel,omitempty"`
	// Timeout (in seconds) after which a kernel is considered idle and
	// ready to be culled. Values of 0 or lower disable culling. Very
	// short timeouts may result in kernels being culled for users
	// with poor network connections.
	// Ref https://jupyter-notebook.readthedocs.io/en/stable/config.html
	CullIdleTimeout *int32 `json:"cullIdleTimeout,omitempty"`

	// The interval (in seconds) on which to check for idle kernels
	// exceeding the cull timeout value.
	CullInterval *int32 `json:"cullInterval,omitempty"`

	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`

	Image string `json:"image,omitempty"`

	ClusterRole *string `json:"clusterRole,omitempty"`
}

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
