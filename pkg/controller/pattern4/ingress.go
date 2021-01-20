/*
 *
 *  * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *  *
 *  * WSO2 Inc. licenses this file to you under the Apache License,
 *  * Version 2.0 (the "License"); you may not use this file except
 *  * in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  * http:www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing,
 *  * software distributed under the License is distributed on an
 *  * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  * KIND, either express or implied. See the License for the
 *  * specific language governing permissions and limitations
 *  * under the License.
 *
 */

package pattern4

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	networkv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func PubDevTmIngress(apimanager *apimv1alpha1.APIManager, x *ingressConfigvalues) *v1beta1.Ingress {

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        x.IngressName,
			Namespace:   apimanager.Namespace,
			Annotations: x.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: x.Hostname,
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
					Hosts: []string{
						x.Hostname,
					},
				},
			},
		},
	}
}

func ExternalGatewayIngress(apimanager *apimv1alpha1.APIManager, x *ingressConfigvalues) *v1beta1.Ingress {

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        x.IngressName,
			Namespace:   apimanager.Namespace,
			Annotations: x.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: x.Hostname,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkv1.IngressBackend{
										ServiceName: "wso2-am-external-gw-svc",
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
					Hosts: []string{
						x.Hostname,
					},
				},
			},
		},
	}
}

//Internal Gateway
func InternalGatewayIngress(apimanager *apimv1alpha1.APIManager, x *ingressConfigvalues) *v1beta1.Ingress {

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        x.IngressName,
			Namespace:   apimanager.Namespace,
			Annotations: x.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: x.Hostname,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkv1.IngressBackend{
										ServiceName: "wso2-am-internal-gw-svc",
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
					Hosts: []string{
						x.Hostname,
					},
				},
			},
		},
	}
}

func DashboardIngress(apimanager *apimv1alpha1.APIManager, x *ingressConfigvalues) *v1beta1.Ingress {

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        x.IngressName,
			Namespace:   apimanager.Namespace,
			Annotations: x.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: x.Hostname,
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
					Hosts: []string{
						x.Hostname,
					},
				},
			},
		},
	}
}
