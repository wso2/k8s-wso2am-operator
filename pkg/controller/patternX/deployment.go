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

package patternX

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func ApimXDeployment(apimanager *apimv1alpha1.APIManager, r *apimv1alpha1.Profile, x *configvalues) *appsv1.Deployment {

	labels := map[string]string{
		"deployment": r.Name,
	}
	apimXVolumeMount, apimXVolume := getApimXVolumes(apimanager, *r, x)
	initContainers := getMysqlInitContainers(apimanager, &apimXVolume, &apimXVolumeMount)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        r.Deployment.Replicas,
			MinReadySeconds: x.Minreadysec,

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
							Name:  r.Name + "-container",
							Image: x.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},
								InitialDelaySeconds: x.Livedelay,
								PeriodSeconds:       x.Liveperiod,
								FailureThreshold:    x.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},

								InitialDelaySeconds: x.Readydelay,
								PeriodSeconds:       x.Readyperiod,
								FailureThreshold:    x.Readythres,
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
							VolumeMounts: apimXVolumeMount,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},
					Volumes: apimXVolume,
				},
			},
		},
	}
}

// for handling analytics-dashboard deployment
func DashboardXDeployment(apimanager *apimv1alpha1.APIManager, r *apimv1alpha1.Profile, x *configvalues) *appsv1.Deployment {

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9643",
	}

	labels := map[string]string{
		"deployment": r.Name,
	}
	runasuser := int64(802)

	dashVolumeMount, dashVolume := getDashboardXVolumes(apimanager, *r)
	initContainers := getMysqlInitContainers(apimanager, &dashVolume, &dashVolumeMount)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: x.Minreadysec,

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
							Name:  r.Name + "-container",
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
											"${WSO2_SERVER_HOME}/bin/dashboard.sh stop",
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

							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
							},

							Ports: []corev1.ContainerPort{

								{
									ContainerPort: 9713,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9643,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9613,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7713,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9091,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7613,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: dashVolumeMount,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},

					Volumes: dashVolume,
				},
			},
		},
	}
}

// for handling analytics-worker deployment
func WorkerXDeployment(apimanager *apimv1alpha1.APIManager, r *apimv1alpha1.Profile, x *configvalues) *appsv1.Deployment {
	workerVolMounts, workerVols := getWorkerXVolumes(apimanager, *r)
	initContainers := getMysqlInitContainers(apimanager, &workerVols, &workerVolMounts)

	labels := map[string]string{
		"deployment": r.Name,
	}
	runasuser := int64(802)
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depApiVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        apimanager.Spec.Replicas,
			MinReadySeconds: x.Minreadysec,

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
							Name:  r.Name + "-container",
							Image: x.Image,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 7712",
										},
									},
								},
								InitialDelaySeconds: x.Livedelay,
								PeriodSeconds:       x.Liveperiod,
								FailureThreshold:    x.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"nc -z localhost 7712",
										},
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
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
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

							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
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

							VolumeMounts: workerVolMounts,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},

					Volumes: workerVols,
				},
			},
		},
	}
}
