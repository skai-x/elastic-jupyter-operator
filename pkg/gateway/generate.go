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

// DesiredDeploymentWithoutOwner returns the desired deployment
// without owner.
func (g generator) DesiredDeploymentWithoutOwner() *appsv1.Deployment {
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
