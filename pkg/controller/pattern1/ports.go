package pattern1

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
)

func getApim1SvcNPPorts() []corev1.ServicePort {
	var apimsvs1ports []corev1.ServicePort
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       30838, //8280,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30838},
		NodePort:   32004,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       30801, //8243,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30801},
		NodePort:   32003,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       32321, //9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32321},
		NodePort:   32002,
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       32001, //9443
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32001},
		NodePort:   32001,
	})
	return apimsvs1ports

}

func getApim2SvcNPPorts() []corev1.ServicePort {
	var apimsvs2ports []corev1.ServicePort
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       30843,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30843},
		NodePort:   32009,
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       30806,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30806},
		NodePort:   32008,
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       32326,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32326},
		NodePort:   32007,
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       32006,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32006},
		NodePort:   32006,
	})
	return apimsvs2ports

}

func getApim1SvcLBPorts() []corev1.ServicePort {
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

func getApim2SvcLBPorts() []corev1.ServicePort {
	var apimsvs2ports []corev1.ServicePort
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "pass-through-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       8280,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "pass-through-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       8243,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "servlet-http",
		Protocol:   corev1.ProtocolTCP,
		Port:       9763,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
	})
	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return apimsvs2ports

}

func getApim1DeployNPPorts() []corev1.ContainerPort {
	var apim1deployports []corev1.ContainerPort
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 30838, //8280,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 30801,//8243,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 32321,//9763,
		Protocol:      "TCP",
	})
	apim1deployports = append(apim1deployports, corev1.ContainerPort{
		ContainerPort: 32001, //9443,
		Protocol:      "TCP",
	})
	return apim1deployports

}

func getApim2DeployNPPorts() []corev1.ContainerPort {
	var apim2deployports []corev1.ContainerPort
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 30843, //8280,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 30806,//8243,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 32326,//9763,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 32006, //9443,
		Protocol:      "TCP",
	})
	return apim2deployports

}

func getApim1DeployLBPorts() []corev1.ContainerPort {
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

func getApim2DeployLBPorts() []corev1.ContainerPort {
	var apim2deployports []corev1.ContainerPort
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 8280,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 8243,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 9763,
		Protocol:      "TCP",
	})
	apim2deployports = append(apim2deployports, corev1.ContainerPort{
		ContainerPort: 9443,
		Protocol:      "TCP",
	})
	return apim2deployports

}