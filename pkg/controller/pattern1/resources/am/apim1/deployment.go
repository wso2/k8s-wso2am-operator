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

package apim1

import (
	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)
import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func Apim1Deployment(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap) *appsv1.Deployment {

	//defaultmode := int32(0407)

	//apim1cpu1,_:=resource.ParseQuantity(apimanager.Spec.ApimDeployment.Resources.Requests.CPU)
	//apim1mem1, _ := resource.ParseQuantity(apimanager.Spec.ApimDeployment.Resources.Requests.Memory)
	//apim1cpu2, _ :=resource.ParseQuantity(apimanager.Spec.ApimDeployment.Resources.Limits.CPU)
	//apim1mem2, _ := resource.ParseQuantity(apimanager.Spec.ApimDeployment.Resources.Limits.Memory)
	liveDelay := apimanager.Spec.Profiles.ApiManager1.Deployment.LivenessProbe.InitialDelaySeconds
	livePeriod := apimanager.Spec.Profiles.ApiManager1.Deployment.LivenessProbe.PeriodSeconds
	liveThres := apimanager.Spec.Profiles.ApiManager1.Deployment.LivenessProbe.FailureThreshold
	readyDelay := apimanager.Spec.Profiles.ApiManager1.Deployment.ReadinessProbe.InitialDelaySeconds
	readyPeriod := apimanager.Spec.Profiles.ApiManager1.Deployment.ReadinessProbe.PeriodSeconds
	readyThres := apimanager.Spec.Profiles.ApiManager1.Deployment.ReadinessProbe.FailureThreshold



	ControlConfigData := configMap.Data

	amMinReadySeconds:= "amMinReadySeconds"
	maxSurge:= "maxSurge"
	maxUnavailable:= "amMaxUnavailable"
	//amProbeInitialDelaySeconds := "amProbeInitialDelaySeconds"
	//periodSeconds:= "periodSeconds"
	imagePullPolicy:= "imagePullPolicy"
	//amRequestsCPU:= "amRequestsCPU"
	//amRequestsMemory:= "amRequestsMemory"
	//amLimitsCPU:= "amLimitsCPU"
	//amLimitsMemory:= "amLimitsMemory"


	minReadySec,_ := strconv.ParseInt(ControlConfigData[amMinReadySeconds], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData[maxSurge], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData[maxUnavailable], 10, 32)
	//liveDelay,_ := strconv.ParseInt(ControlConfigData[amProbeInitialDelaySeconds], 10, 32)
	//livePeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	//readyDelay,_ := strconv.ParseInt(ControlConfigData[amProbeInitialDelaySeconds], 10, 32)
	//readyPeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	imagePull,_ := ControlConfigData[imagePullPolicy]
	//reqCPU := ControlConfigData[amRequestsCPU]
	//reqMem := ControlConfigData[amRequestsMemory]
	//limitCPU := ControlConfigData[amLimitsCPU]
	//limitMem := ControlConfigData[amLimitsMemory]






	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apim-1-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Profiles.ApiManager1.Deployment.Replicas,
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
							Name:  "wso2-pattern-1-am",
							Image: apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Image,//"wso2/wso2am:3.0.0",
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
								InitialDelaySeconds: liveDelay, //int32(liveDelay)
								PeriodSeconds:      livePeriod, //int32(livePeriod),
								FailureThreshold: liveThres,
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

								InitialDelaySeconds: readyDelay,//int32(readyDelay),
								PeriodSeconds:  readyPeriod,     //int32(readyPeriod),
								FailureThreshold: readyThres,

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
							//		corev1.ResourceCPU:apim1cpu1,
							//		corev1.ResourceMemory:apim1mem1,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:apim1cpu2,
							//		corev1.ResourceMemory:apim1mem2,
							//	},
							//},

							ImagePullPolicy:corev1.PullPolicy(imagePull),

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
									Name: "wso2am-pattern-1-am-1-conf",
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
							Name: "wso2am-pattern-1-am-1-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										//Name: "wso2am-pattern-1-am-1-conf",
										Name:apimanager.Spec.Profiles.ApiManager1.DeploymentConfigmap,
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
					},
				},
			},
		},
	}
}