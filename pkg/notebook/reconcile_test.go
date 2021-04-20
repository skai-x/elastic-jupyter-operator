package notebook

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
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

var cfg *rest.Config
var k8sClient client.Client
var k8sManager manager.Manager
var testEnv *envtest.Environment

var s *runtime.Scheme
var (
	log = ctrl.Log.WithName("controllers").WithName("JupyterNotebook")
	rec = record.NewFakeRecorder(1024 * 1024)
	// rec = k8sManager.GetEventRecorderFor("reconciler")
	// scheme = scheme.Scheme
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
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = kubeflowtkestackiov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
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
		It("Should fail to create reconciler", func() {
			_, err := NewReconciler(k8sClient, log, rec, s, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("JupyterNotebook only have template", func() {
		It("Should create reconciler successfully", func() {
			r, err := NewReconciler(k8sClient, log, rec, s, notebookWithTemplate)
			Expect(err).ToNot(HaveOccurred())
			Expect(r).ToNot(BeNil())
		})

		It("Should reconcile deployment as desired", func() {
			var r *Reconciler
			var err error
			r, err = NewReconciler(k8sClient, log, rec, s, notebookWithTemplate)
			Expect(err).ToNot(HaveOccurred())
			Expect(r).ToNot(BeNil())
			println("my client is ", r.cli)

			err = r.cli.Create(context.TODO(), notebookWithTemplate)
			Expect(err).ToNot(HaveOccurred())

			err = r.reconcileDeployment()
			Expect(err).ToNot(HaveOccurred())
			// actual := &appsv1.Deployment{}
			// err := k8sClient.Get(context.Background(),
			// 	types.NamespacedName{Name: notebookWithTemplate.GetName(), Namespace: notebookWithTemplate.GetNamespace()}, actual)
			// Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("JupyterNotebook without template and notebook", func() {
		r, err := NewReconciler(k8sClient, log, rec, s, notebookWithTemplate)
		It("Should fail to generate deployment", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(r).ToNot(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
