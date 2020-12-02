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

package pattern2

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//service ports

func getGatewaySpecificSvcPorts() []corev1.ServicePort {
	var gatewayPorts []corev1.ServicePort
	gatewayPorts = append(gatewayPorts, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       8280,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
	})
	gatewayPorts = append(gatewayPorts, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       8243,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
	})
	gatewayPorts = append(gatewayPorts, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
	})
	gatewayPorts = append(gatewayPorts, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return gatewayPorts
}

// func getGatewaySpecificSvcNPorts() []corev1.ServicePort {
// 	var gatewayNPorts []corev1.ServicePort
// 	gatewayNPorts = append(gatewayNPorts, corev1.ServicePort{
// 		Name:       "pass-through-http",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       8280,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
// 		NodePort:   32004,
// 	})
// 	gatewayNPorts = append(gatewayNPorts, corev1.ServicePort{
// 		Name:       "pass-through-https",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       8243,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
// 		NodePort:   32003,
// 	})
// 	gatewayNPorts = append(gatewayNPorts, corev1.ServicePort{
// 		Name:       "servlet-http",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9763,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
// 		NodePort:   32002,
// 	})
// 	gatewayNPorts = append(gatewayNPorts, corev1.ServicePort{
// 		Name:       "servlet-https",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9443,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
// 		NodePort:   32001,
// 	})
// 	return gatewayNPorts
// }

func getKeyManagerSpecificSvcPorts() []corev1.ServicePort {
	var keyManagerPorts []corev1.ServicePort
	keyManagerPorts = append(keyManagerPorts, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return keyManagerPorts
}

// func getKeyManagerSpecificSvcNPorts() []corev1.ServicePort {
// 	var keyManagerNPorts []corev1.ServicePort
// 	keyManagerNPorts = append(keyManagerNPorts, corev1.ServicePort{
// 		Name:       "servlet-https",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9443,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
// 		NodePort:   32201,
// 	})
// 	return keyManagerNPorts
// }

func getPubDevTmCommonSvcPorts() []corev1.ServicePort {
	var pubDevTmPorts []corev1.ServicePort
	pubDevTmPorts = append(pubDevTmPorts, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9703},
	})
	pubDevTmPorts = append(pubDevTmPorts, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return pubDevTmPorts
}

// func getPubDevTmCommonSvcNPorts() []corev1.ServicePort {
// 	var pubDevTmNPorts []corev1.ServicePort
// 	pubDevTmNPorts = append(pubDevTmNPorts, corev1.ServicePort{
// 		Name:       "servlet-http",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9703,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9703},
// 		NodePort:   32302,
// 	})
// 	pubDevTmNPorts = append(pubDevTmNPorts, corev1.ServicePort{
// 		Name:       "servlet-https",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9443,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
// 		NodePort:   32301,
// 	})
// 	return pubDevTmNPorts
// }

func getPubDevTmSpecificSvcPorts() []corev1.ServicePort {
	var pubDevTmPorts []corev1.ServicePort
	pubDevTmPorts = append(pubDevTmPorts, corev1.ServicePort{
		Name:       "binary",
		Protocol:   corev1.ProtocolTCP,
		Port:       9611,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
	})
	pubDevTmPorts = append(pubDevTmPorts, corev1.ServicePort{
		Name:       "binary-secure",
		Protocol:   corev1.ProtocolTCP,
		Port:       9711,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
	})
	pubDevTmPorts = append(pubDevTmPorts, corev1.ServicePort{
		Name:     "jms-tcp",
		Protocol: corev1.ProtocolTCP,
		Port:     5672,
	})
	return pubDevTmPorts
}

// func getPubDevTmSpecificSvcNPorts() []corev1.ServicePort {
// 	var pubDevTm1Ports []corev1.ServicePort
// 	pubDevTm1Ports = append(pubDevTm1Ports, corev1.ServicePort{
// 		Name:       "binary",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9611,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
// 		NodePort:   32403,
// 	})
// 	pubDevTm1Ports = append(pubDevTm1Ports, corev1.ServicePort{
// 		Name:       "binary-secure",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9711,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
// 		NodePort:   32402,
// 	})
// 	pubDevTm1Ports = append(pubDevTm1Ports, corev1.ServicePort{
// 		Name:     "jms-tcp",
// 		Protocol: corev1.ProtocolTCP,
// 		Port:     5672,
// 		NodePort: 32401,
// 	})
// 	return pubDevTm1Ports
// }

// func pubDevTm2SpecificSvcPorts() []corev1.ServicePort {
// 	var pubDevTm2Ports = []corev1.ServicePort
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "binary",
// 		Protocol:    corev1.ProtocolTCP,
// 		Port:        9611,
// 		TargetPort:  9611,
// 	})
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "binary-secure",
// 		Protocol:    corev1.ProtocolTCP,
// 		Port:        9711,
// 		TargetPort:  9711,
// 	})
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "jms-tcp",
// 		Protocol:    5672,
// 		Port:        5672,
// 	})
// 	return pubDevTm2Ports
// }

