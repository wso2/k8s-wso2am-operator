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
package mysql

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/wso2/k8s-wso2am-operator/pkg/controller/pattern1"
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
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
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
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"mysqladmin",
											"ping",
										},
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:     10,
								FailureThreshold: 5,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											 "sh",
											 "-c",
											 "mysqladmin ping -u wso2carbon -p${MYSQL_PASSWORD}",
										},
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:5,
								FailureThreshold:1,
							},
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


