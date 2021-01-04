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

package pattern3

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//GatewayService is for handling gateway-sevice...
func GatewayService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-gateway",
	}
	gatewayports := getGatewaySpecificSvcPorts()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-gw-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    gatewayports,
		},
	}
}

//KeyManagerService is for handling key manager service...
func KeyManagerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-km",
	}
	keymanagerports := getKeyManagerSpecificSvcPorts()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-km-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    keymanagerports,
		},
	}
}

//TrafficManagerService is for handling key manager service...
func TrafficManagerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-tm",
	}
	tmports := getTrafficManagerSpecificSvcPorts()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-tm-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    tmports,
		},
	}
}

//TrafficManagerHeadlessService is for handling key manager service...
func TrafficManagerHeadlessService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-tm",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-tm-headless-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector:  labels,
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       "binary",
					Protocol:   corev1.ProtocolTCP,
					Port:       9611,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
				},
				{
					Name:       "binary-secure",
					Protocol:   corev1.ProtocolTCP,
					Port:       9711,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
				},
				{
					Name:       "jms-tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       5672,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
				},
			},
		},
	}
}

// DashboardService for handling analytics-dashboard service...
func DashboardService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-analytics-dashboard",
	}
	dashports := getDashBoardPorts()

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
			Ports:    dashports,
		},
	}
}

//WorkerService is for handling analytics-worker service
func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-analytics-worker",
	}
	workerports := getWorkerPorts()

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
			Ports:    workerports,
		},
	}
}

//WorkerHeadlessService is for...
func WorkerHeadlessService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-analytics-worker",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-worker-headless-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector:  labels,
			ClusterIP: "None",
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
			},
		},
	}
}

//PubCommonService is for common service
func PubCommonService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-am-publisher",
	}

	pubcommonsvsports := getPubCommonSvcPorts()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-publisher-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    pubcommonsvsports,
		},
	}
}

//DevCommonService is for common service
func DevCommonService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2-am-devportal",
	}

	devcommonsvsports := getDevportalCommonSvcPorts()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-devportal-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    devcommonsvsports,
		},
	}
}
