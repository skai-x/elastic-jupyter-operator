// Tencent is pleased to support the open source community by making TKEStack
// available.
//
// Copyright (C) 2012-2020 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License. You may obtain a copy of the
// License at
//
// https://opensource.org/licenses/Apache-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

package notebook

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

type Reconciler struct {
	cli      client.Client
	log      logr.Logger
	recorder record.EventRecorder
	scheme   *runtime.Scheme

	instance *v1alpha1.JupyterNotebook
	gen      *generator
}

func NewReconciler(cli client.Client,
	l logr.Logger,
	r record.EventRecorder, s *runtime.Scheme,
	i *v1alpha1.JupyterNotebook) (*Reconciler, error) {
	g, err := newGenerator(i)
	if err != nil {
		return nil, err
	}
	return &Reconciler{
		cli:      cli,
		log:      l,
		recorder: r,
		scheme:   s,
		instance: i,
		gen:      g,
	}, nil
}

func (r Reconciler) Reconcile() error {
	if err := r.reconcileDeployment(); err != nil {
		return err
	}
	return nil
}

func (r Reconciler) reconcileDeployment() error {
	desired, err := r.gen.DesiredDeploymentWithoutOwner()
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return err
	}

	actual := &appsv1.Deployment{}
	err = r.cli.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)

	// Create deployment if not found
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating deployment", "namespace", desired.Namespace, "name", desired.Name)
		if err := r.cli.Create(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to create the deployment",
				"deployment", desired.Name)
			return err
		}
	} else if err != nil {
		r.log.Error(err, "failed to get the expected deployment",
			"deployment", desired.Name)
		return err
	}

	// Update deployment from desired to actural
	if !equality.Semantic.DeepEqual(desired.Spec, actual.Spec) {
		if err := r.cli.Update(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to update deployment")
			return err
		} else {
			r.log.Info("deployment updated")
		}
	}

	return nil
}
