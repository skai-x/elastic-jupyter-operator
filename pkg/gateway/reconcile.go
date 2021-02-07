package gateway

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/tkestack/jupyter-operator/api/v1alpha1"
)

type Reconciler struct {
	cli      client.Client
	log      logr.Logger
	recorder record.EventRecorder
	scheme   *runtime.Scheme

	instance *v1alpha1.JupyterGateway
	gen      *generator
}

func NewReconciler(cli client.Client, l logr.Logger,
	r record.EventRecorder, s *runtime.Scheme,
	i *v1alpha1.JupyterGateway) (*Reconciler, error) {
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
	serviceAccountName, err := r.reconcileRBAC()
	if err != nil {
		return err
	}
	if err := r.reconcileDeployment(serviceAccountName); err != nil {
		return err
	}
	if err := r.reconcileService(); err != nil {
		return err
	}
	return nil
}

func (r Reconciler) reconcileRBAC() (string, error) {
	sa, err := r.reconcileServiceAccount()
	if err != nil {
		return "", err
	}
	if err := r.reconcileRoleBinding(sa); err != nil {
		return "", err
	}
	return sa.Name, nil
}

func (r Reconciler) reconcileRoleBinding(
	sa *v1.ServiceAccount) error {
	desired := r.gen.DesiredRoleBinding(sa)

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return err
	}

	actual := &rbacv1.RoleBinding{}
	err := r.cli.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating rolebinding",
			"namespace", desired.Namespace, "name", desired.Name)

		if err := r.cli.Create(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to create the rolebinding",
				"rolebinding", desired.Name)
			return err
		}
	} else if err != nil {
		r.log.Error(err, "failed to get the expected rolebinding",
			"rolebinding", desired.Name)
		return err
	}
	return nil
}

func (r Reconciler) reconcileServiceAccount() (*v1.ServiceAccount, error) {
	desired := r.gen.DesiredServiceAccountWithoutOwner()

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return nil, err
	}

	actual := &v1.ServiceAccount{}
	err := r.cli.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating serviceaccount", "namespace", desired.Namespace, "name", desired.Name)

		if err := r.cli.Create(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to create the serviceaccount",
				"serviceaccount", desired.Name)
			return nil, err
		}
	} else if err != nil {
		r.log.Error(err, "failed to get the expected serviceaccount",
			"serviceaccount", desired.Name)
		return nil, err
	}
	// When the sa is created, actual is nil. Thus actual cannot be used to build rolebinding.
	return desired, nil
}

func (r Reconciler) reconcileService() error {
	desired := r.gen.DesiredServiceWithoutOwner()

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return err
	}

	actual := &v1.Service{}
	err := r.cli.Get(context.TODO(),
		types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, actual)
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating service", "namespace", desired.Namespace, "name", desired.Name)

		if err := r.cli.Create(context.TODO(), desired); err != nil {
			r.log.Error(err, "Failed to create the serivce",
				"service", desired.Name)
			return err
		}
	} else if err != nil {
		r.log.Error(err, "failed to get the expected service",
			"service", desired.Name)
		return err
	}
	return nil
}

func (r Reconciler) reconcileDeployment(sa string) error {
	desired := r.gen.DesiredDeploymentWithoutOwner(sa)

	if err := controllerutil.SetControllerReference(
		r.instance, desired, r.scheme); err != nil {
		r.log.Error(err,
			"Set controller reference error, requeuing the request")
		return err
	}

	actual := &appsv1.Deployment{}
	err := r.cli.Get(context.TODO(),
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

	if !equality.Semantic.DeepEqual(r.instance.Status.DeploymentStatus, actual.Status) {
		r.instance.Status.DeploymentStatus = actual.Status
		if err := r.cli.Status().Update(context.TODO(), r.instance); err != nil {
			r.log.Error(err, "failed to update status",
				"namespace", r.instance.Namespace,
				"jupytergateway", r.instance.Name)
			return err
		}
	}
	return nil
}
