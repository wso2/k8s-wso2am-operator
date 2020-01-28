/*
 *
 *  * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *  *
 *  * WSO2 Inc. licenses this file to you under the Apache License,
 *  * Version 2.0 (the "License"); you may not use this file except
 *  * in compliance with the License.external
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

package pattern1

import (
	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func Apim1Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}

	//apimsvs1ports:= []corev1.ServicePort{}
	apimsvs1ports := getApimSvcCIPorts()
	servType :="ClusterIP"
	//if apimanager.Spec.Service.Type=="NodePort"{
	//	apimsvs1ports = getApim1SvcNPPorts()
	//	servType = "NodePort"
	//} else{
	//	apimsvs1ports = getApim1SvcLBPorts()
	//	servType = "LoadBalancer"
	//}



	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-1-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:    corev1.ServiceType(servType),
			Ports: apimsvs1ports,

		},
	}
}

// for handling apim-instance-2 service
func Apim2Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-2",
	}
	apimsvs2ports := getApimSvcCIPorts()
	servType :="ClusterIP"
	//apimsvs2ports:= []corev1.ServicePort{}
	//servType :=""
	//if apimanager.Spec.Service.Type=="NodePort"{
	//	apimsvs2ports = getApim2SvcNPPorts()
	//	servType = "NodePort"
	//} else {
	//	apimsvs2ports = getApim2SvcLBPorts()
	//	servType = "LoadBalancer"
	//}


	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-2-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports: apimsvs2ports,

		},
	}
}

// for handling analytics-dashboard service
func DashboardService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	servType :=""
	dashports := []corev1.ServicePort{}
	if apimanager.Spec.Service.Type=="NodePort"{
		dashports = getDashNPPorts()
		servType = "NodePort"
	} else if apimanager.Spec.Service.Type=="LoadBalancer"{
		dashports = getDashLBPorts()
		servType = "LoadBalancer"
	} else if apimanager.Spec.Service.Type=="ClusterIP" {
		dashports = getDashCIPorts()
		servType = "ClusterIP"
	} else {
		dashports = getDashLBPorts()
		servType = "LoadBalancer"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-dashboard-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:    corev1.ServiceType(servType),
			Ports: 	dashports,

		},
	}
}

// for handling analytics-worker service
func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	servType :=""
	if apimanager.Spec.Service.Type=="NodePort"{
		servType = "NodePort"
	} else if apimanager.Spec.Service.Type=="LoadBalancer"{
		servType = "LoadBalancer"
	} else {
		servType = "LoadBalancer"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-worker-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports: []corev1.ServicePort{
				{
					Name:       "thrift",
					Protocol:   corev1.ProtocolTCP,
					Port:       7612,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
					//NodePort:   32020,
				},
				{
					Name:       "thrift-ssl",
					Protocol:   corev1.ProtocolTCP,
					Port:       7712,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
					//NodePort:   32011,
				},
				{
					Name:       "rest-api-port-1",
					Protocol:   corev1.ProtocolTCP,
					Port:       9444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
					//NodePort:   32012,
				},
				{
					Name:       "rest-api-port-2",
					Protocol:   corev1.ProtocolTCP,
					Port:       9091,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
					//NodePort:   32013,
				},
				{
					Name:       "rest-api-port-3",
					Protocol:   corev1.ProtocolTCP,
					Port:       7071,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
					//NodePort:   32014,
				},
				{
					Name:       "rest-api-port-4",
					Protocol:   corev1.ProtocolTCP,
					Port:       7444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
					//NodePort:   32015,
				},
			},
		},
	}
}

//for common service
func ApimCommonService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
	}

	apimcommonsvsports:= []corev1.ServicePort{}
	servType :=""
	if apimanager.Spec.Service.Type=="NodePort"{
		apimcommonsvsports = getApimCommonSvcNPPorts()
		servType = "NodePort"
	}else if apimanager.Spec.Service.Type=="LoadBalancer"{
		apimcommonsvsports = getApimCommonSvcLBPorts()
		servType = "LoadBalancer"
	}else if apimanager.Spec.Service.Type=="ClusterIP"{
		apimcommonsvsports = getApimCommonSvcCIPorts()
		servType = "ClusterIP"
	} else{
		apimcommonsvsports = getApimCommonSvcLBPorts()
		servType = "LoadBalancer"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-svc",
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


