package kernel

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	labelNS     = "namespace"
	labelKernel = "kernel"
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

func (g generator) DesiredDeployment() (*v1.Deployment, error) {
	labels := g.labels()

	d := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.k.Name,
			Namespace: g.k.Namespace,
			Labels:    labels,
		},
		Spec: v1.DeploymentSpec{
			Template: *g.k.Spec.Template,
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

	return d, nil
}

func (g generator) labels() map[string]string {
	return map[string]string{
		labelNS:     g.k.Namespace,
		labelKernel: g.k.Name,
	}
}
