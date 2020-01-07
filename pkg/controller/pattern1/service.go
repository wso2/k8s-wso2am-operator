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
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "apim-1-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			ExternalIPs: []string{"192.168.99.101"},
			Ports: []corev1.ServicePort{
				{
					Name:       "pass-through-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8280,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
					NodePort:   32004,
					//NodePort:9611,
				},
				{
					Name:       "pass-through-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       8243,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
					NodePort:   32003,
					//NodePort:   9711,

				},
				{
					Name:       "servlet-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       9763,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
					NodePort:   32002,
					//NodePort:   5672,

				},
				{
					Name:       "servlet-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       9443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
					NodePort:   32001,
					//NodePort:   9443,

				},
				{
					Name:       "jms-provider",
					Protocol:   corev1.ProtocolTCP,
					Port:       5672,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
					NodePort:   32005,
					//NodePort:   5672,

				},
				//{
				//	Name:       "binary",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9611,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
				//	NodePort:   32002,
				//},
				//{
				//	Name:       "binary-secure",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9711,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
				//	NodePort:   32003,
				//},
				//{
				//	Name:       "jms-tcp",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       5672,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
				//	NodePort:   32004,
				//},
				//{
				//	Name:       "servlet-https",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9443,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
				//	NodePort:   32001,
				//},
			},
		},
	}
}

// for handling apim-instance-2 service
func Apim2Service(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-2",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "apim-2-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			ExternalIPs: []string{"192.168.99.101"},
			Ports: []corev1.ServicePort{
				{
					Name:       "pass-through-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       9611,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
					NodePort:   32009,
				},
				{
					Name:       "pass-through-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       9711,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
					NodePort:   32008,
				},
				{
					Name:       "servlet-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       5672,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
					NodePort:   32007,
				},
				{
					Name:       "servlet-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       9443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
					NodePort:   32006,
				},
				//{
				//	Name:       "jms-provider",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       5672,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
				//	NodePort:   32010,
				//},
				//{
				//	Name:       "binary",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9611,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
				//	NodePort:   32002,
				//},
				//{
				//	Name:       "binary-secure",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9711,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
				//	NodePort:   32003,
				//},
				//{
				//	Name:       "jms-tcp",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       5672,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
				//	NodePort:   32004,
				//},
				//{
				//	Name:       "servlet-https",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9443,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
				//	NodePort:   32001,
				//},
			},
		},
	}
}

// for handling analytics-dashboard service
func DashboardService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "analytics-dash-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			// values are fetched from wso2-apim.yaml file
			// Type: apimanager.Spec.ServType,
			ExternalIPs: []string{"192.168.99.101"},
			// ExternalIPs: apimanager.Spec.ExternalIps,
			Ports: []corev1.ServicePort{
				//{
				//	Name:       "analytics-dashboard",
				//	Protocol:   corev1.ProtocolTCP,
				//	Port:       9643,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
				//	NodePort:   32009,
				//},
				{
					Name:       "analytics-dashboard",
					Protocol:   corev1.ProtocolTCP,
					Port:       32201,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32201},
					NodePort:   32201,
				},
			},
		},
	}
}

// for handling analytics-worker service
func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "wso2apim-analytics-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			// values are fetched from wso2-apim.yaml file
			// Type: apimanager.Spec.ServType,
			ExternalIPs: []string{"192.168.99.101"},
			// ExternalIPs: apimanager.Spec.ExternalIps,
			Ports: []corev1.ServicePort{
				{
					Name:       "thrift",
					Protocol:   corev1.ProtocolTCP,
					Port:       7612,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
					NodePort:   32020,
				},
				{
					Name:       "thrift-ssl",
					Protocol:   corev1.ProtocolTCP,
					Port:       7712,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
					NodePort:   32011,
				},
				{
					Name:       "rest-api-port-1",
					Protocol:   corev1.ProtocolTCP,
					Port:       9444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
					NodePort:   32012,
				},
				{
					Name:       "rest-api-port-2",
					Protocol:   corev1.ProtocolTCP,
					Port:       9091,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
					NodePort:   32013,
				},
				{
					Name:       "rest-api-port-3",
					Protocol:   corev1.ProtocolTCP,
					Port:       7071,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
					NodePort:   32014,
				},
				{
					Name:       "rest-api-port-4",
					Protocol:   corev1.ProtocolTCP,
					Port:       7444,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
					NodePort:   32015,
				},
			},
		},
	}
}

// // for handling mysql service
func MysqlService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2apim-with-analytics-rdbms-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql-port",
					Protocol:   corev1.ProtocolTCP,
					Port:       3306,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 3306},
				},
			},
		},
	}
}
