package gateway

import (
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/jupyter-operator/api/v1alpha1"
)

const (
	defaultImage         = "elyra/enterprise-gateway:dev"
	defaultContainerName = "gateway"
	defaultPortName      = "gateway"
	defaultKernel        = "python_kubernetes"
	defaultPort          = 8888

	LabelGateway = "gateway"
	LabelNS      = "namespace"
)

var (
	defaultKernels = "'r_kubernetes','python_kubernetes','python_tf_kubernetes','python_tf_gpu_kubernetes','scala_kubernetes','spark_r_kubernetes','spark_python_kubernetes','spark_scala_kubernetes'"
)

type Generator struct {
	gateway *v1alpha1.JupyterGateway
}

func NewGenerator(gateway *v1alpha1.JupyterGateway) (*Generator, error) {
	if gateway == nil {
		return nil, fmt.Errorf("Got nil when initializing Generator")
	}
	g := &Generator{
		gateway: gateway,
	}

	return g, nil
}

func (g Generator) TemplateWithoutOwner() *appsv1.Deployment {
	labels := map[string]string{
		LabelNS:      g.gateway.Namespace,
		LabelGateway: g.gateway.Name,
	}

	selector := &metav1.LabelSelector{
		MatchLabels: labels,
	}
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.gateway.Namespace,
			Name:      g.gateway.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
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
							},
						},
					},
				},
			},
		},
	}

	return d
}

func (g Generator) kernels() string {
	if g.gateway.Spec.Kernels != nil {
		return strings.Join(g.gateway.Spec.Kernels, ",")
	}
	return defaultKernels
}

func (g Generator) defaultKernel() string {
	if g.gateway.Spec.DefaultKernel != nil {
		return *g.gateway.Spec.DefaultKernel
	}
	return defaultKernel
}
