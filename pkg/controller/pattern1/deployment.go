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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)


// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func Apim1Deployment(apimanager *apimv1alpha1.APIManager, x *configvalues) *appsv1.Deployment {

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}

	apim1VolumeMount, apim1Volume := getApim1Volumes(apimanager)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apim-1-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Profiles[0].Deployment.Replicas,
			MinReadySeconds:x.Minreadysec,
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
				Spec:corev1.PodSpec{
					HostAliases: []corev1.HostAlias{
						{
							IP: "127.0.0.1",
							Hostnames: []string{
								"wso2-am",
								"wso2-gateway",
							},
						},
					},

					Containers: []corev1.Container{
						{
							Name:  "wso2-pattern-1-am",
							Image: x.Amimage,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},
								InitialDelaySeconds: x.Livedelay,
								PeriodSeconds:     x.Liveperiod,
								FailureThreshold: x.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},

								InitialDelaySeconds: x.Readydelay,
								PeriodSeconds:  x.Readyperiod,
								FailureThreshold: x.Readythres,

							},

							Lifecycle: &corev1.Lifecycle{
								PreStop:&corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
										},
									},
								},
							},

							Resources:corev1.ResourceRequirements{
								//Requests:corev1.ResourceList{
								//	corev1.ResourceCPU:x.Reqcpu,
								//	corev1.ResourceMemory:x.Reqmem,
								//},
								//Limits:corev1.ResourceList{
								//	corev1.ResourceCPU:x.Limitcpu,
								//	corev1.ResourceMemory:x.Limitmem,
								//},
							},

							ImagePullPolicy:corev1.PullPolicy(x.Imagepull),

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8280,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8243,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9763,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9443,
									Protocol:      "TCP",
								},
							},
							Env: []corev1.EnvVar{
								// {
								// 	Name:  "HOST_NAME",
								// 	Value: "wso2-am",
								// },
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},

							VolumeMounts:apim1VolumeMount,

						},
					},

					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: apim1Volume,
				},
			},
		},
	}
}

// for handling apim-instance-2 deployment
func Apim2Deployment(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *appsv1.Deployment {
	//apim1cpu, _ :=resource.ParseQuantity("2000m")
	//apim1mem, _ := resource.ParseQuantity("2Gi")
	//apim1cpu2, _ :=resource.ParseQuantity("3000m")
	//apim1mem2, _ := resource.ParseQuantity("3Gi")
	//defaultmode := int32(0407)

	ControlConfigData := configMap.Data


	amMinReadySeconds:= "amMinReadySeconds"
	maxSurge:= "maxSurge"
	maxUnavailable:= "amMaxUnavailable"
	amProbeInitialDelaySeconds := "amProbeInitialDelaySeconds"
	periodSeconds:= "periodSeconds"
	imagePullPolicy:= "imagePullPolicy"
	//amRequestsCPU:= "amRequestsCPU"
	//amRequestsMemory:= "amRequestsMemory"
	//amLimitsCPU:= "amLimitsCPU"
	//amLimitsMemory:= "amLimitsMemory"


	minReadySec,_ := strconv.ParseInt(ControlConfigData[amMinReadySeconds], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData[maxSurge], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData[maxUnavailable], 10, 32)
	liveDelay,_ := strconv.ParseInt(ControlConfigData[amProbeInitialDelaySeconds], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData[amProbeInitialDelaySeconds], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	imagePull,_ := ControlConfigData[imagePullPolicy]
	//reqCPU := ControlConfigData[amRequestsCPU]
	//reqMem := ControlConfigData[amRequestsMemory]
	//limitCPU := ControlConfigData[amLimitsCPU]
	//limitMem := ControlConfigData[amLimitsMemory]


	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-2",
	}
	am2ConfigMap := "wso2am-pattern-1-am-2-conf"
	am2ConfigMapFromYaml := apimanager.Spec.Profiles[1].Deployment.Configmaps.DeploymentConfigmap
	if am2ConfigMapFromYaml != ""{
		am2ConfigMap = am2ConfigMapFromYaml
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apim-2-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:int32(minReadySec),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxSurges),
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxUnavail),
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
					HostAliases: []corev1.HostAlias{
						{
							IP: "127.0.0.1",
							Hostnames: []string{
								"wso2-am",
								"wso2-gateway",
							},
						},
					},

					Containers: []corev1.Container{
						{
							Name:  "wso2-pattern-1-am",
							Image: "wso2/wso2am:3.0.0",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},
								InitialDelaySeconds: int32(liveDelay),
								PeriodSeconds:       int32(livePeriod),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},

								InitialDelaySeconds: int32(readyDelay),
								PeriodSeconds:       int32(readyPeriod),
							},

							Lifecycle: &corev1.Lifecycle{
								PreStop:&corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
										},
									},
								},
							},

							//Resources:corev1.ResourceRequirements{
							//	Requests:corev1.ResourceList{
							//		corev1.ResourceCPU:apim1cpu,
							//		corev1.ResourceMemory:apim1mem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:apim1cpu2,
							//		corev1.ResourceMemory:apim1mem2,
							//	},
							//},

							ImagePullPolicy: corev1.PullPolicy(imagePull),

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8280,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8243,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9763,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9443,
									Protocol:      "TCP",
								},
							},
							Env: []corev1.EnvVar{
								// {
								// 	Name:  "HOST_NAME",
								// 	Value: "wso2-am",
								// },
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
									MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",
								},
								{
									Name: "wso2am-pattern-1-am-volume-claim-executionplans",
									MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
								},
								{
									Name: "wso2am-pattern-1-am-2-conf",
									MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
									SubPath:"deployment.toml",
								},
								//{
								//	Name: "wso2am-pattern-1-am-conf-entrypoint",
								//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
								//	SubPath:"docker-entrypoint.sh",
								//},
							},
						},
					},

					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName:"pvc-synapse-configs",
								},
							},
						},
						{
							Name: "wso2am-pattern-1-am-volume-claim-executionplans",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-execution-plans",
								},
							},
						},
						{
							Name: "wso2am-pattern-1-am-2-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: am2ConfigMap,//apimanager.Spec.Profiles.ApiManager2.DeploymentConfigmap,//"wso2am-pattern-1-am-2-conf",
									},
								},
							},
						},
						//{
						//	Name: "wso2am-pattern-1-am-conf-entrypoint",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "wso2am-pattern-1-am-conf-entrypoint",
						//			},
						//			DefaultMode:&defaultmode,
						//		},
						//	},
						//},
						//{
						//	Name: "mysql-jdbc-driver",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "mysql-jdbc-driver-cm",
						//			},
						//		},
						//	},
						//},
					},
				},
			},
		},
	}
}

