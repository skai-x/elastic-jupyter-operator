package kernel

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

type Reconciler struct {
	cli      client.Client
	log      logr.Logger
	recorder record.EventRecorder
	scheme   *runtime.Scheme

	instance *v1alpha1.JupyterKernel
	gen      *generator
}

func NewReconciler(cli client.Client, l logr.Logger,
	r record.EventRecorder, s *runtime.Scheme,
	i *v1alpha1.JupyterKernel) (*Reconciler, error) {
	g, err := newGenerator(i)
	if err != nil {
		return nil, err
	}
	return &Reconciler{
		cli:      cli,
		log:      l,
		recorder: r,
		scheme:   s,
		instance: i,
		gen:      g,
	}, nil
}

func (r Reconciler) Reconcile() error {
	if err := r.reconcileDeployment(); err != nil {
		return err
	}

	return nil
}

func (r Reconciler) reconcileDeployment() error {
	desired, err := r.gen.DesiredDeployment()
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return err
	}

	actual := &appsv1.Deployment{}
	err = r.cli.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating deployment", "namespace", desired.Namespace, "name", desired.Name)

		if err := r.cli.Create(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to create the deployment",
				"deployment", desired.Name)
			return err
		}
	} else if err != nil {
		r.log.Error(err, "failed to get the expected deployment",
			"deployment", desired.Name)
		return err
	}

	// TODO(gaocegege): Update status.
	return nil
}
