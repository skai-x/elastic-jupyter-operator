package notebook

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.
var (
	cfg        *rest.Config
	k8sClient  client.Client
	k8sManager manager.Manager
	testEnv    *envtest.Environment
	s          *runtime.Scheme

	log = ctrl.Log.WithName("controllers").WithName("JupyterNotebook")
	rec = record.NewFakeRecorder(1024 * 1024)
)

const (
	timeout  = time.Second * 10
	duration = time.Second * 10
	interval = time.Millisecond * 250
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Noteboook Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = kubeflowtkestackiov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())
	s = k8sManager.GetScheme()
	Expect(s).ToNot(BeNil())

	close(done)
}, 60)

var _ = Describe("JupyterNotebook controller", func() {

	Context("Nil JupyterNotebook", func() {
		It("Should fail to NewReconciler", func() {
			_, err := NewReconciler(k8sClient, log, rec, s, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("JupyterNotebook without template and notebook", func() {
		It("Should fail to reconcileDeployment", func() {
			var r *Reconciler
			var err error
			r, err = NewReconciler(k8sClient, log, rec, s, emptyNotebook)
			Expect(err).ToNot(HaveOccurred())
			Expect(r).ToNot(BeNil())
			err = r.reconcileDeployment()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("JupyterNotebook only have template", func() {
		It("Should reconcile deployment as desired", func() {
			var r *Reconciler
			var err error
			r, err = NewReconciler(k8sClient, log, rec, s, notebookWithTemplate)
			Expect(err).ToNot(HaveOccurred())
			Expect(r).ToNot(BeNil())

			err = r.cli.Create(context.TODO(), notebookWithTemplate)
			Expect(err).ToNot(HaveOccurred())

			err = r.reconcileDeployment()
			Expect(err).ToNot(HaveOccurred())

			By("Expecting template name")
			Eventually(func() string {
				actual := &kubeflowtkestackiov1alpha1.JupyterNotebook{}
				if err := k8sClient.Get(context.Background(),
					types.NamespacedName{Name: notebookWithTemplate.GetName(), Namespace: notebookWithTemplate.GetNamespace()}, actual); err == nil {
					return actual.Spec.Template.Spec.Containers[0].Name
				}
				return ""
			}, timeout, interval).Should(Equal(notebookWithTemplate.Spec.Template.Spec.Containers[0].Name))

			err = r.Reconcile()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
