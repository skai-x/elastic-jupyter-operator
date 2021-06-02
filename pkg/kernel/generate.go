package kernel

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	labelNS       = "namespace"
	labelKernel   = "kernel"
	envKernelID   = "KERNEL_ID"
	labelKernelID = "kernel_id"
)

// generator defines the generator which is used to generate
// desired specs.
type generator struct {
	k *v1alpha1.JupyterKernel
}

// newGenerator creates a new Generator.
func newGenerator(k *v1alpha1.JupyterKernel) (
	*generator, error) {
	if k == nil {
		return nil, fmt.Errorf("Got nil when initializing Generator")
	}
	g := &generator{
		k: k,
	}

	return g, nil
}

func (g generator) DesiredDeployment() (*appsv1.Deployment, error) {
	labels := g.labels()

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.k.Name,
			Namespace: g.k.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Template: g.k.Spec.Template,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	if d.Spec.Template.Labels == nil {
		d.Spec.Template.Labels = make(map[string]string)
	}
	// Set the labels to the pod template.
	for k, v := range labels {
		d.Spec.Template.Labels[k] = v
	}

	// Update the metadata.
	g.hackLabelID(&d.Spec.Template)

	return d, nil
}

func (g generator) labels() map[string]string {
	return map[string]string{
		labelNS:     g.k.Namespace,
		labelKernel: g.k.Name,
	}
}

// hackLabelID copies the ID from environment variables to
// metadata.
// TODO(gaocegege): Use newer version of controller-tools to avoid it.
// https://github.com/kubernetes-sigs/controller-tools/issues/448
func (g generator) hackLabelID(pod *v1.PodTemplateSpec) {
	if pod.Spec.Containers == nil || len(pod.Spec.Containers) == 0 {
		return
	}
	for _, env := range pod.Spec.Containers[0].Env {
		if env.Name == envKernelID {
			if pod.Labels == nil {
				pod.Labels = make(map[string]string)
			}
			pod.Labels[labelKernelID] = env.Value
			return
		}
	}
}
