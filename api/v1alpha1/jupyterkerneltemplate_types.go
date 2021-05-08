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

// JupyterKernelTemplateSpec defines the desired state of JupyterKernelTemplate
type JupyterKernelTemplateSpec struct {
	Template *v1.PodTemplate `json:"template,omitempty"`
}

// JupyterKernelTemplateStatus defines the observed state of JupyterKernelTemplate
type JupyterKernelTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JupyterKernelTemplate is the Schema for the jupyterkerneltemplates API
type JupyterKernelTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterKernelTemplateSpec   `json:"spec,omitempty"`
	Status JupyterKernelTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JupyterKernelTemplateList contains a list of JupyterKernelTemplate
type JupyterKernelTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterKernelTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterKernelTemplate{}, &JupyterKernelTemplateList{})
}
