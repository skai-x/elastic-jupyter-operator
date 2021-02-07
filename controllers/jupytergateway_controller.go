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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/tkestack/jupyter-operator/api/v1alpha1"
	kubeflowtkestackiov1alpha1 "github.com/tkestack/jupyter-operator/api/v1alpha1"
	"github.com/tkestack/jupyter-operator/pkg/gateway"
)

// JupyterGatewayReconciler reconciles a JupyterGateway object
type JupyterGatewayReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupytergateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeflow.tkestack.io,resources=jupytergateways/status,verbs=get;update;patch

func (r *JupyterGatewayReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("jupytergateway", req.NamespacedName)

	original := &v1alpha1.JupyterGateway{}

	err := r.Get(context.TODO(), req.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Error(err, "Failed to get the object, requeuing the request")
		return reconcile.Result{}, err
	}
	instance := original.DeepCopy()

	g, err := gateway.NewGenerator(instance)
	if err != nil {
		// Error reading the object - requeue the request.
		r.Log.Error(err, "Failed to create the generator, requeuing the request")
		return reconcile.Result{}, err
	}
	desired := g.TemplateWithoutOwner()

	if err := controllerutil.SetControllerReference(
		instance, desired, r.Scheme); err != nil {
		r.Log.Error(err,
			"Set controller reference error, requeuing the request")
		return reconcile.Result{}, err
	}

	actual := &appsv1.Deployment{}
	err = r.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating deployment", "namespace", desired.Namespace, "name", desired.Name)

		if err := r.Create(context.TODO(), desired); err != nil {
			r.Log.Error(err, "CreateSuggestion failed",
				"deployment", desired.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		r.Log.Error(err, "failed to get the expected deployment",
			"deployment", desired.Name)
		return reconcile.Result{}, err
	}

	if !equality.Semantic.DeepEqual(instance.Status.DeploymentStatus, actual.Status) {
		instance.Status.DeploymentStatus = actual.Status
		if err := r.Status().Update(context.TODO(), instance); err != nil {
			r.Log.Error(err, "failed to update status",
				"namespace", instance.Namespace,
				"jupytergateway", instance.Name)
			return reconcile.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *JupyterGatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeflowtkestackiov1alpha1.JupyterGateway{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &v1alpha1.JupyterGateway{},
			}).
		Complete(r)
}
