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

package pattern1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//service ports

func getApimCommonSvcNPPorts() []corev1.ServicePort {
	var apimsvs1ports []corev1.ServicePort
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       8280,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
		NodePort:   32004,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       8243,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
		NodePort:   32003,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
		NodePort:   32002,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
		NodePort:   32001,
	})
	return apimsvs1ports
}

func getApimCommonSvcPorts() []corev1.ServicePort {
	var apimsvs1ports []corev1.ServicePort
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       8280,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       8243,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return apimsvs1ports

}

func getApimSpecificSvcPorts() []corev1.ServicePort {
	var apimsvs1ports []corev1.ServicePort
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "thrift-binary",
		Protocol:   corev1.ProtocolTCP,
		Port:       9611,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "thrift-binary-secure",
		Protocol:   corev1.ProtocolTCP,
		Port:       9711,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "jms-tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       5672,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "https-servlet",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return apimsvs1ports

}

func getDashBoardNPPorts() []corev1.ServicePort {
	var dashports []corev1.ServicePort

	dashports = append(dashports, corev1.ServicePort{
		Name:       "analytics-dashboard",
		Protocol:   corev1.ProtocolTCP,
		Port:       32201,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32201},
		NodePort:   32201,
	})

	return dashports
}

func getDashBoardPorts() []corev1.ServicePort {

	var dashports []corev1.ServicePort
	dashports = append(dashports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9643,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
	})

	return dashports
}

//container ports

func getApimDeployLBPorts() []corev1.ContainerPort {
	var apim1deployports []corev1.ContainerPort
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 8280,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 8243,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      "TCP",
	})
	return apim1deployports

}
func getApimDeployNPPorts() []corev1.ContainerPort {
	var apim1deployports []corev1.ContainerPort
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 30838, //8280,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 30801, //8243,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 32321, //9763,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 32001, //9443,
		Protocol:      "TCP",
	})
	return apim1deployports

}

func getDashDeployLBPorts() []corev1.ContainerPort {
	var dashdeployports []corev1.ContainerPort
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 9713,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 9643,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 9613,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 7713,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 9091,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 7613,
		Protocol:      "TCP",
	})
	return dashdeployports
}
func getDashDeployNPPorts() []corev1.ContainerPort {
	var dashdeployports []corev1.ContainerPort
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 32271,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 32201,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 32171,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 30271,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 32649,
		Protocol:      "TCP",
	})
	dashdeployports = append(dashdeployports, corev1.ContainerPort{
		ContainerPort: 30171,
		Protocol:      "TCP",
	})
	return dashdeployports
}
