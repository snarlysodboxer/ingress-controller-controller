package stub

import (
	"context"

	"github.com/snarlysodboxer/ingress-controller-controller/pkg/manifests"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	// corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func NewHandler(m *Metrics) sdk.Handler {
	return &Handler{
		metrics: m,
	}
}

type Metrics struct {
	operatorErrors prometheus.Counter
}

type Handler struct {
	// Metrics example
	metrics *Metrics

	// Fill me TODO
}

func (handler *Handler) Handle(ctx context.Context, event sdk.Event) error {
	if event.Deleted {
		return nil
	}

	// TODO ensure the sdk doesn't need this
	// switch object := event.Object.(type) {
	// case *corev1.Service:
	err, services := manifests.GetAllServices()
	if err != nil {
		logrus.Errorf("Error listing Services: %v", err)
	}

	annotatedServices := manifests.GetAnnotatedServices(services)

	err, configs := manifests.BuildConfigs(annotatedServices)
	if err != nil {
		logrus.Errorf("Error building ingress configs: %v\n", err)
	}

	calculatedIngresses := manifests.NewIngressList(configs)

	err, ingresses := manifests.GetAllIngresses()
	if err != nil {
		logrus.Errorf("Error listing Ingresses: %v", err)
	}

	annotatedIngresses := manifests.GetAnnotatedIngresses(ingresses)

	orphans := manifests.GetOrphanedIngresses(calculatedIngresses, annotatedIngresses)
	err = sdk.Delete(&orphans)
	if err != nil {
		logrus.Errorf("Error deleting Ingresses: %v", err)
	}

	err = applyObject(handler, &calculatedIngresses)
	if err != nil {
		logrus.Errorf("Error applying Ingresses: %v", err)
	}

	logrus.Debug("Handled event loop for Service")
	// }

	return nil
}

func applyObject(handler *Handler, obj sdk.Object) error {
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	name, _, err := k8sutil.GetNameAndNamespace(obj)
	err = sdk.Create(obj)
	switch {
	case err != nil && errors.IsAlreadyExists(err):
		err = sdk.Update(obj)
		if err != nil {
			logrus.Errorf("Failed to update %s '%s' : %v", kind, name, err)
			handler.metrics.operatorErrors.Inc()
		}
	case err != nil:
		logrus.Errorf("Failed to apply %s '%s' : %v", kind, name, err)
		handler.metrics.operatorErrors.Inc()
		return err
	}
	logrus.Debugf("Reconciled %s '%s'", kind, name)

	return nil
}

func createObject(handler *Handler, obj sdk.Object) error {
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	name, _, err := k8sutil.GetNameAndNamespace(obj)
	err = sdk.Create(obj)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to apply %s '%s' : %v", kind, name, err)
		handler.metrics.operatorErrors.Inc()
		return err
	}
	logrus.Debugf("Handled event loop for %s '%s'", kind, name)

	return nil
}

func RegisterOperatorMetrics() (*Metrics, error) {
	operatorErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "icc_operator_reconcile_errors_total",
		Help: "Number of errors that occurred while reconciling Ingress manifests",
	})
	err := prometheus.Register(operatorErrors)
	if err != nil {
		return nil, err
	}

	return &Metrics{operatorErrors: operatorErrors}, nil
}
