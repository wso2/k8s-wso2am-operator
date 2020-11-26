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

package pattern2

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func PubDevTm1Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-1",
	}

	pubdevtm1ports := getPubDevTmSpecificSvcPorts()
	servType := "ClusterIP"

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-pattern-2-am-1-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports:    pubdevtm1ports,
		},
	}
}

// for handling pub-dev-tm-2 service
func PubDevTm2Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-2",
	}
	pubdevtm2ports := getPubDevTmSpecificSvcPorts()
	servType := "ClusterIP"

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-pattern-2-am-2-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports:    pubdevtm2ports,
		},
	}
}

func GatewayService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am-gateway",
	}
	gatewayports := getGatewaySpecificSvcPorts()
	servType := "ClusterIP"

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-pattern-2-am-gateway-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports:    gatewayports,
		},
	}
}

func KeyManagerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-km-service",
	}
	keymanagerports := getKeyManagerSpecificSvcPorts()
	servType := "ClusterIP"

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-pattern-2-km-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(servType),
			Ports:    keymanagerports,
		},
	}
}

// for handling analytics-dashboard service
func DashboardService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am-analytics-dashboard",
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
		servType = "ClusterIP"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-pattern2--am-analytics-dashboard-service",
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
	workerports := getWorkerPorts()
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
			Ports:    workerports,
		},
	}
}

//for common service
func PubDevTmCommonService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
	}

	pubdevtmcommonsvsports := []corev1.ServicePort{}
	servType := ""
	if apimanager.Spec.Service.Type == "NodePort" {
		pubdevtmcommonsvsports = getPubDevTmCommonSvcNPorts()
		servType = "NodePort"
	} else if apimanager.Spec.Service.Type == "LoadBalancer" {
		pubdevtmcommonsvsports = getPubDevTmCommonSvcPorts()
		servType = "LoadBalancer"
	} else if apimanager.Spec.Service.Type == "ClusterIP" {
		pubdevtmcommonsvsports = getPubDevTmCommonSvcPorts()
		servType = "ClusterIP"
	} else {
		pubdevtmcommonsvsports = getPubDevTmCommonSvcPorts()
		servType = "ClusterIP"
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
			Ports:    pubdevtmcommonsvsports,
		},
	}
}
