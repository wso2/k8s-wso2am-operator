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
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func Apim1Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node":       "wso2am-pattern-1-am-1",
	}

	apimsvs1ports := getApimSpecificSvcPorts()
	servType := "ClusterIP"

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
			Type:     corev1.ServiceType(servType),
			Ports:    apimsvs1ports,
		},
	}
}

// for handling apim-instance-2 service
func Apim2Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node":       "wso2am-pattern-1-am-2",
	}
	apimsvs2ports := getApimSpecificSvcPorts()
	servType := "ClusterIP"

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
			Ports:    apimsvs2ports,
		},
	}
}

// for handling analytics-dashboard service
func DashboardService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	servType := ""
	dashports := []corev1.ServicePort{}
	if apimanager.Spec.Service.Type == "NodePort" {
		dashports = getDashBoardNPPorts()
		servType = "NodePort"
	} else if apimanager.Spec.Service.Type == "LoadBalancer" {
		dashports = getDashBoardPorts()
		servType = "LoadBalancer"
	} else if apimanager.Spec.Service.Type == "ClusterIP" {
		dashports = getDashBoardPorts()
		servType = "ClusterIP"
	} else {
		dashports = getDashBoardPorts()
		servType = "LoadBalancer"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-dashboard-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports:    dashports,
		},
	}
}

// for handling analytics-worker service
func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	servType := "ClusterIP"

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-worker-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
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

//for common service
func ApimCommonService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
	}

	apimcommonsvsports := []corev1.ServicePort{}
	servType := ""
	if apimanager.Spec.Service.Type == "NodePort" {
		apimcommonsvsports = getApimCommonSvcNPPorts()
		servType = "NodePort"
	} else if apimanager.Spec.Service.Type == "LoadBalancer" {
		apimcommonsvsports = getApimCommonSvcPorts()
		servType = "LoadBalancer"
	} else if apimanager.Spec.Service.Type == "ClusterIP" {
		apimcommonsvsports = getApimCommonSvcPorts()
		servType = "ClusterIP"
	} else {
		apimcommonsvsports = getApimCommonSvcPorts()
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
			Type:     corev1.ServiceType(servType),
			Ports:    apimcommonsvsports,
		},
	}
}
