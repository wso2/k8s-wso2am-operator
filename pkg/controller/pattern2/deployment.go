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

package pattern2

import (
	"strconv"
	"strings"

	"k8s.io/klog"
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// PubDev1Deployment creates a new Deployment for a PubDevTm instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it...
func PubDev1Deployment(apimanager *apimv1alpha1.APIManager, x *configvalues, num int) *appsv1.Deployment {
	klog.Info("PubDevTm-1 Depl Starting.........")
	useMysql := true
	allowAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.AllowAnalytics != "" {
		allowAnalytics, _ = strconv.ParseBool(apimanager.Spec.AllowAnalytics)
	}

	klog.Info("PubDevTm-1 AllowAnalytics: ", allowAnalytics)

	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-1",
	}

	pubDevTm1VolumeMount, pubDevTm1Volume := getDevPubTm1Volumes(apimanager, num)
	pubDevTm1deployports := getPubDevTmContainerPorts()

	klog.Info(pubDevTm1VolumeMount)

	// cmdstring := []string{
	// 	"/bin/sh",
	// 	"-c",
	// 	"nc -z 9443",
	// }

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &pubDevTm1Volume, &pubDevTm1VolumeMount)
	}

	klog.Info("PubDevTm-1 Containers Done")

	if allowAnalytics {
		getInitContainers([]string{"init-am-analytics-worker"}, &initContainers)
		klog.Info("Pub-Dev-Tm-1 Containers", initContainers[1])
	}

	klog.Info("PubDevTm-1 Container Phase 1 Done")

	//initContainers = append(initContainers, getInitContainers([]string{"init-am-analytics-worker"}))
	pubDev1SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(x.SecurityContext), ":")

	AssignSecurityContext(securityContextString, pubDev1SecurityContext)
	klog.Info("PubDevTm-1 Security Context Done")
	klog.Info("Ready Delay", x.Readydelay)
	klog.Info("Live Delay", x.Livedelay)
	return &appsv1.Deployment{
		// TypeMeta: metav1.TypeMeta{
		// 	APIVersion: depApiVersion,
		// 	Kind:       deploymentKind,
		// },
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
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
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
				// RollingUpdate: &appsv1.RollingUpdateDeployment{
				// 	MaxSurge: &intstr.IntOrString{
				// 		Type:   intstr.Int,
				// 		IntVal: x.Maxsurge,
				// 	},
				// 	MaxUnavailable: &intstr.IntOrString{
				// 		Type:   intstr.Int,
				// 		IntVal: x.Maxunavail,
				// 	},
				// },
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
							Name:  "wso2-pattern-2-am-1",
							Image: x.Image,
							// LivenessProbe: &corev1.Probe{
							// 	InitialDelaySeconds: x.Livedelay,
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},
							// 	PeriodSeconds:       x.Liveperiod,
							// 	FailureThreshold:    x.Livethres,
							// },
							// ReadinessProbe: &corev1.Probe{
							// 	InitialDelaySeconds: x.Readydelay,
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},

							// 	PeriodSeconds:       x.Readyperiod,
							// 	FailureThreshold:    x.Readythres,
							// },

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
							SecurityContext: pubDev1SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(x.Imagepull),
							Ports:           pubDevTm1deployports,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "JVM_MEM_OPTS",
									Value:  x.JvmMemOpts,
								},
								{
									Name:  "ALLOW_ANALYTICS",
									Value: apimanager.Spec.AllowAnalytics,
								},
							},
							VolumeMounts: pubDevTm1VolumeMount,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},

					Volumes: pubDevTm1Volume,
				},
			},
		},
	}
}

