package pattern1

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	networkv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApimIngress(apimanager *apimv1alpha1.APIManager) *v1beta1.Ingress  {

	//labels := map[string]string{
	//	"deployment": "wso2am-pattern-1-am",
	//}
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:     "wso2-am-p1-ingress",
			Namespace: apimanager.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/affinity": "cookie",
				"nginx.ingress.kubernetes.io/session-cookie-name": "route",
				"nginx.ingress.kubernetes.io/session-cookie-hash": "sha1",
			},
			//Labels:   labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: "wso2apim",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkv1.IngressBackend{
										ServiceName: "wso2-am-svc",
										ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
									},
								},
							},
						},
					},
				},
			},
			TLS: []v1beta1.IngressTLS{
				{
					Hosts:[]string{
						"wso2apim",
					},
				},
			},
		},
	}
}

func GatewayIngress(apimanager *apimv1alpha1.APIManager) *v1beta1.Ingress  {

	//labels := map[string]string{
	//	"deployment": "wso2am-pattern-1-am",
	//}
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:     "wso2-am-gateway-p1-ingress",
			Namespace: apimanager.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
			},
			//Labels:   labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: "wso2apim-gateway",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkv1.IngressBackend{
										ServiceName: "wso2-am-svc",
										ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
									},
								},
							},
						},
					},
				},
			},
			TLS: []v1beta1.IngressTLS{
				{
					Hosts:[]string{
						"wso2apim-gateway",
					},
				},
			},
		},
	}
}

func DashboardIngress(apimanager *apimv1alpha1.APIManager) *v1beta1.Ingress  {

	//labels := map[string]string{
	//	"deployment": "wso2am-pattern-1-am",
	//}
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:     "wso2-am-analytics-dashboard-p1-ingress",
			Namespace: apimanager.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
			},
			//Labels:   labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: "wso2apim-analytics-dashboard",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkv1.IngressBackend{
										ServiceName: "wso2-am-analytics-dashboard-svc",
										ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
									},
								},
							},
						},
					},
				},
			},
			TLS: []v1beta1.IngressTLS{
				{
					Hosts:[]string{
						"wso2apim-analytics-dashboard",
					},
				},
			},
		},
	}
}
