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

package pattern1

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func Apim1Deployment(apimanager *apimv1alpha1.APIManager, x *configvalues, num int) *appsv1.Deployment {

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node":       "wso2am-pattern-1-am-1",
	}

	apim1VolumeMount, apim1Volume := getApim1Volumes(apimanager, num)
	apim1deployports := getApimContainerPorts()

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if x.UseMysqlPod {
		// Downloading mysql connector
		// init container
		mysqlConnectorContainer := corev1.Container{}
		mysqlConnectorContainer.Name = "init-mysql-connector-download"
		mysqlConnectorContainer.Image = "busybox:1.32"
		downloadCmdStr := `set -e
              connector_version=8.0.17
              wget https://repo1.maven.org/maven2/mysql/mysql-connector-java/${connector_version}/mysql-connector-java-${connector_version}.jar -P /mysql-connector-jar/`
		mysqlConnectorContainer.Command = []string{"/bin/sh", "-c", downloadCmdStr}
		mysqlConnectorContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "mysql-connector-jar",
				MountPath: "/mysql-connector-jar",
			},
		}
		initContainers = append(initContainers, mysqlConnectorContainer)
		// volume for downloaded mysql connector
		apim1Volume = append(apim1Volume, corev1.Volume{
			Name: "mysql-connector-jar",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// volume mount for downloaded mysql connector
		apim1VolumeMount = append(apim1VolumeMount, corev1.VolumeMount{
			Name:      "mysql-connector-jar",
			MountPath: "/home/wso2carbon/wso2-artifact-volume/lib",
		})

		// Checking for the availability of MySQL Server deployment
		// init container
		mysqlWaitContainer := corev1.Container{}
		mysqlWaitContainer.Name = "init-mysql"
		mysqlWaitContainer.Image = "busybox:1.32"
		executionStr := "echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"mysql-svc\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";"
		mysqlWaitContainer.Command = []string{"/bin/sh", "-c", executionStr}
		initContainers = append(initContainers, mysqlWaitContainer)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-1-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: x.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: x.Maxsurge,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: x.Maxunavail,
					},
				},
			},

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers: []corev1.Container{
						{
							Name:  "wso2-pattern-1-am",
							Image: x.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},
								InitialDelaySeconds: x.Livedelay,
								PeriodSeconds:       x.Liveperiod,
								FailureThreshold:    x.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},

								InitialDelaySeconds: x.Readydelay,
								PeriodSeconds:       x.Readyperiod,
								FailureThreshold:    x.Readythres,
							},

							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    x.Reqcpu,
									corev1.ResourceMemory: x.Reqmem,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    x.Limitcpu,
									corev1.ResourceMemory: x.Limitmem,
								},
							},
							ImagePullPolicy: corev1.PullPolicy(x.Imagepull),
							Ports:           apim1deployports,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: apim1VolumeMount,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},

					Volumes: apim1Volume,
				},
			},
		},
	}
}

