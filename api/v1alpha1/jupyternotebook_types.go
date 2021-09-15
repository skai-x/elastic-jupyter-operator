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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JupyterNotebookSpec defines the desired state of JupyterNotebook
type JupyterNotebookSpec struct {
	Gateway *v1.ObjectReference `json:"gateway,omitempty"`
	Auth    *JupyterAuth        `json:"auth,omitempty"`

	Template *v1.PodTemplateSpec `json:"template,omitempty"`
}

// JupyterAuth defines how to deal with jupyter notebook tokens or passwords.
// https://jupyter-notebook.readthedocs.io/en/stable/security.html
type JupyterAuth struct {
	// TODO(gaocegege): Is this field necessary since we make Token and Password a pointer?
	Mode     ModeJupyterAuth `json:"mode,omitempty"`
	Token    *string         `json:"token,omitempty"`
	Password *string         `json:"password,omitempty"`
}

type ModeJupyterAuth string

const (
	ModeJupyterAuthEnable ModeJupyterAuth = "enable"
	// ModeJupyterAuthDisable disables authentication altogether by setting the token
	// and password to empty strings, but this is NOT RECOMMENDED, unless authentication
	// or access restrictions are handled at a different layer in your web application
	ModeJupyterAuthDisable ModeJupyterAuth = "disable"
)

// JupyterNotebookStatus defines the observed state of JupyterNotebook
type JupyterNotebookStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JupyterNotebook is the Schema for the jupyternotebooks API
type JupyterNotebook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterNotebookSpec   `json:"spec,omitempty"`
	Status JupyterNotebookStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JupyterNotebookList contains a list of JupyterNotebook
type JupyterNotebookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterNotebook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterNotebook{}, &JupyterNotebookList{})
}
