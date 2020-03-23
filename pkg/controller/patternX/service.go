/*
 *
 *  * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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

package patternX

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func ApimXService(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile) *corev1.Service {
	labels := map[string]string{
		"deployment": r.Name,
	}

	apimcommonsvsports := []corev1.ServicePort{}
	servType := ""
	if apimanager.Spec.Service.Type == "NodePort"{
		apimcommonsvsports = getApimCommonSvcNPPorts()
		servType = "NodePort"
	}else if apimanager.Spec.Service.Type == "LoadBalancer"{
		apimcommonsvsports = getApimCommonSvcPorts()
		servType = "LoadBalancer"
	}else if apimanager.Spec.Service.Type == "ClusterIP"{
		apimcommonsvsports = getApimCommonSvcPorts()
		servType = "ClusterIP"
	} else{
		apimcommonsvsports = getApimCommonSvcPorts()
		servType = "LoadBalancer"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Service.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:    corev1.ServiceType(servType),
			Ports: apimcommonsvsports,
		},
	}
}

func DashboardXService(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile) *corev1.Service {
	labels := map[string]string{
		"deployment": r.Name,
	}


	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Service.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:   "LoadBalancer",
			Ports: 	[]corev1.ServicePort{
				{
					Name:       "analytics-dashboard",
					Protocol:   corev1.ProtocolTCP,
					Port:       9643,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},

				},
			},
		},
	}
}

func WorkerXService(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile) *corev1.Service {
	labels := map[string]string{
		"deployment": r.Name,
	}


	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:     r.Service.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "LoadBalancer",
			Ports: []corev1.ServicePort{
				{
					Name:       "thrift",
					Protocol:   corev1.ProtocolTCP,
					Port:       7612,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},

				},
				{
					Name:       "thrift-ssl",
					Protocol:   corev1.ProtocolTCP,
					Port:       7712,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},

				},
				{
					Name:       "rest-api-port-1",
					Protocol:   corev1.ProtocolTCP,
					Port:       9444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},

				},
				{
					Name:       "rest-api-port-2",
					Protocol:   corev1.ProtocolTCP,
					Port:       9091,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},

				},
				{
					Name:       "rest-api-port-3",
					Protocol:   corev1.ProtocolTCP,
					Port:       7071,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},

				},
				{
					Name:       "rest-api-port-4",
					Protocol:   corev1.ProtocolTCP,
					Port:       7444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},

				},
			},
		},
	}
}