// for handling analytics-dashboard deployment
func DashboardDeployment(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *appsv1.Deployment {

	ControlConfigData := configMap.Data

	analyticsMinReadySeconds:= "analyticsMinReadySeconds"
	maxSurge:= "maxSurge"
	maxUnavailable:= "maxUnavailable"
	analyticsProbeInitialDelaySeconds := "analyticsProbeInitialDelaySeconds"
	periodSeconds:= "periodSeconds"
	imagePullPolicy:= "imagePullPolicy"
	//analyticsCPU:= "analyticsCPU"
	//analyticsMemory:= "analyticsMemory"


	minReadySec,_ := strconv.ParseInt(ControlConfigData[analyticsMinReadySeconds], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData[maxSurge], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData[maxUnavailable], 10, 32)
	liveDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	imagePull,_ := ControlConfigData[imagePullPolicy]
	//reqCPU := ControlConfigData[analyticsCPU]
	//reqMem := ControlConfigData[analyticsMemory]
	//limitCPU := ControlConfigData[analyticsCPU]
	//limitMem := ControlConfigData[analyticsMemory]





	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	runasuser := int64(802)
	defaultMode := int32(0407)
	//dashcpu, _ :=resource.ParseQuantity("2000m")
	//dashmem, _ := resource.ParseQuantity("4Gi")


	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.DeploymentName,
			Name:      "analytics-dash-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:int32(minReadySec),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxSurges),
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxUnavail),
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

					//InitContainers: []corev1.Container{
					//	{
					//		Name: "init-apim-analytics-db",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
					//			"sh",
					//			"-c",
					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
					//		},
					//
					//	},
					//	{
					//		Name: "init-am-analytics-worker",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
					//		},
					//
					//	},
					//},
					Containers: []corev1.Container{
						{
							Name:  "wso2am-pattern-1-analytics-dashboard",
							Image: "wso2/wso2am-analytics-dashboard:3.0.0",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 32201",
										},
									},
								},
								InitialDelaySeconds: int32(liveDelay),
								PeriodSeconds:       int32(livePeriod),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 32201",
										},
									},
								},

								InitialDelaySeconds: int32(readyDelay),
								PeriodSeconds:       int32(readyPeriod),

							},

							Lifecycle: &corev1.Lifecycle{
								PreStop:&corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/dashboard.sh stop",
										},
									},
								},
							},

							//Resources:corev1.ResourceRequirements{
							//	Requests:corev1.ResourceList{
							//		corev1.ResourceCPU:dashcpu,
							//		corev1.ResourceMemory:dashmem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:dashcpu,
							//		corev1.ResourceMemory:dashmem,
							//	},
							//},

							ImagePullPolicy:corev1.PullPolicy(imagePull),

							SecurityContext: &corev1.SecurityContext{
								RunAsUser:&runasuser,
							},

							Ports: []corev1.ContainerPort{
								//{
								//	ContainerPort: 9713,
								//	Protocol:      "TCP",
								//},
								{
									ContainerPort: 9643,
									Protocol:      "TCP",
								},
								//{
								//	ContainerPort: 9613,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 7713,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 9091,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 7613,
								//	Protocol:      "TCP",
								//},


								/////////////////////////////////
								//{
								//	ContainerPort: 32269,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 32169,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 30269,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 30169,
								//	Protocol:      "TCP",
								//},
								//{
								//{
								//	ContainerPort: 9713,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 9643,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 9613,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 7713,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 9091,
								//	Protocol:      "TCP",
								//},
								//{
								//	ContainerPort: 7613,
								//	Protocol:      "TCP",
								//},

							},

							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-analytics-dashboard-conf",
									MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
									SubPath:"deployment.yaml",
								},
								//{
								//	Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
								//	MountPath: "/home/wso2carbon/wso2-config-volume/wso2/dashboard/bin/carbon.sh",
								//	SubPath:"carbon.sh",
								//},
								//{
								//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
								//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
								//	SubPath:"docker-entrypoint.sh",
								//},
							},
						},
					},

					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-analytics-dashboard-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "dash-conf",
									},
									DefaultMode:&defaultMode,
								},
							},
						},
						//{
						//	Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
						//			},
						//			DefaultMode:&defaultMode,
						//		},
						//	},
						//},
						//{
						//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
						//			},
						//			DefaultMode:&defaultMode,
						//		},
						//	},
						//},
					},
				},
			},
		},
	}
}

