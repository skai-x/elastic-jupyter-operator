package notebook

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		},
	}

	notebookWithTemplate = &v1alpha1.JupyterNotebook{
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
		},
	}

	completeNotebook = &v1alpha1.JupyterNotebook{
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
	type test struct {
		input       *v1alpha1.JupyterNotebook
		expectedErr error
		expectedGen *generator
	}

	tests := []test{
		{input: nil, expectedErr: errors.New("Got nil when initializing Generator"), expectedGen: nil},
		{input: notebookWithTemplate, expectedErr: nil, expectedGen: &generator{nb: notebookWithTemplate}},
		{input: notebookWithGateway, expectedErr: nil, expectedGen: &generator{nb: notebookWithGateway}},
		{input: completeNotebook, expectedErr: nil, expectedGen: &generator{nb: completeNotebook}},
		{input: emptyNotebook, expectedErr: nil, expectedGen: &generator{nb: emptyNotebook}},
	}

	for _, tc := range tests {
		gen, err := newGenerator(tc.input)
		if !reflect.DeepEqual(tc.expectedErr, err) {
			t.Errorf("expected: %v, got: %v", tc.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(tc.expectedGen, gen) {
			t.Errorf("expected: %v, got: %v", tc.expectedGen, gen)
		}
	}
}

func TestDesiredDeploymentWithoutOwner(t *testing.T) {
	type test struct {
		gen           *generator
		expectedErr   error
		expectedImage string
		expectedArgs  []string
	}

	tests := []test{
		{gen: &generator{nb: notebookWithTemplate}, expectedErr: nil, expectedImage: DefaultImage, expectedArgs: nil},
		{gen: &generator{nb: notebookWithGateway}, expectedErr: nil, expectedImage: DefaultImageWithGateway,
			expectedArgs: []string{"start-notebook.sh", "--gateway-url", fmt.Sprintf("http://%s.%s:%d", GatewayName, GatewayNamespace, 8888)}},
		{gen: &generator{nb: completeNotebook}, expectedErr: nil, expectedImage: DefaultImage,
			expectedArgs: []string{"--gateway-url", fmt.Sprintf("http://%s.%s:%d", GatewayName, GatewayNamespace, 8888)}},
		{gen: &generator{nb: emptyNotebook}, expectedErr: errors.New("no gateway and template applied")},
	}

	for i, tc := range tests {
		d, err := tc.gen.DesiredDeploymentWithoutOwner()
		if !reflect.DeepEqual(tc.expectedErr, err) {
			t.Errorf("expected: %v, got: %v", tc.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(tc.expectedImage, d.Spec.Template.Spec.Containers[0].Image) {
			t.Errorf("expected: %v, got: %v", tc.expectedImage, d.Spec.Template.Spec.Containers[0].Image)
		}
		if err == nil && !reflect.DeepEqual(tc.expectedArgs, d.Spec.Template.Spec.Containers[0].Args) {
			t.Errorf("i= %d expected: %v, got: %v", i, tc.expectedArgs, d.Spec.Template.Spec.Containers[0].Args)
		}
	}
}

func TestLable(t *testing.T) {
	type test struct {
		gen      *generator
		expected string
	}

	tests := []test{
		{gen: &generator{nb: notebookWithTemplate}, expected: JupyterNotebookName},
	}

	for _, tc := range tests {
		mp := tc.gen.labels()
		if !reflect.DeepEqual(tc.expected, mp["notebook"]) {
			t.Errorf("expected: %v, got: %v", tc.expected, mp["notebook"])
		}
	}
}
