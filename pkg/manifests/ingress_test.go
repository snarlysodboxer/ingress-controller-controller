package manifests

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetAnnotatedServices(t *testing.T) {
	input := inputList()
	result := GetAnnotatedServices(input)
	expected := expectedList()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}
}

func inputList() corev1.ServiceList {
	return corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceList",
			APIVersion: "core/v1",
		},
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "core/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "default",
					Annotations: map[string]string{
						"github.com/snarlysodboxer/ingress-controller-controller": "fdsa", // TODO set as constant
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port:     80,
							NodePort: 80,
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "core/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "another",
					Namespace: "default",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port:     80,
							NodePort: 80,
						},
					},
				},
			},
		},
	}
}

func expectedList() corev1.ServiceList {
	return corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceList",
			APIVersion: "core/v1",
		},
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "core/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "default",
					Annotations: map[string]string{
						"github.com/snarlysodboxer/ingress-controller-controller": "fdsa", // TODO set as constant
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port:     80,
							NodePort: 80,
						},
					},
				},
			},
		},
	}
}
