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
	"strconv"
	"strings"

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

	apim1SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(x.SecurityContext), ":")

	AssignSecurityContext(securityContextString, apim1SecurityContext)

	initContainers := getMysqlInitContainers(apimanager, &apim1Volume, &apim1VolumeMount)

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
							SecurityContext: apim1SecurityContext,
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

	initContainers := getMysqlInitContainers(apimanager, &apim2Volume, &apim2VolumeMount)

	// Checking for the availability of API Manager Server 1 deployment
	apim1InitContainer := corev1.Container{}
	apim1InitContainer.Name = "init-apim-1"
	apim1InitContainer.Image = "busybox:1.32"
	executionStr := "echo -e \"Checking for the availability of API Manager Server deployment\"; while ! nc -z \"wso2-am-1-svc\" 9711; do sleep 1; printf \"-\"; done; echo -e \"  >> APIM Server has started\";"
	apim1InitContainer.Command = []string{"/bin/sh", "-c", executionStr}
	initContainers = append(initContainers, apim1InitContainer)

	apim2SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, apim2SecurityContext)

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
							SecurityContext: apim2SecurityContext,
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

	dashVolumeMount, dashVolume := getAnalyticsDashVolumes(apimanager, num)

	initContainers := getMysqlInitContainers(apimanager, &dashVolume, &dashVolumeMount)

	dashbordSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, dashbordSecurityContext)

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
							SecurityContext: dashbordSecurityContext,
							ImagePullPolicy: corev1.PullPolicy(y.Imagepull),
							Ports:           dashdeployports,
							VolumeMounts:    dashVolumeMount,
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

	workerVolMounts, workerVols := getAnalyticsWorkerVolumes(apimanager, num)

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}

	workerContainerPorts := getWorkerContainerPorts()

	initContainers := getMysqlInitContainers(apimanager, &workerVols, &workerVolMounts)

	workerSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, workerSecurityContext)

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
							SecurityContext: workerSecurityContext,
							ImagePullPolicy: corev1.PullPolicy(y.Imagepull),
							Ports:           workerContainerPorts,
							VolumeMounts:    workerVolMounts,
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
