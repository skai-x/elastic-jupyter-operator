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

// JupyterKernelSpec defines the desired state of JupyterKernel
type JupyterKernelCRDSpec struct {
	Template *v1.PodTemplate `json:"template,omitempty"`
}

// JupyterKernelStatus defines the observed state of JupyterKernel
type JupyterKernelStatus struct {
	// Conditions is an array of current observed job conditions.
	Conditions []JupyterKernelCondition `json:"conditions"`

	// Represents time when the job was acknowledged by the job controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the job was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

type JupyterKernelCondition struct {
	// Type of job condition.
	Type JupyterKernelConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

type JupyterKernelConditionType string

const (
	JupyterKernelRunning   JupyterKernelConditionType = "Running"
	JupyterKernelFailed    JupyterKernelConditionType = "Failed"
	JupyterKernelSucceeded JupyterKernelConditionType = "Succeeded"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JupyterKernel is the Schema for the jupyterkernels API
type JupyterKernel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterKernelCRDSpec `json:"spec,omitempty"`
	Status JupyterKernelStatus  `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JupyterKernelList contains a list of JupyterKernel
type JupyterKernelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterKernel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterKernel{}, &JupyterKernelList{})
}
