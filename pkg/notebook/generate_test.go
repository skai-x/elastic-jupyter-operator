package notebook

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	JupyterNotebookName      = "jupyternotebook-sample"
	JupyterNotebookNamespace = "default"
	DefaultContainerName     = "notebook"
	DefaultImage             = "busysandbox"
	DefaultImageWithGateway  = "jupyter/base-notebook:python-3.8.6"
	GatewayName              = "gateway"
	GatewayNamespace         = "default"
)

var (
	podSpec = v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  DefaultContainerName,
				Image: DefaultImage,
			},
		},
	}

	emptyNotebook = &v1alpha1.JupyterNotebook{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.tkestack.io/v1alpha1",
			Kind:       "JupyterNotebook",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		},
	}

	notebookWithTemplate = &v1alpha1.JupyterNotebook{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.tkestack.io/v1alpha1",
			Kind:       "JupyterNotebook",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		},
		Spec: v1alpha1.JupyterNotebookSpec{
			Template: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "notebook"},
				},
				Spec: podSpec,
			},
		},
	}

	notebookWithGateway = &v1alpha1.JupyterNotebook{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.tkestack.io/v1alpha1",
			Kind:       "JupyterNotebook",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		},
		Spec: v1alpha1.JupyterNotebookSpec{
			Gateway: &v1.ObjectReference{
				Kind: "JupyterGateway",
			},
		},
	}

	completeNotebook = &v1alpha1.JupyterNotebook{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.tkestack.io/v1alpha1",
			Kind:       "JupyterNotebook",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		},
		Spec: v1alpha1.JupyterNotebookSpec{
			Gateway: &v1.ObjectReference{
				Kind:      "JupyterGateway",
				Namespace: GatewayNamespace,
				Name:      GatewayName,
			},
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
		},
	}
)

func TestGenerate(t *testing.T) {
	// Test newGenerator
	var err error
	var g1, g2, g3, g4 *generator
	_, err = newGenerator(nil)
	if err == nil {
		t.Errorf("Expect error to have occurred")
	}
	g1, err = newGenerator(notebookWithTemplate)
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	g2, err = newGenerator(notebookWithGateway)
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	g3, err = newGenerator(completeNotebook)
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	g4, err = newGenerator(emptyNotebook)
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}

	// Test DesiredDeploymentWithoutOwner
	// Generate deployment with template
	var d1 *appsv1.Deployment
	d1, err = g1.DesiredDeploymentWithoutOwner()
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	if d1.Spec.Template.Spec.Containers[0].Name != DefaultContainerName {
		t.Errorf("Actual: %s, Expected: %s", d1.Spec.Template.Spec.Containers[0].Name, DefaultContainerName)
	}

	// Generate deployment with gateway
	var d2 *appsv1.Deployment
	d2, err = g2.DesiredDeploymentWithoutOwner()
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	if d2.Spec.Template.Spec.Containers[0].Image != DefaultImageWithGateway {
		t.Errorf("Actual: %s, Expected: %s", d2.Spec.Template.Spec.Containers[0].Image, DefaultImageWithGateway)
	}

	// Generate deployment with both template and gateway
	var d3 *appsv1.Deployment
	d3, err = g3.DesiredDeploymentWithoutOwner()
	if err != nil {
		t.Errorf("Expect error not to have occurred")
	}
	if d3.Spec.Template.Spec.Containers[0].Image != DefaultImage {
		t.Errorf("Actual: %s, Expected: %s", d3.Spec.Template.Spec.Containers[0].Image, DefaultImage)
	}
	s := []string{"--gateway-url", fmt.Sprintf("http://%s.%s:%d", GatewayName, GatewayNamespace, 8888)}
	if !reflect.DeepEqual(d3.Spec.Template.Spec.Containers[0].Args, s) {
		t.Errorf("Actual: %s, Expected: %s", d3.Spec.Template.Spec.Containers[0].Args, s)
	}

	// Generate deployment
	_, err = g4.DesiredDeploymentWithoutOwner()
	if err == nil {
		t.Errorf("Expect error to have occurred")
	}

	// Test lables
	mp := g1.labels()
	if mp["notebook"] != JupyterNotebookName {
		t.Errorf("Actual: %s, Expected: %s", mp["notebook"], JupyterNotebookName)
	}
}
