package notebook

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
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
				Name:            DefaultContainerName,
				Image:           DefaultImage,
				ImagePullPolicy: v1.PullIfNotPresent,
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
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Generate Suite",
		[]Reporter{printer.NewlineReporter{}})

	By("Get generator")
	var err error
	var g1, g2, g3, g4 *generator
	_, err = newGenerator(nil)
	Expect(err).To(HaveOccurred())
	g1, err = newGenerator(notebookWithTemplate)
	Expect(err).NotTo(HaveOccurred())
	g2, err = newGenerator(notebookWithGateway)
	Expect(err).NotTo(HaveOccurred())
	g3, err = newGenerator(completeNotebook)
	Expect(err).NotTo(HaveOccurred())
	g4, err = newGenerator(emptyNotebook)
	Expect(err).NotTo(HaveOccurred())

	By("Generate deployment with template")
	var d1 *appsv1.Deployment
	d1, err = g1.DesiredDeploymentWithoutOwner()
	Expect(err).NotTo(HaveOccurred())
	Expect(d1.Spec.Template.Spec.Containers[0].Name).Should(Equal(DefaultContainerName))

	By("Generate deployment with gateway")
	var d2 *appsv1.Deployment
	d2, err = g2.DesiredDeploymentWithoutOwner()
	Expect(err).NotTo(HaveOccurred())
	Expect(d2.Spec.Template.Spec.Containers[0].Image).Should(Equal(DefaultImageWithGateway))

	By("Generate deployment with both template and gateway")
	var d3 *appsv1.Deployment
	d3, err = g3.DesiredDeploymentWithoutOwner()
	Expect(err).NotTo(HaveOccurred())
	Expect(d3.Spec.Template.Spec.Containers[0].Image).Should(Equal(DefaultImage))
	s := []string{"--gateway-url", fmt.Sprintf("http://%s.%s:%d", GatewayName, GatewayNamespace, 8888)}
	Expect(d3.Spec.Template.Spec.Containers[0].Args).Should(Equal(s))

	By("Generate deployment")
	_, err = g4.DesiredDeploymentWithoutOwner()
	Expect(err).To(HaveOccurred())

	By("Check labels")
	mp := g1.labels()
	Expect(mp["notebook"]).Should(Equal(JupyterNotebookName))
}
