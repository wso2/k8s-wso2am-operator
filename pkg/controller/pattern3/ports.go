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

package pattern3

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

func getTrafficManagerSpecificSvcPorts() []corev1.ServicePort {
	var tmPorts []corev1.ServicePort
	tmPorts = append(tmPorts, corev1.ServicePort{
		Name:       "service",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return tmPorts
}

func getPubCommonSvcPorts() []corev1.ServicePort {
	var pubPorts []corev1.ServicePort
	pubPorts = append(pubPorts, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return pubPorts
}

func getDevportalCommonSvcPorts() []corev1.ServicePort {
	var devPorts []corev1.ServicePort
	devPorts = append(devPorts, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return devPorts
}

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
	workerPorts = append(workerPorts, corev1.ServicePort{
		Name:       "rest-api-port-1",
		Protocol:   corev1.ProtocolTCP,
		Port:       7444,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
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

func getPubContainerPorts() []corev1.ContainerPort {
	var pubdeployports []corev1.ContainerPort
	pubdeployports = append(pubdeployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      corev1.ProtocolTCP,
	})
	pubdeployports = append(pubdeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      corev1.ProtocolTCP,
	})
	return pubdeployports
}

func getDevportalContainerPorts() []corev1.ContainerPort {
	var devdeployports []corev1.ContainerPort
	devdeployports = append(devdeployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      corev1.ProtocolTCP,
	})
	devdeployports = append(devdeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      corev1.ProtocolTCP,
	})
	return devdeployports
}

func getTmContainerPorts() []corev1.ContainerPort {
	var tmdeployports []corev1.ContainerPort
	tmdeployports = append(tmdeployports, corev1.ContainerPort{
		ContainerPort: 9611,
		Protocol:      corev1.ProtocolTCP,
	})
	tmdeployports = append(tmdeployports, corev1.ContainerPort{
		ContainerPort: 9711,
		Protocol:      corev1.ProtocolTCP,
	})
	tmdeployports = append(tmdeployports, corev1.ContainerPort{
		ContainerPort: 5672,
		Protocol:      corev1.ProtocolTCP,
	})
	tmdeployports = append(tmdeployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      corev1.ProtocolTCP,
	})
	return tmdeployports
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