// for handling analytics-worker deployment
func WorkerDeployment(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *appsv1.Deployment {

	ControlConfigData := configMap.Data

	analyticsMinReadySeconds:= "analyticsMinReadySeconds"
	maxSurge:= "maxSurge"
	maxUnavailable:= "maxUnavailable"
	analyticsProbeInitialDelaySeconds := "analyticsProbeInitialDelaySeconds"
	periodSeconds:= "periodSeconds"
	imagePullPolicy:= "imagePullPolicy"
	//analyticsCPU:= "analyticsCPU"
	//analyticsMemory:= "analyticsMemory"


	minReadySec,_ := strconv.ParseInt(ControlConfigData[analyticsMinReadySeconds], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData[maxSurge], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData[maxUnavailable], 10, 32)
	liveDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	imagePull,_ := ControlConfigData[imagePullPolicy]
	//reqCPU := ControlConfigData[analyticsCPU]
	//reqMem := ControlConfigData[analyticsMemory]
	//limitCPU := ControlConfigData[analyticsCPU]
	//limitMem := ControlConfigData[analyticsMemory]




	//workercpu,_ := resource.ParseQuantity("2000m")
	//workermem,_ := resource.ParseQuantity("4Gi")
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	runasuser := int64(802)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.DeploymentName,
			Name:      "analytics-worker-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:int32(minReadySec),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxSurges),
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxUnavail),
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

					//InitContainers: []corev1.Container{
					//	{
					//		Name: "init-apim-analytics-db",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
					//			"sh",
					//			"-c",
					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
					//		},
					//
					//	},
					//	{
					//		Name: "init-am-analytics-worker",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
					//		},
					//
					//	},
					//},
					Containers: []corev1.Container{
						{
							Name:  "wso2am-pattern-1-analytics-worker",
							Image: "wso2/wso2am-analytics-worker:3.0.0",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 7712",
										},
									},
								},
								InitialDelaySeconds: int32(liveDelay),
								PeriodSeconds:       int32(livePeriod),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 7712",
										},
									},
								},

								InitialDelaySeconds: int32(readyDelay),
								PeriodSeconds:       int32(readyPeriod),

							},

							Lifecycle: &corev1.Lifecycle{
								PreStop:&corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
										},
									},
								},
							},

							//Resources:corev1.ResourceRequirements{
							//	Requests:corev1.ResourceList{
							//		corev1.ResourceCPU:workercpu,
							//		corev1.ResourceMemory:workermem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:workercpu,
							//		corev1.ResourceMemory:workermem,
							//	},
							//},

							ImagePullPolicy: corev1.PullPolicy(imagePull),

							SecurityContext: &corev1.SecurityContext{
								RunAsUser:&runasuser,
							},

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9764,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9444,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7612,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7712,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9091,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7071,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7444,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7575,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7576,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7577,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-analytics-worker-conf",
									MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
									//SubPath:"deployment.yaml",
								},
							},
						},
					},
					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-analytics-worker-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "worker-conf",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

//  for handling mysql deployment
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




