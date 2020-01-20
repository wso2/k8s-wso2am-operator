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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"

	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"strconv"
	//"k8s.io/apimachinery/pkg/api/resource"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func ApimXDeployment(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile ) *appsv1.Deployment {

	labels := map[string]string{
		"deployment": r.Name,

	}
	apimXVolumeMount, apimXVolume := getApimXVolumes(apimanager,*r)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: r.Deployment.Replicas,
			MinReadySeconds:r.Deployment.MinReadySeconds,

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
							Name:  r.Name+"container",
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
								InitialDelaySeconds: r.Deployment.LivenessProbe.InitialDelaySeconds,
								PeriodSeconds:      r.Deployment.LivenessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.LivenessProbe.FailureThreshold,
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

								InitialDelaySeconds: r.Deployment.ReadinessProbe.InitialDelaySeconds,
								PeriodSeconds:  r.Deployment.ReadinessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.ReadinessProbe.FailureThreshold,

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

							ImagePullPolicy:corev1.PullPolicy(r.Deployment.ImagePullPolicy),

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
							VolumeMounts: apimXVolumeMount,
						},
					},

					Volumes: apimXVolume,

				},
			},
		},
	}
}

// for handling analytics-dashboard deployment
func DashboardXDeployment(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile) *appsv1.Deployment {

	cmdstring := []string{}
	if apimanager.Spec.Service.Type=="NodePort"{
		cmdstring = []string{
			"/bin/sh",
			"-c",
			"nc -z localhost 32201",
		}
	} else {
		cmdstring = []string{
			"/bin/sh",
			"-c",
			"nc -z localhost 32201",
		}
	}

	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-dashboard",
	}
	runasuser := int64(802)
	//defaultMode := int32(0407)

	dashVolumeMount, dashVolume := getDashboardXVolumes(apimanager,*r)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:r.Deployment.MinReadySeconds,

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
							Name:  "wso2am-pattern-1-analytics-dashboard",
							Image: "wso2/wso2am-analytics-dashboard:3.0.0",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:cmdstring,
									},
								},
								InitialDelaySeconds: r.Deployment.LivenessProbe.InitialDelaySeconds,
								PeriodSeconds:    r.Deployment.LivenessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.LivenessProbe.FailureThreshold,

							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:cmdstring,
									},
								},

								InitialDelaySeconds: r.Deployment.ReadinessProbe.InitialDelaySeconds,
								PeriodSeconds:  r.Deployment.ReadinessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.ReadinessProbe.FailureThreshold,

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
							//		corev1.ResourceCPU:y.Reqcpu,
							//		corev1.ResourceMemory:y.Reqmem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:y.Limitcpu,
							//		corev1.ResourceMemory:y.Limitmem,
							//	},
							//},

							ImagePullPolicy:corev1.PullPolicy(r.Deployment.ImagePullPolicy),

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

							VolumeMounts: dashVolumeMount,

						},
					},

					// ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: dashVolume,

				},
			},
		},
	}
}

// for handling analytics-worker deployment
func WorkerXDeployment(apimanager *apimv1alpha1.APIManager,r *apimv1alpha1.Profile) *appsv1.Deployment {
	workervolumemounts, workervolume := getWorkerXVolumes(apimanager, *r)


	labels := map[string]string{
		"deployment": r.Name,
	}
	runasuser := int64(802)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:r.Deployment.MinReadySeconds,


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
								InitialDelaySeconds: r.Deployment.LivenessProbe.InitialDelaySeconds,
								PeriodSeconds:     r.Deployment.LivenessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.LivenessProbe.FailureThreshold,
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

								InitialDelaySeconds: r.Deployment.ReadinessProbe.InitialDelaySeconds,
								PeriodSeconds:     r.Deployment.ReadinessProbe.PeriodSeconds,
								FailureThreshold: r.Deployment.ReadinessProbe.FailureThreshold,

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
							//		corev1.ResourceCPU:y.Reqcpu,
							//		corev1.ResourceMemory:y.Reqmem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:y.Limitcpu,
							//		corev1.ResourceMemory:y.Limitmem,
							//	},
							//},

							ImagePullPolicy: corev1.PullPolicy(r.Deployment.ImagePullPolicy),

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

							VolumeMounts: workervolumemounts,
						},
					},
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: workervolume,
				},
			},
		},
	}
}