func Apim2Deployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {

	apim2VolumeMount, apim2Volume := getApim2Volumes(apimanager, num)
	apim2deployports := getApimContainerPorts()

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node":       "wso2am-pattern-1-am-2",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if z.UseMysqlPod {
		// Downloading mysql connector
		// init container
		mysqlConnectorContainer := corev1.Container{}
		mysqlConnectorContainer.Name = "init-mysql-connector-download"
		mysqlConnectorContainer.Image = "busybox:1.32"
		downloadCmdStr := `set -e
              connector_version=8.0.17
              wget https://repo1.maven.org/maven2/mysql/mysql-connector-java/${connector_version}/mysql-connector-java-${connector_version}.jar -P /mysql-connector-jar/`
		mysqlConnectorContainer.Command = []string{"/bin/sh", "-c", downloadCmdStr}
		mysqlConnectorContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "mysql-connector-jar",
				MountPath: "/mysql-connector-jar",
			},
		}
		initContainers = append(initContainers, mysqlConnectorContainer)
		// volume for downloaded mysql connector
		apim2Volume = append(apim2Volume, corev1.Volume{
			Name: "mysql-connector-jar",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// volume mount for downloaded mysql connector
		apim2VolumeMount = append(apim2VolumeMount, corev1.VolumeMount{
			Name:      "mysql-connector-jar",
			MountPath: "/home/wso2carbon/wso2-artifact-volume/lib",
		})

		// Checking for the availability of MySQL Server deployment
		// init container
		mysqlWaitContainer := corev1.Container{}
		mysqlWaitContainer.Name = "init-mysql"
		mysqlWaitContainer.Image = "busybox:1.32"
		executionStr := "echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"mysql-svc\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";"
		mysqlWaitContainer.Command = []string{"/bin/sh", "-c", executionStr}
		initContainers = append(initContainers, mysqlWaitContainer)
	}

	// Checking for the availability of API Manager Server 1 deployment
	apim1InitContainer := corev1.Container{}
	apim1InitContainer.Name = "init-apim-1"
	apim1InitContainer.Image = "busybox:1.32"
	executionStr := "echo -e \"Checking for the availability of API Manager Server deployment\"; while ! nc -z \"wso2-am-1-svc\" 9711; do sleep 1; printf \"-\"; done; echo -e \"  >> APIM Server has started\";"
	apim1InitContainer.Command = []string{"/bin/sh", "-c", executionStr}
	initContainers = append(initContainers, apim1InitContainer)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-2-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: z.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: z.Maxsurge,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: z.Maxunavail,
					},
				},
			},

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers: []corev1.Container{
						{
							Name:  "wso2-pattern-1-am",
							Image: z.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},
								InitialDelaySeconds: z.Livedelay,
								PeriodSeconds:       z.Liveperiod,
								FailureThreshold:    z.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},
								InitialDelaySeconds: z.Readydelay,
								PeriodSeconds:       z.Readyperiod,
								FailureThreshold:    z.Readythres,
							},

							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
										},
									},
								},
							},

							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    z.Reqcpu,
									corev1.ResourceMemory: z.Reqmem,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    z.Limitcpu,
									corev1.ResourceMemory: z.Limitmem,
								},
							},
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           apim2deployports,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: apim2VolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: apim2Volume,
				},
			},
		},
	}
}

// for handling analytics-dashboard deployment
func DashboardDeployment(apimanager *apimv1alpha1.APIManager, y *configvalues, num int) *appsv1.Deployment {

	dashdeployports := getDashContainerPorts()

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9643",
	}

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	runasuser := int64(802)
	//defaultMode := int32(0407)

	dashVolumeMount, dashVolume := getAnalyticsDashVolumes(apimanager, num)

	initContainers := []corev1.Container{}

	if y.UseMysqlPod {
		// Downloading mysql connector
		// init container
		mysqlConnectorContainer := corev1.Container{}
		mysqlConnectorContainer.Name = "init-mysql-connector-download"
		mysqlConnectorContainer.Image = "busybox:1.32"
		downloadCmdStr := `set -e
              connector_version=8.0.17
              wget https://repo1.maven.org/maven2/mysql/mysql-connector-java/${connector_version}/mysql-connector-java-${connector_version}.jar -P /mysql-connector-jar/`
		mysqlConnectorContainer.Command = []string{"/bin/sh", "-c", downloadCmdStr}
		mysqlConnectorContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "mysql-connector-jar",
				MountPath: "/mysql-connector-jar",
			},
		}
		initContainers = append(initContainers, mysqlConnectorContainer)
		// volume for downloaded mysql connector
		dashVolume = append(dashVolume, corev1.Volume{
			Name: "mysql-connector-jar",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// volume mount for downloaded mysql connector
		dashVolumeMount = append(dashVolumeMount, corev1.VolumeMount{
			Name:      "mysql-connector-jar",
			MountPath: "/home/wso2carbon/wso2-artifact-volume/lib",
		})

		// Checking for the availability of MySQL Server deployment
		// init container
		mysqlWaitContainer := corev1.Container{}
		mysqlWaitContainer.Name = "init-mysql"
		mysqlWaitContainer.Image = "busybox:1.32"
		executionStr := "echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"mysql-svc\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";"
		mysqlWaitContainer.Command = []string{"/bin/sh", "-c", executionStr}
		initContainers = append(initContainers, mysqlWaitContainer)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-dashboard-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: y.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: y.Maxsurge,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: y.Maxunavail,
					},
				},
			},

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers: []corev1.Container{
						{
							Name:  "wso2am-pattern-1-analytics-dashboard",
							Image: y.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},
								InitialDelaySeconds: y.Livedelay,
								PeriodSeconds:       y.Liveperiod,
								FailureThreshold:    y.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},

								InitialDelaySeconds: y.Readydelay,
								PeriodSeconds:       y.Readyperiod,
								FailureThreshold:    y.Readythres,
							},

							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/dashboard.sh stop",
										},
									},
								},
							},

							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    y.Reqcpu,
									corev1.ResourceMemory: y.Reqmem,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    y.Limitcpu,
									corev1.ResourceMemory: y.Limitmem,
								},
							},
							ImagePullPolicy: corev1.PullPolicy(y.Imagepull),
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
							},
							Ports:        dashdeployports,
							VolumeMounts: dashVolumeMount,
						},
					},
					ServiceAccountName: y.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: y.ImagePullSecret,
						},
					},
					Volumes: dashVolume,
				},
			},
		},
	}
}

