package manifests

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	prodConfig = `name: production
host: this.example.com
path: /*
service: web
port: 80`
	stagingConfig = `name: staging
host: that.example.com
path: /that
service: web
port: 80`
	stagingConfig2 = `name: staging
host: other.example.com
path: /asdf
service: web
port: 80`
	stagingConfig3 = `name: staging
host: other.example.com
path: /fdsa
service: web2
port: 80`
)

func TestGetAnnotatedServices(t *testing.T) {
	serviceList := newServiceList()
	result := GetAnnotatedServices(serviceList)
	expected := expectedServiceList()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}
}

func TestBuildConfigs(t *testing.T) {
	serviceList := newServiceList()
	annotatedList := GetAnnotatedServices(serviceList)
	err, result := BuildConfigs(annotatedList)
	if err != nil {
		t.Errorf("Error building ingress configs: %v\n", err)
	}
	expected := expectedIngressConfigs()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}
}

func TestNewIngressList(t *testing.T) {
	serviceList := newServiceList()
	annotatedList := GetAnnotatedServices(serviceList)
	err, configs := BuildConfigs(annotatedList)
	if err != nil {
		t.Errorf("Error building ingress configs: %v\n", err)
	}
	result := NewIngressList(configs)
	expected := expectedIngressList()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}
}

func TestGetAnnotatedIngresses(t *testing.T) {
	ingressList := newIngressList()
	result := GetAnnotatedIngresses(ingressList)
	expected := expectedIngressList()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}
}

func TestGetOrphanedIngresses(t *testing.T) {
	// list with more observed than desired
	ingressList := newIngressList()
	desired := GetAnnotatedIngresses(ingressList)
	observed := desired
	ai := aIngress()
	expected := v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Items: []v1beta1.Ingress{ai},
	}
	observed.Items = append(observed.Items, ai)
	result := GetOrphanedIngresses(desired, observed)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v\n", expected, result)
	}

	// equal lists
	expected = v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
	}
	result = GetOrphanedIngresses(desired, desired)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, result)
	}

	// list with more desired than observed
	ingressList = newIngressList()
	desired = GetAnnotatedIngresses(ingressList)
	observed = desired
	ai = aIngress()
	expected = v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
	}
	desired.Items = append(desired.Items, ai)
	result = GetOrphanedIngresses(desired, observed)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v\n", expected, result)
	}
}

func expectedIngressConfigs() []ingressConfig {
	return []ingressConfig{
		{
			Name: "production",
			HostConfigs: []hostConfig{
				{
					Host: "this.example.com",
					PathConfigs: []pathConfig{
						{
							Path:    "/*",
							Service: "web",
							Port:    80,
						},
					},
				},
			},
		},
		{
			Name: "staging",
			HostConfigs: []hostConfig{
				{
					Host: "that.example.com",
					PathConfigs: []pathConfig{
						{
							Path:    "/that",
							Service: "web",
							Port:    80,
						},
					},
				},
				{
					Host: "other.example.com",
					PathConfigs: []pathConfig{
						{
							Path:    "/asdf",
							Service: "web",
							Port:    80,
						},
						{
							Path:    "/fdsa",
							Service: "web2",
							Port:    80,
						},
					},
				},
			},
		},
	}
}

func newServiceList() corev1.ServiceList {
	return corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: prodConfig,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "two",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig2,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "three",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig3,
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
					APIVersion: "v1",
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

func expectedServiceList() corev1.ServiceList {
	return corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: prodConfig,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "two",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig2,
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
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "three",
					Namespace: "default",
					Annotations: map[string]string{
						configAnnotationKey: stagingConfig3,
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

func expectedIngressList() v1beta1.IngressList {
	servicePort := intstr.FromInt(80)
	return v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Items: []v1beta1.Ingress{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "extensions/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "production",
					Namespace: "default",
					Annotations: map[string]string{
						ingressAnnotationKey: "true",
					},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "this.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/*",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "extensions/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "staging",
					Namespace: "default",
					Annotations: map[string]string{
						ingressAnnotationKey: "true",
					},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "that.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/that",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
						{
							Host: "other.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/asdf",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
										{
											Path: "/fdsa",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web2",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func newIngressList() v1beta1.IngressList {
	servicePort := intstr.FromInt(80)
	return v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		Items: []v1beta1.Ingress{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "extensions/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "production",
					Namespace: "default",
					Annotations: map[string]string{
						ingressAnnotationKey: "true",
					},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "this.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/*",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "extensions/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "staging",
					Namespace: "default",
					Annotations: map[string]string{
						ingressAnnotationKey: "true",
					},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "that.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/that",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
						{
							Host: "other.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/asdf",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
										{
											Path: "/fdsa",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web2",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "extensions/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-other",
					Namespace: "default",
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "nope.example.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/*",
											Backend: v1beta1.IngressBackend{
												ServiceName: "web",
												ServicePort: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func aIngress() v1beta1.Ingress {
	servicePort := intstr.FromInt(80)
	return v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "a-ingress",
			Namespace: "default",
			Annotations: map[string]string{
				ingressAnnotationKey: "true",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: "a-ingress.example.com",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/*",
									Backend: v1beta1.IngressBackend{
										ServiceName: "web",
										ServicePort: servicePort,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
