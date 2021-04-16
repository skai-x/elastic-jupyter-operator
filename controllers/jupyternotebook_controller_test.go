package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

var _ = Describe("JupyterNotebook controller", func() {
	const (
		JupyterNotebookName      = "jupyternotebook-sample"
		JupyterNotebookNamespace = "default"
		DefaultContainerName     = "notebook"
		DefaultImage             = "busysandbox"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
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

		NotebookWithTemplate = &kubeflowtkestackiov1alpha1.JupyterNotebook{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "kubeflow.tkestack.io/v1alpha1",
				Kind:       "JupyterNotebook",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      JupyterNotebookName,
				Namespace: JupyterNotebookNamespace,
			},
			Spec: kubeflowtkestackiov1alpha1.JupyterNotebookSpec{
				Template: &v1.PodTemplateSpec{
					Spec: podSpec,
				},
			},
		}

		key = types.NamespacedName{
			Name:      JupyterNotebookName,
			Namespace: JupyterNotebookNamespace,
		}
	)

	Context("When updating JupyterNotebook", func() {
		It("Should generate correct JupyterNotebook image", func() {
			By("Create new notebook")
			Expect(k8sClient.Create(context.Background(), NotebookWithTemplate)).Should(Succeed())

			By("Expecting obtained image to be same as default image")
			Eventually(func() string {
				actual := &kubeflowtkestackiov1alpha1.JupyterNotebook{}
				if err := k8sClient.Get(context.Background(), key, actual); err == nil {
					if actual != nil && actual.Spec.Template != nil {
						println(actual)
						return actual.Spec.Template.Name
					}
				}
				return ""
			}, timeout, interval).Should(Equal(DefaultContainerName))
		}, timeout.Hours())
	})
})
