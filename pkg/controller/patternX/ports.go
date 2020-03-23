package patternX

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
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

