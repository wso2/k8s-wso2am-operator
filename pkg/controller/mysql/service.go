package mysql

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// // for handling mysql service
func MysqlService(apimanager *apimv1alpha1.APIManager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-svc",
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