// func pubDevTm2SpecificSvcNPorts() []corev1.ServicePort {
// 	var pubDevTm2Ports = []corev1.ServicePort
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "binary",
// 		Protocol:    corev1.ProtocolTCP,
// 		Port:        9611,
// 		TargetPort:  9611,
// 		NodePort:    32503
// 	})
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "binary-secure",
// 		Protocol:    corev1.ProtocolTCP,
// 		Port:        9711,
// 		TargetPort:  9711,
// 		NodePort:    32502
// 	})
// 	pubDevTm2Ports = append(pubDevTm2Ports, corev1.ServicePort{
// 		Name:        "jms-tcp",
// 		Protocol:    5672,
// 		Port:        5672,
// 		NodePort:    32501
// 	})
// 	return pubDevTm2Ports
// }

// func getDashBoardNPPorts() []corev1.ServicePort {
// 	var dashports []corev1.ServicePort

// 	dashports = append(dashports, corev1.ServicePort{
// 		Name:       "analytics-dashboard",
// 		Protocol:   corev1.ProtocolTCP,
// 		Port:       9643,
// 		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
// 		NodePort:   32201,
// 	})

// 	return dashports
// }

func getDashBoardPorts() []corev1.ServicePort {

	var dashports []corev1.ServicePort
	dashports = append(dashports, corev1.ServicePort{
		Name:       "analytics-dashboard",
		Protocol:   corev1.ProtocolTCP,
		Port:       9643,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
	})

	return dashports
}

func getWorkerPorts() []corev1.ServicePort {
	var workerPorts []corev1.ServicePort
	workerPorts = append(workerPorts, corev1.ServicePort{
		Name:       "thrift-ssl",
		Protocol:   corev1.ProtocolTCP,
		Port:       7712,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
	})
	// workerPorts = append(workerPorts, corev1.ServicePort{
	// 	Name:       "thrift-ssl",
	// 	Protocol:   corev1.ProtocolTCP,
	// 	Port:       7712,
	// 	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712}
	// })
	workerPorts = append(workerPorts, corev1.ServicePort{
		Name:       "rest-api-port-1",
		Protocol:   corev1.ProtocolTCP,
		Port:       7444,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
	})
	return workerPorts
}

func getWorkerNPorts() []corev1.ServicePort {
	var workerPorts []corev1.ServicePort
	workerPorts = append(workerPorts, corev1.ServicePort{
		Name:       "rest-api-port-1",
		Protocol:   corev1.ProtocolTCP,
		Port:       7444,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 744},
		NodePort:   32501,
	})
	workerPorts = append(workerPorts, corev1.ServicePort{
		Name:       "thrift-ssl",
		Protocol:   corev1.ProtocolTCP,
		Port:       7712,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
	})
	return workerPorts
}

//container ports

func getGatewayContainerPorts() []corev1.ContainerPort {
	var gatewaydeployports []corev1.ContainerPort
	gatewaydeployports = append(gatewaydeployports, corev1.ContainerPort{
		ContainerPort: 8280,
		Protocol:      "TCP",
	})
	gatewaydeployports = append(gatewaydeployports, corev1.ContainerPort{
		ContainerPort: 8243,
		Protocol:      "TCP",
	})
	gatewaydeployports = append(gatewaydeployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      "TCP",
	})
	gatewaydeployports = append(gatewaydeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      "TCP",
	})
	return gatewaydeployports

}

func getKeyManagerContainerPorts() []corev1.ContainerPort {
	var keymanagerdeployports []corev1.ContainerPort
	keymanagerdeployports = append(keymanagerdeployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      "TCP",
	})
	keymanagerdeployports = append(keymanagerdeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      "TCP",
	})
	return keymanagerdeployports

}

func getPubDevTmContainerPorts() []corev1.ContainerPort {
	var pubdevtmdeployports []corev1.ContainerPort
	pubdevtmdeployports = append(pubdevtmdeployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      corev1.ProtocolTCP,
	})
	pubdevtmdeployports = append(pubdevtmdeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      corev1.ProtocolTCP,
	})
	pubdevtmdeployports = append(pubdevtmdeployports, corev1.ContainerPort{
		ContainerPort: 9711,
		Protocol:      corev1.ProtocolTCP,
	})
	pubdevtmdeployports = append(pubdevtmdeployports, corev1.ContainerPort{
		ContainerPort: 9611,
		Protocol:      corev1.ProtocolTCP,
	})
	pubdevtmdeployports = append(pubdevtmdeployports, corev1.ContainerPort{
		ContainerPort: 5672,
		Protocol:      corev1.ProtocolTCP,
	})
	return pubdevtmdeployports
}

func getDashContainerPorts() []corev1.ContainerPort {
	var dashdeployports []corev1.ContainerPort
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 9643,
		Protocol:      "TCP",
	})
	return dashdeployports
}

func getWorkerContainerPorts() []corev1.ContainerPort {

	var workerdeployports []corev1.ContainerPort
	workerdeployports = append(workerdeployports, corev1.ContainerPort{
		ContainerPort: 7612,
		Protocol:      "TCP",
	})
	workerdeployports = append(workerdeployports, corev1.ContainerPort{
		ContainerPort: 7712,
		Protocol:      "TCP",
	})
	workerdeployports = append(workerdeployports, corev1.ContainerPort{
		ContainerPort: 7444,
		Protocol:      "TCP",
	})
	return workerdeployports
}
