package mysql

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/wso2-incubator/wso2am-k8s-operator/pkg/controller/pattern1"

	//"strconv"
	//v1 "k8s.io/api/core/v1"

)

//  for handling mysql deployment
func MysqlDeployment(apimanager *apimv1alpha1.APIManager) *appsv1.Deployment {

	mysqlvolumemount, mysqlvolume := pattern1.GetMysqlVolumes(apimanager)


	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	runasuser := int64(999)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-"+apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "wso2apim-with-analytics-mysql",
							Image: "mysql:5.7",
							ImagePullPolicy: "Always",
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "root",
								},
								{
									Name:  "MYSQL_USER",
									Value: "wso2carbon",
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: "wso2carbon",
								},

							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: mysqlvolumemount,

							Args: []string{
								"--max-connections",
								"10000",
							},

						},
					},

					Volumes:mysqlvolume,

					// ServiceAccountName: "wso2am-pattern-1-svc-account",
				},
			},
		},
	}
}


