package manifests

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

// List all `Service` objects
func GetAllServices() (error, corev1.ServiceList) {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},
	}
	err := sdk.List("default", &serviceList) // TODO set the namespace via config
	if err != nil {
		logrus.Errorf("Failed to query Services : %v", err)
		return err, corev1.ServiceList{}
	}

	return nil, serviceList
}

func GetAnnotatedServices(sl corev1.ServiceList) corev1.ServiceList {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceList",
			APIVersion: "core/v1",
		},
	}
	for _, service := range sl.Items {
		if service.ObjectMeta.Annotations["github.com/snarlysodboxer/ingress-controller-controller"] != "" { // TODO how to set?
			serviceList.Items = append(serviceList.Items, service)
		}
	}

	return serviceList
}
