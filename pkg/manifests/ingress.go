package manifests

import (
	"gopkg.in/yaml.v2"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const configAnnotationKey = "ingress-controller-controller.alpha.davidamick.com/config"
const ingressAnnotationKey = "ingress-controller-controller.alpha.davidamick.com/managed"

type yamlConfig struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Path    string `yaml:"path"`
	Service string `yaml:"service"`
	Port    int    `yaml:"port"`
}

type ingressConfig struct {
	Name        string
	HostConfigs []hostConfig
}

type hostConfig struct {
	Host        string
	PathConfigs []pathConfig
}

type pathConfig struct {
	Path    string
	Service string
	Port    int
}

// List all `Service` objects
func GetAllServices() (error, corev1.ServiceList) {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceList",
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

// Find the `Service`s that have the right annotation
func GetAnnotatedServices(sl corev1.ServiceList) corev1.ServiceList {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceList",
			APIVersion: "core/v1",
		},
	}
	for _, service := range sl.Items {
		if service.ObjectMeta.Annotations[configAnnotationKey] != "" { // TODO how to set?
			serviceList.Items = append(serviceList.Items, service)
		}
	}

	return serviceList
}

// NewIngressList calculates a list of `Ingress`s from the annotations
func NewIngressList(configs []ingressConfig) v1beta1.IngressList {
	ingresses := []v1beta1.Ingress{}
	for _, config := range configs {
		rules := []v1beta1.IngressRule{}
		for _, hostConfig := range config.HostConfigs {
			rule := newRule(hostConfig)
			rules = append(rules, rule)
		}
		ingress := newIngress(config.Name, rules)
		ingresses = append(ingresses, ingress)
	}

	return v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressList",
			APIVersion: "extensions/v1beta1",
		},
		Items: ingresses,
	}
}

// List all `Ingress` objects
func GetAllIngresses() (error, v1beta1.IngressList) {
	ingressList := v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressList",
			APIVersion: "core/v1",
		},
	}
	err := sdk.List("default", &ingressList) // TODO set the namespace via config
	if err != nil {
		logrus.Errorf("Failed to query Ingresses : %v", err)
		return err, v1beta1.IngressList{}
	}

	return nil, ingressList
}

func GetAnnotatedIngresses(sl v1beta1.IngressList) v1beta1.IngressList {
	ingressList := v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressList",
			APIVersion: "extensions/v1beta1",
		},
	}
	for _, ingress := range sl.Items {
		if ingress.ObjectMeta.Annotations[ingressAnnotationKey] == "true" {
			ingressList.Items = append(ingressList.Items, ingress)
		}
	}

	return ingressList
}

func GetOrphanedIngresses(desired, observed v1beta1.IngressList) v1beta1.IngressList {
	orphaned := v1beta1.IngressList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressList",
			APIVersion: "extensions/v1beta1",
		},
	}
	for _, desiredItem := range desired.Items {
		found := false
		item := v1beta1.Ingress{}
		for _, observedItem := range observed.Items {
			if desiredItem.Name == observedItem.Name {
				found = true
				item = observedItem
			}
		}
		if !found {
			orphaned.Items = append(orphaned.Items, item)
		}
	}

	return orphaned
}

//   expects all services passed to be annotated
func BuildConfigs(sl corev1.ServiceList) (error, []ingressConfig) {
	nameMap := map[string][]yamlConfig{}
	for _, service := range sl.Items {
		yc := yamlConfig{}
		err := yaml.Unmarshal([]byte(service.ObjectMeta.Annotations[configAnnotationKey]), &yc)
		if err != nil {
			return err, []ingressConfig{}
		}
		nameMap[yc.Name] = append(nameMap[yc.Name], yc)
	}

	configs := []ingressConfig{}

	for name, yConfigs := range nameMap {
		hostMap := map[string][]pathConfig{}
		for _, yConfig := range yConfigs {
			hostMap[yConfig.Host] = append(hostMap[yConfig.Host], pathConfig{
				Path:    yConfig.Path,
				Service: yConfig.Service,
				Port:    yConfig.Port,
			})
		}

		hostConfigs := []hostConfig{}
		for hostName, pathConfigs := range hostMap {
			hc := hostConfig{Host: hostName, PathConfigs: pathConfigs}
			hostConfigs = append(hostConfigs, hc)
		}

		ic := ingressConfig{
			Name:        name,
			HostConfigs: hostConfigs,
		}
		configs = append(configs, ic)
	}

	return nil, configs
}

func newIngress(name string, rules []v1beta1.IngressRule) v1beta1.Ingress {
	return v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Annotations: map[string]string{
				ingressAnnotationKey: "true",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: rules,
		},
	}
}

func newRule(config hostConfig) v1beta1.IngressRule {
	pathConfigs := []v1beta1.HTTPIngressPath{}
	for _, pathConfig := range config.PathConfigs {
		pc := v1beta1.HTTPIngressPath{
			Path: pathConfig.Path,
			Backend: v1beta1.IngressBackend{
				ServiceName: pathConfig.Service,
				ServicePort: intstr.FromInt(pathConfig.Port),
			},
		}
		pathConfigs = append(pathConfigs, pc)
	}
	return v1beta1.IngressRule{
		Host: config.Host,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: pathConfigs,
			},
		},
	}
}
