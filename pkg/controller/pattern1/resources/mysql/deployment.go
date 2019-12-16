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

package mysql
import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// mysqlDeployment creates a new Deployment for a MySQL
func MysqlDeployment(apimanager *apimv1alpha1.APIManager) *appsv1.Deployment {
	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	runasuser := int64(999)
	mysqlreplics := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2apim-with-analytics-mysql-deployment",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &mysqlreplics,
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
							ImagePullPolicy: "IfNotPresent",
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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "mysql-dbscripts",
									MountPath: "/docker-entrypoint-initdb.d",
								},
								{
									Name: "apim-rdbms-persistent-storage",
									MountPath: "/var/lib/mysql",
								},
							},
							Args: []string{
								"--max-connections",
								"10000",
							},

						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "mysql-dbscripts",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-dbscripts",
									},
								},
							},
						},
						{
							Name: "apim-rdbms-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName:"pvc-mysql",
								},
							},
						},
					},
					ServiceAccountName: "wso2am-pattern-1-svc-account",
				},
			},
		},
	}
}