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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JupyterKernelSpecSpec defines the desired state of JupyterKernelSpec
type JupyterKernelSpecSpec struct {
	Language    string      `json:"language,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Image       string      `json:"image,omitempty"`
	Env         []v1.EnvVar `json:"env,omitempty"`
	Command     []string    `json:"command,omitempty"`
	// TODO(gaocegege): Support resources and so on.

	Template *v1.PodTemplate `json:"template,omitempty"`
}

// JupyterKernelSpecStatus defines the observed state of JupyterKernelSpec
type JupyterKernelSpecStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JupyterKernelSpec is the Schema for the jupyterkernelspecs API
type JupyterKernelSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterKernelSpecSpec   `json:"spec,omitempty"`
	Status JupyterKernelSpecStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JupyterKernelSpecList contains a list of JupyterKernelSpec
type JupyterKernelSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterKernelSpec `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterKernelSpec{}, &JupyterKernelSpecList{})
}