func PubDev2Deployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {

	useMysql := true
	allowAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.AllowAnalytics != "" {
		allowAnalytics, _ = strconv.ParseBool(apimanager.Spec.AllowAnalytics)
	}

	pubDevTm2VolumeMount, pubDevTm2Volume := getDevPubTm2Volumes(apimanager, num)
	pubDevTm2deployports := getPubDevTmContainerPorts()

	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-2",
	}

	// cmdstring := []string{
	// 	"/bin/sh",
	// 	"-c",
	// 	"nc -z localhost 9443",
	// }

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &pubDevTm2Volume, &pubDevTm2VolumeMount)
	}

	if allowAnalytics {
		getInitContainers([]string{"init-am-analytics-worker"}, &initContainers)
		klog.Info("Pub-Dev-Tm-2 Containers", initContainers[1])
	}

	//initContainers = append(initContainers, getAnalyticsWorkerInitContainers([]string{"init-am-analytics-worker"}))
	pubDev2SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, pubDev2SecurityContext)

	// Checking for the availability of API Manager Server 1 deployment

	// apim1InitContainer := corev1.Container{}
	// apim1InitContainer.Name = "init-apim-1"
	// apim1InitContainer.Image = "busybox:1.32"
	// executionStr := "echo -e \"Checking for the availability of API Manager Server deployment\"; while ! nc -z \"wso2-am-1-svc\" 9711; do sleep 1; printf \"-\"; done; echo -e \"  >> APIM Server has started\";"
	// apim1InitContainer.Command = []string{"/bin/sh", "-c", executionStr}
	//initContainers = append(initContainers, apim1InitContainer)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
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
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
				// RollingUpdate: &appsv1.RollingUpdateDeployment{
				// 	MaxSurge: &intstr.IntOrString{
				// 		Type:   intstr.Int,
				// 		IntVal: z.Maxsurge,
				// 	},
				// 	MaxUnavailable: &intstr.IntOrString{
				// 		Type:   intstr.Int,
				// 		IntVal: z.Maxunavail,
				// 	},
				// },
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
							Name:  "wso2-pattern-2-am-2",
							Image: z.Image,
							// LivenessProbe: &corev1.Probe{
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},
							// 	InitialDelaySeconds: z.Livedelay,
							// 	PeriodSeconds:       z.Liveperiod,
							// 	FailureThreshold:    z.Livethres,
							// },
							// ReadinessProbe: &corev1.Probe{
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},
							// 	InitialDelaySeconds: z.Readydelay,
							// 	PeriodSeconds:       z.Readyperiod,
							// 	FailureThreshold:    z.Readythres,
							// },

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
							SecurityContext: pubDev2SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           pubDevTm2deployports,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "JVM_MEM_OPTS",
									Value: z.JvmMemOpts,
								},
								{
									Name:  "ALLOW_ANALYTICS",
									Value: apimanager.Spec.AllowAnalytics,
								},
							},
							VolumeMounts: pubDevTm2VolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: pubDevTm2Volume,
				},
			},
		},
	}
}

func GatewayDeployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {
	useMysql := true
	allowAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.AllowAnalytics != "" {
		allowAnalytics, _ = strconv.ParseBool(apimanager.Spec.AllowAnalytics)
	}

	gatewayVolumeMount, gatewayVolume := getgatewayVolumes(apimanager, num)
	gatewaydeployports := getGatewayContainerPorts()

	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-2",
	}

	// cmdstring := []string{
	// 	"/bin/sh",
	// 	"-c",
	// 	"nc -z localhost 9443",
	// }

	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &gatewayVolume, &gatewayVolumeMount)
	}

	if allowAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km", "init-apim-1", "init-apim-2"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km", "init-apim-1", "init-apim-2"}, &initContainers)
	}

	gatewaySecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, gatewaySecurityContext)

	// Checking for the availability of API Manager Server 1 deployment

	// apim1InitContainer := corev1.Container{}
	// apim1InitContainer.Name = "init-apim-1"
	// apim1InitContainer.Image = "busybox:1.32"
	// executionStr := "echo -e \"Checking for the availability of API Manager Server deployment\"; while ! nc -z \"wso2-am-1-svc\" 9711; do sleep 1; printf \"-\"; done; echo -e \"  >> APIM Server has started\";"
	// apim1InitContainer.Command = []string{"/bin/sh", "-c", executionStr}
	// initContainers = append(initContainers, apim1InitContainer)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-gw-" + apimanager.Name,
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
							Name:  "wso2-pattern-2-gw",
							Image: z.Image,
							// LivenessProbe: &corev1.Probe{
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},
							// 	InitialDelaySeconds: z.Livedelay,
							// 	PeriodSeconds:       z.Liveperiod,
							// 	FailureThreshold:    z.Livethres,
							// },
							// ReadinessProbe: &corev1.Probe{
							// 	Handler: corev1.Handler{
							// 		Exec: &corev1.ExecAction{
							// 			Command: cmdstring,
							// 		},
							// 	},
							// 	InitialDelaySeconds: z.Readydelay,
							// 	PeriodSeconds:       z.Readyperiod,
							// 	FailureThreshold:    z.Readythres,
							// },

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
							SecurityContext: gatewaySecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           gatewaydeployports,
							Env: []corev1.EnvVar{
								{
									Name:  "PROFILE_NAME",
									Value: "gateway-worker",
								},
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "JVM_MEM_OPTS",
									Value: z.JvmMemOpts,
								},
								{
									Name:  "ALLOW_ANALYTICS",
									Value: strconv.FormatBool(allowAnalytics),
								},
							},
							VolumeMounts: gatewayVolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: gatewayVolume,
				},
			},
		},
	}
}

func KeyManagerDeployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.StatefulSet {
	useMysql := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}

	kmVolumeMount, kmVolume := getKeyManagerVolumes(apimanager, num)
	kmdeployports := getKeyManagerContainerPorts()

	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
		"node":       "wso2am-pattern-2-am-2",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &kmVolume, &kmVolumeMount)
	}

	keyManagerSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, keyManagerSecurityContext)

	// Checking for the availability of API Manager Server 1 deployment

	// apim1InitContainer := corev1.Container{}
	// apim1InitContainer.Name = "init-apim-1"
	// apim1InitContainer.Image = "busybox:1.32"
	// executionStr := "echo -e \"Checking for the availability of API Manager Server deployment\"; while ! nc -z \"wso2-am-1-svc\" 9711; do sleep 1; printf \"-\"; done; echo -e \"  >> APIM Server has started\";"
	// apim1InitContainer.Command = []string{"/bin/sh", "-c", executionStr}
	// initContainers = append(initContainers, apim1InitContainer)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       statefulsetKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-km-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.StatefulSetSpec{
			Replicas:    apimanager.Spec.Replicas,
			ServiceName: "wso2-am-km-svc",
			// MinReadySeconds: z.Minreadysec,
			// Strategy: appsv1.DeploymentStrategy{
			// 	Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
			// 	RollingUpdate: &appsv1.RollingUpdateDeployment{
			// 		MaxSurge: &intstr.IntOrString{
			// 			Type:   intstr.Int,
			// 			IntVal: z.Maxsurge,
			// 		},
			// 		MaxUnavailable: &intstr.IntOrString{
			// 			Type:   intstr.Int,
			// 			IntVal: z.Maxunavail,
			// 		},
			// 	},
			// },

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
							Name:  "wso2-pattern-2-km",
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
							SecurityContext: keyManagerSecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           kmdeployports,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "PROFILE_NAME",
									Value: "api-key-manager",
								},
							},
							VolumeMounts: kmVolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: kmVolume,
				},
			},
		},
	}
}

// for handling analytics-dashboard deployment
func DashboardDeployment(apimanager *apimv1alpha1.APIManager, y *configvalues, num int) *appsv1.Deployment {
	useMysql := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}

	dashdeployports := getDashContainerPorts()

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9643",
	}

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	//runasuser := int64(802)
	//defaultMode := int32(0407)

	dashVolumeMount, dashVolume := getAnalyticsDashVolumes(apimanager, num)

	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &dashVolume, &dashVolumeMount)
	}

	getInitContainers([]string{"init-am"}, &initContainers)

	dashboardSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, dashboardSecurityContext)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
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
							Name:  "wso2am-pattern-2-analytics-dashboard",
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
							SecurityContext: dashboardSecurityContext,
							Ports:           dashdeployports,
							// Env: []corev1.EnvVar{
							// 	{
							// 		Name:  "PROFILE_NAME",
							// 		Value: "gateway-worker",
							// 	},
							// },
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
func WorkerDeployment(apimanager *apimv1alpha1.APIManager, y *configvalues, num int) *appsv1.StatefulSet {
	useMysql := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	workerVolMounts, workerVols := getAnalyticsWorkerVolumes(apimanager, num)

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	//runasuser := int64(802)

	workerContainerPorts := getWorkerContainerPorts()
	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &workerVols, &workerVolMounts)
	}

	workerSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, workerSecurityContext)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       statefulsetKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-worker-statefulset",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    apimanager.Spec.Replicas,
			ServiceName: "wso2-am-analytics-worker-headless-svc",
			// MinReadySeconds: y.Minreadysec,
			// Strategy: appsv1.DeploymentStrategy{
			// 	Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
			// 	RollingUpdate: &appsv1.RollingUpdateDeployment{
			// 		MaxSurge: &intstr.IntOrString{
			// 			Type:   intstr.Int,
			// 			IntVal: y.Maxsurge,
			// 		},
			// 		MaxUnavailable: &intstr.IntOrString{
			// 			Type:   intstr.Int,
			// 			IntVal: y.Maxunavail,
			// 		},
			// 	},
			// },

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

							SecurityContext: workerSecurityContext,
							Ports:           workerContainerPorts,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name: "NODE_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							VolumeMounts: workerVolMounts,
						},
					},
					ServiceAccountName: y.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: y.ImagePullSecret,
						},
					},
					Volumes: workerVols,
				},
			},
		},
	}
}

// AssignSecurityContext is ..
func AssignSecurityContext(securityContextString []string, securityContext *corev1.SecurityContext) {

	if securityContextString[0] == "runAsUser" {
		securityContextVal, _ := strconv.ParseInt(securityContextString[1], 10, 64)
		securityContext.RunAsUser = &securityContextVal
	} else if securityContextString[0] == "runAsGroup" {
		securityContextVal, _ := strconv.ParseInt(securityContextString[1], 10, 64)
		securityContext.RunAsGroup = &securityContextVal
	} else if securityContextString[0] == "runAsNonRoot" {
		securityContextVal, _ := strconv.ParseBool(securityContextString[1])
		securityContext.RunAsNonRoot = &securityContextVal
	} else if securityContextString[0] == "privileged" {
		securityContextVal, _ := strconv.ParseBool(securityContextString[1])
		securityContext.Privileged = &securityContextVal
	} else if securityContextString[0] == "readOnlyRootFilesystem" {
		securityContextVal, _ := strconv.ParseBool(securityContextString[1])
		securityContext.ReadOnlyRootFilesystem = &securityContextVal
	} else if securityContextString[0] == "allowPrivilegeEscalation" {
		securityContextVal, _ := strconv.ParseBool(securityContextString[1])
		securityContext.AllowPrivilegeEscalation = &securityContextVal
	}
}
