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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
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

	ControlConfigData := configMap.Data

	//amRequestsCPU:= "amRequestsCPU"
	//amRequestsMemory:= "amRequestsMemory"
	//amLimitsCPU:= "amLimitsCPU"
	//amLimitsMemory:= "amLimitsMemory"

	amImages:=ControlConfigData["amImage"]
	//imageFromYaml:= apimanager.Spec.Profiles[0].Image
	//if imageFromYaml != ""{
	//	amImages=imageFromYaml
	//}

	liveDelay,_ := strconv.ParseInt(ControlConfigData["amProbeInitialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles[0].LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
		//utilruntime.HandleError(fmt.Errorf("Live delay not present"))
	}

	livePeriod,_ := strconv.ParseInt(ControlConfigData["periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles[0].LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(liveDelayFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles[0].LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}


	readyDelay,_ := strconv.ParseInt(ControlConfigData["amProbeInitialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles[0].ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles[0].ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles[0].ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
	}

	minReadySec,_ := strconv.ParseInt(ControlConfigData["amMinReadySeconds"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles[0].MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}

	maxSurges,_ := strconv.ParseInt(ControlConfigData["maxSurge"], 10, 32)
	maxSurgeFromYaml := apimanager.Spec.Profiles[0].Strategy.RollingUpdate.MaxSurge
	if maxSurgeFromYaml !=0{
		maxSurges = int64(maxSurgeFromYaml)
	}

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["maxUnavailable"], 10, 32)
	maxUnavailFromYaml := apimanager.Spec.Profiles[0].Strategy.RollingUpdate.MaxUnavailable
	if maxUnavailFromYaml !=0{
		maxUnavail = int64(maxUnavailFromYaml)
	}

	imagePull,_ := ControlConfigData["imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles[0].ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}
	//reqCPU := 	resource.MustParse(ControlConfigData["amRequestsCPU"])
	//reqCPUyaml,_:= resource.ParseQuantity(apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Resources.Requests.CPU)
	//if reqCPUyaml == resource.MustParse("0"){
	//
	//}
	//////////////////
	//reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Resources.Requests.CPU)
	//if reqCPUFromYaml != resource.MustParse("0"){
	//	reqCPU = 	reqCPUFromYaml
	//}
	//reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Resources.Requests.Memory)
	//if reqMemFromYaml == resource.MustParse("0"){
	//	reqMemFromYaml = 	resource.MustParse(ControlConfigData["amRequestsMemory"])
	//}
	//limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Resources.Limits.CPU)
	//if limitCPUFromYaml == resource.MustParse("0"){
	//	limitCPUFromYaml = 	resource.MustParse(ControlConfigData["amLimitsCPU"])
	//}
	//limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Resources.Limits.Memory)
	//if limitMemFromYaml == resource.MustParse("0"){
	//	limitMemFromYaml = 	resource.MustParse(ControlConfigData["amLimitsMemory"])
	//}

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}
	am1ConfigMap := "wso2am-pattern-1-am-1-conf"
	am1ConfigMapFromYaml := apimanager.Spec.Profiles[0].DeploymentConfigmap
	if am1ConfigMapFromYaml != ""{
		am1ConfigMap = am1ConfigMapFromYaml
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
			Replicas: apimanager.Spec.Profiles[0].Replicas,
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
							Image: amImages,//apimanager.Spec.Profiles.ApiManager1.ApimDeployment.Image,//"wso2/wso2am:3.0.0",
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
								PeriodSeconds:      int32(livePeriod),
								FailureThreshold: int32(liveThres),
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
								PeriodSeconds:  int32(readyPeriod),
								FailureThreshold: int32(readyThres),

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
							//		corev1.ResourceCPU:reqCPU,
							//		corev1.ResourceMemory:reqMemFromYaml,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:limitCPUFromYaml,
							//		corev1.ResourceMemory:limitMemFromYaml,
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
									MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/am1-am2-deployment.toml",
									SubPath:"am1-am2-deployment.toml",
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

										Name:am1ConfigMap,
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