// for handling analytics-worker deployment
func WorkerDeployment(apimanager *apimv1alpha1.APIManager, y *configvalues, num int) *appsv1.Deployment {

	workervolumemounts, workervolume := getAnalyticsWorkerVolumes(apimanager, num)

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	runasuser := int64(802)

	workerContainerPorts := getWorkerContainerPorts()

	initContainers := []corev1.Container{}

	if y.UseMysqlPod {
		// Downloading mysql connector
		// init container
		mysqlConnectorContainer := corev1.Container{}
		mysqlConnectorContainer.Name = "init-mysql-connector-download"
		mysqlConnectorContainer.Image = "busybox:1.32"
		downloadCmdStr := `set -e
              connector_version=8.0.17
              wget https://repo1.maven.org/maven2/mysql/mysql-connector-java/${connector_version}/mysql-connector-java-${connector_version}.jar -P /mysql-connector-jar/`
		mysqlConnectorContainer.Command = []string{"/bin/sh", "-c", downloadCmdStr}
		mysqlConnectorContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "mysql-connector-jar",
				MountPath: "/mysql-connector-jar",
			},
		}
		initContainers = append(initContainers, mysqlConnectorContainer)
		// volume for downloaded mysql connector
		workervolume = append(workervolume, corev1.Volume{
			Name: "mysql-connector-jar",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// volume mount for downloaded mysql connector
		workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
			Name:      "mysql-connector-jar",
			MountPath: "/home/wso2carbon/wso2-artifact-volume/lib",
		})

		// Checking for the availability of MySQL Server deployment
		// init container
		mysqlWaitContainer := corev1.Container{}
		mysqlWaitContainer.Name = "init-mysql"
		mysqlWaitContainer.Image = "busybox:1.32"
		executionStr := "echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"mysql-svc\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";"
		mysqlWaitContainer.Command = []string{"/bin/sh", "-c", executionStr}
		initContainers = append(initContainers, mysqlWaitContainer)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-worker-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: y.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: y.Maxsurge,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: y.Maxunavail,
					},
				},
			},

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers: []corev1.Container{
						{
							Name:  "wso2am-pattern-1-analytics-worker",
							Image: y.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9444",
										},
									},
								},
								InitialDelaySeconds: y.Livedelay,
								PeriodSeconds:       y.Liveperiod,
								FailureThreshold:    y.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9444",
										},
									},
								},

								InitialDelaySeconds: y.Readydelay,
								PeriodSeconds:       y.Readyperiod,
								FailureThreshold:    y.Readythres,
							},

							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
										},
									},
								},
							},

							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    y.Reqcpu,
									corev1.ResourceMemory: y.Reqmem,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    y.Limitcpu,
									corev1.ResourceMemory: y.Limitmem,
								},
							},

							ImagePullPolicy: corev1.PullPolicy(y.Imagepull),

							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
							},
							Ports:        workerContainerPorts,
							VolumeMounts: workervolumemounts,
						},
					},
					ServiceAccountName: y.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: y.ImagePullSecret,
						},
					},
					Volumes: workervolume,
				},
			},
		},
	}
}
