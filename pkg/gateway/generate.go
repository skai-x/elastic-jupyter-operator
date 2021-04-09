/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package gateway

import (
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	defaultImage         = "ccr.ccs.tencentyun.com/kubeflow-oteam/enterprise-gateway:2.5.0"
	defaultContainerName = "gateway"
	// TODO: This default kernel is not used by now
	defaultKernelImage        = "ccr.ccs.tencentyun.com/kubeflow-oteam/jupyter-kernel-py:2.5.0"
	defaultPortName           = "gateway"
	defaultKernel             = "python_kubernetes"
	defaultPort               = 8888
	defaultGatewayClusterRole = "enterprise-gateway-controller"
	defaultServiceAccount     = "enterprise-gateway-sa"

	LabelGateway = "gateway"
	LabelNS      = "namespace"

	cullTimeoutOpt = "--MappingKernelManager.cull_idle_timeout"
	cullInterval   = "--MappingKernelManager.cull_interval"
)

var (
	defaultKernels = "'r_kubernetes','python_kubernetes','python_tf_kubernetes','python_tf_gpu_kubernetes','scala_kubernetes','spark_r_kubernetes','spark_python_kubernetes','spark_scala_kubernetes'"
)

// generator defines the generator which is used to generate
// desired specs.
type generator struct {
	gateway *v1alpha1.JupyterGateway
}

// newGenerator creates a new Generator.
func newGenerator(gateway *v1alpha1.JupyterGateway) (
	*generator, error) {
	if gateway == nil {
		return nil, fmt.Errorf("Got nil when initializing Generator")
	}
	g := &generator{
		gateway: gateway,
	}

	return g, nil
}

// DesiredServiceWithoutOwner returns desired service without
// owner.
func (g generator) DesiredServiceWithoutOwner() *v1.Service {
	labels := g.labels()
	s := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.gateway.Namespace,
			Name:      g.gateway.Name,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Selector:        labels,
			Type:            v1.ServiceTypeClusterIP,
			SessionAffinity: v1.ServiceAffinityClientIP,
			Ports: []v1.ServicePort{
				{
					Name:     defaultPortName,
					Port:     defaultPort,
					Protocol: v1.ProtocolTCP,
				},
			},
		},
	}
	return s
}

func (g generator) DesiredRoleBinding(
	sa *v1.ServiceAccount) *rbacv1.RoleBinding {
	labels := g.labels()
	crb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.gateway.Namespace,
			Name:      g.gateway.Name,
			Labels:    labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Name:     defaultGatewayClusterRole,
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	return crb
}

func (g generator) DesiredServiceAccountWithoutOwner() *v1.ServiceAccount {
	labels := g.labels()
	sa := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind: "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.gateway.Namespace,
			Name:      g.gateway.Name,
			Labels:    labels,
		},
	}
	return sa
}

// DesiredDeploymentWithoutOwner returns the desired deployment
// without owner.
func (g generator) DesiredDeploymentWithoutOwner(
	sa string) *appsv1.Deployment {
	labels := g.labels()
	selector := &metav1.LabelSelector{
		MatchLabels: labels,
	}
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.gateway.Namespace,
			Name:      g.gateway.Name,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					ServiceAccountName: sa,
					Containers: []v1.Container{
						{
							Name:            defaultContainerName,
							Image:           defaultImage,
							ImagePullPolicy: v1.PullIfNotPresent,
							Ports: []v1.ContainerPort{
								{
									Name:          defaultPortName,
									ContainerPort: defaultPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "EG_DEFAULT_KERNEL_NAME",
									Value: g.defaultKernel(),
								},
								{
									Name:  "EG_KERNEL_CLUSTER_ROLE",
									Value: g.defaultClusterRole(),
								},
								{
									Name:  "EG_KERNEL_WHITELIST",
									Value: g.kernels(),
								},
								{
									Name:  "EG_PORT",
									Value: strconv.Itoa(defaultPort),
								},
								{
									Name:  "EG_NAMESPACE",
									Value: g.gateway.Namespace,
								},
								{
									// TODO(gaocegege): Make it configurable.
									Name:  "EG_SHARED_NAMESPACE",
									Value: "true",
								},
								{
									// TODO(gaocegege): Make it configurable.
									Name:  "EG_MIRROR_WORKING_DIRS",
									Value: "false",
								},
								{
									Name:  "EG_CULL_IDLE_TIMEOUT",
									Value: "3600",
								},
								{
									Name:  "EG_KERNEL_LAUNCH_TIMEOUT",
									Value: "60",
								},
								{
									Name:  "EG_KERNEL_IMAGE",
									Value: defaultKernelImage,
								},
							},
						},
					},
				},
			},
		},
	}

	if g.gateway.Spec.CullIdleTimeout != nil {
		env := v1.EnvVar{
			Name:  "EG_CULL_IDLE_TIMEOUT",
			Value: strconv.Itoa(int(*g.gateway.Spec.CullIdleTimeout)),
		}
		d.Spec.Template.Spec.Containers[0].Env = append(
			d.Spec.Template.Spec.Containers[0].Env, env)
	}
	if g.gateway.Spec.CullInterval != nil {
		env := v1.EnvVar{
			Name:  "EG_CULL_INTERVAL",
			Value: strconv.Itoa(int(*g.gateway.Spec.CullInterval)),
		}
		d.Spec.Template.Spec.Containers[0].Env = append(
			d.Spec.Template.Spec.Containers[0].Env, env)
	}
	if g.gateway.Spec.Resources != nil {
		d.Spec.Template.Spec.Containers[0].Resources = *g.gateway.Spec.Resources
	}

	return d
}

func (g generator) defaultClusterRole() string {
	if g.gateway.Spec.ClusterRole != nil {
		return *g.gateway.Spec.ClusterRole
	}
	return defaultGatewayClusterRole
}

func (g generator) labels() map[string]string {
	return map[string]string{
		LabelNS:      g.gateway.Namespace,
		LabelGateway: g.gateway.Name,
	}
}

func (g generator) kernels() string {
	if g.gateway.Spec.Kernels != nil {
		return strings.Join(g.gateway.Spec.Kernels, ",")
	}
	return defaultKernels
}

func (g generator) defaultKernel() string {
	if g.gateway.Spec.DefaultKernel != nil {
		return *g.gateway.Spec.DefaultKernel
	}
	return defaultKernel
}
