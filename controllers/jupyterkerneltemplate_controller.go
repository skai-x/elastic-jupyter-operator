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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

// JupyterKernelTemplateReconciler reconciles a JupyterKernelTemplate object
type JupyterKernelTemplateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupyterkerneltemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupyterkerneltemplates/status,verbs=get;update;patch

func (r *JupyterKernelTemplateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("jupyterkerneltemplate", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *JupyterKernelTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeflowtkestackiov1alpha1.JupyterKernelTemplate{}).
		Complete(r)
}
