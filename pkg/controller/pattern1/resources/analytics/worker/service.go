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

package worker

import "k8s.io/apimachinery/pkg/util/intstr"
import (
	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	//apimv1alpha1 "github.com/apim/pkg/apis/apim/v1alpha1"
)

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
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
//
//package worker
//
//
//import "k8s.io/apimachinery/pkg/util/intstr"
//import (
//	//appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	//"k8s.io/apimachinery/pkg/api/resource"
//	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
//	//apimv1alpha1 "github.com/apim/pkg/apis/apim/v1alpha1"
//)
//
//// newService creates a new Service for a Apimanager resource.
//// It expose the service with Nodeport type with minikube ip as the externel ip.
//func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-worker",
//	}
//	return &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.ServiceName,
//			Name:      "wso2apim-analytics-service",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: corev1.ServiceSpec{
//			Selector: labels,
//			Type:     "ClusterIP",
//			//Type:"NodePort",
//			// values are fetched from wso2-apim.yaml file
//			// Type: apimanager.Spec.ServType,
//			//ExternalIPs: []string{"192.168.99.101"},
//			// ExternalIPs: apimanager.Spec.ExternalIps,
//			Ports: []corev1.ServicePort{
//				{
//					Name:       "thrift",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7612,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
//					//NodePort:   31031,
//				},
//				{
//					Name:       "thrift-ssl",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7712,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
//					//NodePort:   31032,
//				},
//				{
//					Name:       "rest-api-port-1",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
//					//NodePort:   31033,
//				},
//				{
//					Name:       "rest-api-port-2",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9091,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
//					//NodePort:   31034,
//				},
//				{
//					Name:       "rest-api-port-3",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7071,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
//					//NodePort:   31035,
//				},
//				{
//					Name:       "rest-api-port-4",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
//					//NodePort:   31036,
//				},
//				{
//					Name:       "rest-api-port-2",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7575,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7575},
//					//NodePort:   31037,
//				},
//				{
//					Name:       "rest-api-port-3",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7576,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7576},
//					//NodePort:   31038,
//				},
//				{
//					Name:       "rest-api-port-4",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7577,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7577},
//					//NodePort:   31039,
//				},
//			},
//		},
//	}
//}
//
//





















//package worker
//
//import "k8s.io/apimachinery/pkg/util/intstr"
//import (
//	//appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	//"k8s.io/apimachinery/pkg/api/resource"
//	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
//	//apimv1alpha1 "github.com/apim/pkg/apis/apim/v1alpha1"
//)
//
//// newService creates a new Service for a Apimanager resource.
//// It expose the service with Nodeport type with minikube ip as the externel ip.
//func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-worker",
//	}
//	return &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.ServiceName,
//			Name:      "wso2apim-analytics-service",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: corev1.ServiceSpec{
//			Selector: labels,
//			Type:     "NodePort",
//			// values are fetched from wso2-apim.yaml file
//			// Type: apimanager.Spec.ServType,
//			ExternalIPs: []string{"192.168.99.101"},
//			// ExternalIPs: apimanager.Spec.ExternalIps,
//			Ports: []corev1.ServicePort{
//				{
//					Name:       "thrift",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7612,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
//					NodePort:   32020,
//				},
//				{
//					Name:       "thrift-ssl",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7712,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
//					NodePort:   32011,
//				},
//				{
//					Name:       "rest-api-port-1",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
//					NodePort:   32012,
//				},
//				{
//					Name:       "rest-api-port-2",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9091,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
//					NodePort:   32013,
//				},
//				{
//					Name:       "rest-api-port-3",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7071,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
//					NodePort:   32014,
//				},
//				{
//					Name:       "rest-api-port-4",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
//					NodePort:   32015,
//				},
//			},
//		},
//	}
//}
////
////package worker
////
////
////import "k8s.io/apimachinery/pkg/util/intstr"
////import (
////	//appsv1 "k8s.io/api/apps/v1"
////	corev1 "k8s.io/api/core/v1"
////	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
////	//"k8s.io/apimachinery/pkg/api/resource"
////	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
////	//apimv1alpha1 "github.com/apim/pkg/apis/apim/v1alpha1"
////)
////
////// newService creates a new Service for a Apimanager resource.
////// It expose the service with Nodeport type with minikube ip as the externel ip.
////func WorkerService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
////	labels := map[string]string{
////		"deployment": "wso2am-pattern-1-analytics-worker",
////	}
////	return &corev1.Service{
////		ObjectMeta: metav1.ObjectMeta{
////			// Name: apimanager.Spec.ServiceName,
////			Name:      "wso2apim-analytics-service",
////			Namespace: apimanager.Namespace,
////			OwnerReferences: []metav1.OwnerReference{
////				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
////			},
////		},
////		Spec: corev1.ServiceSpec{
////			Selector: labels,
////			Type:     "ClusterIP",
////			//Type:"NodePort",
////			// values are fetched from wso2-apim.yaml file
////			// Type: apimanager.Spec.ServType,
////			//ExternalIPs: []string{"192.168.99.101"},
////			// ExternalIPs: apimanager.Spec.ExternalIps,
////			Ports: []corev1.ServicePort{
////				{
////					Name:       "thrift",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7612,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
////					//NodePort:   31031,
////				},
////				{
////					Name:       "thrift-ssl",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7712,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
////					//NodePort:   31032,
////				},
////				{
////					Name:       "rest-api-port-1",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       9444,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
////					//NodePort:   31033,
////				},
////				{
////					Name:       "rest-api-port-2",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       9091,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
////					//NodePort:   31034,
////				},
////				{
////					Name:       "rest-api-port-3",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7071,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
////					//NodePort:   31035,
////				},
////				{
////					Name:       "rest-api-port-4",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7444,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
////					//NodePort:   31036,
////				},
////				{
////					Name:       "rest-api-port-2",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7575,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7575},
////					//NodePort:   31037,
////				},
////				{
////					Name:       "rest-api-port-3",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7576,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7576},
////					//NodePort:   31038,
////				},
////				{
////					Name:       "rest-api-port-4",
////					Protocol:   corev1.ProtocolTCP,
////					Port:       7577,
////					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7577},
////					//NodePort:   31039,
////				},
////			},
////		},
////	}
////}
////
////