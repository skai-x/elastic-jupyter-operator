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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	"github.com/tkestack/elastic-jupyter-operator/pkg/kernelspec"
)

// JupyterKernelSpecReconciler reconciles a JupyterKernelSpec object
type JupyterKernelSpecReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupyterkernelspecs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupyterkernelspecs/status,verbs=get;update;patch

func (r *JupyterKernelSpecReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("jupyterkernelspec", req.NamespacedName)

	original := &v1alpha1.JupyterKernelSpec{}

	err := r.Get(context.TODO(), req.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Error(err, "Failed to get the object, requeuing the request")
		return ctrl.Result{}, err
	}
	instance := original.DeepCopy()

	gr, err := kernelspec.NewReconciler(r.Client, r.Log, r.Recorder, r.Scheme, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := gr.Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *JupyterKernelSpecReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeflowtkestackiov1alpha1.JupyterKernelSpec{}).
		Complete(r)
}
