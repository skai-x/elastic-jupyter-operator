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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	defaultImage         = "jupyter/base-notebook:python-3.8.6"
	defaultContainerName = "notebook"
	defaultPortName      = "notebook"
	defaultPort          = 8888

	LabelNotebook = "notebook"
	LabelNS       = "namespace"
)

type generator struct {
	nb *v1alpha1.JupyterNotebook
}

// newGenerator creates a new Generator.
func newGenerator(nb *v1alpha1.JupyterNotebook) (
	*generator, error) {
	if nb == nil {
		return nil, fmt.Errorf("Got nil when initializing Generator")
	}
	g := &generator{
		nb: nb,
	}

	return g, nil
}

func (g generator) DesiredDeploymentWithoutOwner() (*appsv1.Deployment, error) {
	if g.nb.Spec.Template == nil && g.nb.Spec.Gateway == nil {
		return nil, fmt.Errorf("no gateway and template applied")
	}

	podSpec := v1.PodSpec{}
	podLabels := g.labels()
	labels := g.labels()
	selector := &metav1.LabelSelector{
		MatchLabels: labels,
	}

	if g.nb.Spec.Template != nil {
		if g.nb.Spec.Template.Labels != nil {
			for k, v := range g.nb.Spec.Template.Labels {
				podLabels[k] = v
			}
		}
		podSpec = g.nb.Spec.Template.Spec
	} else {
		podSpec = v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            defaultContainerName,
					Image:           defaultImage,
					ImagePullPolicy: v1.PullIfNotPresent,
					Args: []string{
						"start-notebook.sh",
					},
					Ports: []v1.ContainerPort{
						{
							Name:          defaultPortName,
							ContainerPort: defaultPort,
							Protocol:      v1.ProtocolTCP,
						},
					},
				},
			},
		}
	}

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.nb.Namespace,
			Name:      g.nb.Name,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: podSpec,
			},
		},
	}

	if g.nb.Spec.Gateway != nil {
		gatewayURL := fmt.Sprintf("http://%s.%s:%d",
			g.nb.Spec.Gateway.Name, g.nb.Spec.Gateway.Namespace, defaultPort)
		d.Spec.Template.Spec.Containers[0].Args = append(
			d.Spec.Template.Spec.Containers[0].Args, "--gateway-url", gatewayURL)
	}

	return d, nil
}

func (g generator) labels() map[string]string {
	return map[string]string{
		LabelNS:       g.nb.Namespace,
		LabelNotebook: g.nb.Name,
	}
}
