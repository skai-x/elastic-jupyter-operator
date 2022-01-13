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
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	defaultImage         = "jupyter/base-notebook:python-3.9.7"
	defaultContainerName = "notebook"
	defaultPortName      = "notebook"
	defaultPort          = 8888

	LabelNotebook = "notebook"
	LabelNS       = "namespace"

	argumentGatewayURL       = "--gateway-url"
	argumentNotebookToken    = "--NotebookApp.token"
	argumentNotebookPassword = "--NotebookApp.password"
)

type generator struct {
	nb *v1alpha1.JupyterNotebook
}

// newGenerator creates a new Generator.
func newGenerator(nb *v1alpha1.JupyterNotebook) (
	*generator, error) {
	if nb == nil {
		return nil, fmt.Errorf("the notebook is null")
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
	terminationGracePeriodSeconds := int64(30)

	if g.nb.Spec.Template != nil {
		if g.nb.Spec.Template.Labels != nil {
			for k, v := range g.nb.Spec.Template.Labels {
				podLabels[k] = v
			}
		}
		podSpec = completePodSpec(&g.nb.Spec.Template.Spec)
	} else {
		podSpec = v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:                     defaultContainerName,
					Image:                    defaultImage,
					ImagePullPolicy:          v1.PullIfNotPresent,
					TerminationMessagePath:   v1.TerminationMessagePathDefault,
					TerminationMessagePolicy: v1.TerminationMessageReadFile,
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
			RestartPolicy:                 v1.RestartPolicyAlways,
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
			DNSPolicy:                     v1.DNSClusterFirst,
			SecurityContext:               &v1.PodSecurityContext{},
			SchedulerName:                 v1.DefaultSchedulerName,
		}
	}

	replicas := int32(1)
	revisionHistoryLimit := int32(10)
	progressDeadlineSeconds := int32(600)
	maxUnavailable := intstr.FromInt(25)

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.nb.Namespace,
			Name:      g.nb.Name,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: podSpec,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxUnavailable,
				},
			},
			RevisionHistoryLimit:    &revisionHistoryLimit,
			ProgressDeadlineSeconds: &progressDeadlineSeconds,
		},
	}

	if g.nb.Spec.Gateway != nil {
		gatewayURL := fmt.Sprintf("http://%s.%s:%d",
			g.nb.Spec.Gateway.Name, g.nb.Spec.Gateway.Namespace, defaultPort)
		d.Spec.Template.Spec.Containers[0].Args = append(
			d.Spec.Template.Spec.Containers[0].Args, argumentGatewayURL, gatewayURL)
	}

	// Set the auth configuration to notebook instance.
	if g.nb.Spec.Auth != nil {
		auth := g.nb.Spec.Auth
		// Set the token and password to empty.
		if auth.Mode == v1alpha1.ModeJupyterAuthDisable {
			d.Spec.Template.Spec.Containers[0].Args = append(
				d.Spec.Template.Spec.Containers[0].Args,
				argumentNotebookToken, "",
				argumentNotebookPassword, "",
			)
		} else {
			if auth.Token != nil {
				d.Spec.Template.Spec.Containers[0].Args = append(
					d.Spec.Template.Spec.Containers[0].Args,
					argumentNotebookToken, *auth.Token,
				)
			}
			if auth.Password != nil {
				d.Spec.Template.Spec.Containers[0].Args = append(
					d.Spec.Template.Spec.Containers[0].Args,
					argumentNotebookPassword, *auth.Password,
				)
			}
		}
	}

	return d, nil
}

func (g generator) labels() map[string]string {
	return map[string]string{
		LabelNS:       g.nb.Namespace,
		LabelNotebook: g.nb.Name,
	}
}

func completePodSpec(old *v1.PodSpec) v1.PodSpec {
	new := old.DeepCopy()
	for i := range new.Containers {
		if new.Containers[i].TerminationMessagePath == "" {
			new.Containers[i].TerminationMessagePath = v1.TerminationMessagePathDefault
		}
		if new.Containers[i].TerminationMessagePolicy == v1.TerminationMessagePolicy("") {
			new.Containers[i].TerminationMessagePolicy = v1.TerminationMessageReadFile
		}
		if new.Containers[i].ImagePullPolicy == v1.PullPolicy("") {
			new.Containers[i].ImagePullPolicy = v1.PullIfNotPresent
		}
	}

	if new.RestartPolicy == v1.RestartPolicy("") {
		new.RestartPolicy = v1.RestartPolicyAlways
	}

	if new.TerminationGracePeriodSeconds == nil {
		d := int64(v1.DefaultTerminationGracePeriodSeconds)
		new.TerminationGracePeriodSeconds = &d
	}

	if new.DNSPolicy == v1.DNSPolicy("") {
		new.DNSPolicy = v1.DNSClusterFirst
	}

	if new.SecurityContext == nil {
		new.SecurityContext = &v1.PodSecurityContext{}
	}

	if new.SchedulerName == "" {
		new.SchedulerName = v1.DefaultSchedulerName
	}

	return *new
}
