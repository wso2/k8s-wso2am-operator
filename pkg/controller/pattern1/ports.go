package pattern1

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
)

//service ports

func getApimSvcCIPorts() []corev1.ServicePort {
	var apimsvs1ports []corev1.ServicePort
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "binary",
		Protocol:   corev1.ProtocolTCP,
		Port:       9611,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
	})
	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
		Name:       "binary-secure",
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
		Name:       "servlet-https",
		Protocol:   corev1.ProtocolTCP,
		Port:       9443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
	})
	return apimsvs1ports

}

func getApimCommonSvcNPPorts() []corev1.ServicePort {
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
func getApimCommonSvcLBPorts() []corev1.ServicePort {
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
func getApimCommonSvcCIPorts() []corev1.ServicePort {
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


func getDashNPPorts() []corev1.ServicePort {
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
func getDashLBPorts() []corev1.ServicePort {
	var dashports []corev1.ServicePort
	dashports = append(dashports, corev1.ServicePort{
		Name:       "pass-through-http",
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




//func getApim2DeployLBPorts() []corev1.ContainerPort {
//	var apim2deployports []corev1.ContainerPort
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 8280,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 8243,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 9763,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 9443,
//		Protocol:      "TCP",
//	})
//	return apim2deployports
//
//}

//not necessary
//func getApim2DeployNPPorts() []corev1.ContainerPort {
//	var apim2deployports []corev1.ContainerPort
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 30843, //8280,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 30806,//8243,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 32326,//9763,
//		Protocol:      "TCP",
//	})
//	apim2deployports = append(apim2deployports, corev1.ContainerPort{
//		ContainerPort: 32006, //9443,
//		Protocol:      "TCP",
//	})
//	return apim2deployports
//
//}

//func getApim1SvcNPPorts() []corev1.ServicePort {
//	var apimsvs1ports []corev1.ServicePort
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "pass-through-http",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       30838, //8280,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30838},
//		NodePort:   32004,
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "pass-through-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       30801, //8243,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30801},
//		NodePort:   32003,
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "servlet-http",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       32321, //9763,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32321},
//		NodePort:   32002,
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "servlet-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       32001, //9443
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32001},
//		NodePort:   32001,
//	})
//	return apimsvs1ports
//
//}
//
//func getApim2SvcNPPorts() []corev1.ServicePort {
//	var apimsvs2ports []corev1.ServicePort
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "pass-through-http",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       30843,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30843},
//		NodePort:   32009,
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "pass-through-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       30806,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 30806},
//		NodePort:   32008,
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "servlet-http",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       32326,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32326},
//		NodePort:   32007,
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "servlet-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       32006,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 32006},
//		NodePort:   32006,
//	})
//	return apimsvs2ports
//
//}
//
//func getApim1SvcLBPorts() []corev1.ServicePort {
//	var apimsvs1ports []corev1.ServicePort
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "binary",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9611,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "binary-secure",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9711,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "jms-tcp",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       5672,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
//	})
//	apimsvs1ports = append(apimsvs1ports, corev1.ServicePort{
//		Name:       "servlet-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9443,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
//	})
//	return apimsvs1ports
//
//}
//
//func getApim2SvcLBPorts() []corev1.ServicePort {
//	var apimsvs2ports []corev1.ServicePort
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "binary",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9611,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "binary-secure",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9711,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "jms-tcp",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       5672,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
//	})
//	apimsvs2ports = append(apimsvs2ports, corev1.ServicePort{
//		Name:       "servlet-https",
//		Protocol:   corev1.ProtocolTCP,
//		Port:       9443,
//		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
//	})
//	return apimsvs2ports
//
//}
