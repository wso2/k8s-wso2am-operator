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

package pattern3

import (
	"strconv"
	"strings"

	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// PubDev1Deployment creates a new Deployment for a PubDevTm instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it...
func Pub1Deployment(apimanager *apimv1alpha1.APIManager, x *configvalues, num int) *appsv1.Deployment {
	useMysql := true
	enableAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, _ = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
	}

	labels := map[string]string{
		"deployment": "wso2-am-pubisher",
	}

	pub1VolumeMount, pub1Volume := getPub1Volumes(apimanager, num)
	pub1deployports := getPubContainerPorts()

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &pub1Volume, &pub1VolumeMount)
	}

	if enableAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km"}, &initContainers)
	}

	//initContainers = append(initContainers, getInitContainers([]string{"init-am-analytics-worker"}))
	pub1SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(x.SecurityContext), ":")

	AssignSecurityContext(securityContextString, pub1SecurityContext)

	envVariables := []corev1.EnvVar{
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
			Value: x.JvmMemOpts,
		},
		{
			Name:  "ENABLE_ANALYTICS",
			Value: strconv.FormatBool(enableAnalytics),
		},
		{
			Name:  "PROFILE_NAME",
			Value: "api-publisher",
		},
	}

	for _, env := range x.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-publisher-1-deployment-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:        &x.Replicas,
			MinReadySeconds: x.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
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
							Name:  "wso2am-publisher-1",
							Image: x.Image,
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: x.Livedelay,
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},
								PeriodSeconds:    x.Liveperiod,
								FailureThreshold: x.Livethres,
							},
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: x.Readydelay,
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: cmdstring,
									},
								},

								PeriodSeconds:    x.Readyperiod,
								FailureThreshold: x.Readythres,
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
							SecurityContext: pub1SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(x.Imagepull),
							Ports:           pub1deployports,
							Env:             envVariables,
							VolumeMounts:    pub1VolumeMount,
						},
					},
					ServiceAccountName: x.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: x.ImagePullSecret,
						},
					},

					Volumes: pub1Volume,
				},
			},
		},
	}
}

//pub-2 deployment
func Pub2Deployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {

	useMysql := true
	enableAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, _ = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
	}

	pub2VolumeMount, pub2Volume := getPub2Volumes(apimanager, num)
	pub2deployports := getPubContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-am-publisher",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &pub2Volume, &pub2VolumeMount)
	}

	if enableAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km"}, &initContainers)
	}

	//initContainers = append(initContainers, getAnalyticsWorkerInitContainers([]string{"init-am-analytics-worker"}))
	pub2SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, pub2SecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "ENABLE_ANALYTICS",
			Value: strconv.FormatBool(enableAnalytics),
		},
		{
			Name:  "PROFILE_NAME",
			Value: "api-publisher",
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-publisher-2-deployment-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas:        &z.Replicas,
			MinReadySeconds: z.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
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
							Name:  "wso2am-publisher-2",
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
							SecurityContext: pub2SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           pub2deployports,
							Env:             envVariables,
							VolumeMounts:    pub2VolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: pub2Volume,
				},
			},
		},
	}
}

//devportal-1 deployment
func Devportal1Deployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {

	useMysql := true
	enableAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, _ = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
	}

	dev1VolumeMount, dev1Volume := getDev1Volumes(apimanager, num)
	dev1deployports := getDevportalContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-am-devportal",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &dev1Volume, &dev1VolumeMount)
	}

	if enableAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km"}, &initContainers)
	}

	//initContainers = append(initContainers, getAnalyticsWorkerInitContainers([]string{"init-am-analytics-worker"}))
	dev1SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, dev1SecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "ENABLE_ANALYTICS",
			Value: strconv.FormatBool(enableAnalytics),
		},
		{
			Name:  "PROFILE_NAME",
			Value: "api-devportal",
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-devportal-1-deployment-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas:        &z.Replicas,
			MinReadySeconds: z.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
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
							Name:  "wso2am-devportal-1",
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
							SecurityContext: dev1SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           dev1deployports,
							Env:             envVariables,
							VolumeMounts:    dev1VolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: dev1Volume,
				},
			},
		},
	}
}

//devportal-2 deployment
func Devportal2Deployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {

	useMysql := true
	enableAnalytics := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, _ = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
	}

	dev2VolumeMount, dev2Volume := getDev2Volumes(apimanager, num)
	dev2deployports := getDevportalContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-am-devportal",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := []corev1.Container{}

	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &dev2Volume, &dev2VolumeMount)
	}

	if enableAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km"}, &initContainers)
	}

	//initContainers = append(initContainers, getAnalyticsWorkerInitContainers([]string{"init-am-analytics-worker"}))
	dev2SecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, dev2SecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "ENABLE_ANALYTICS",
			Value: strconv.FormatBool(enableAnalytics),
		},
		{
			Name:  "PROFILE_NAME",
			Value: "api-devportal",
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-devportal-2-deployment-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas:        &z.Replicas,
			MinReadySeconds: z.Minreadysec,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RecreateDeploymentStrategyType),
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
							Name:  "wso2am-devportal-2",
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
							SecurityContext: dev2SecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           dev2deployports,
							Env:             envVariables,
							VolumeMounts:    dev2VolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: dev2Volume,
				},
			},
		},
	}
}

//gateway deployment
func GatewayDeployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.Deployment {
	enableAnalytics := true
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, _ = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
	}

	gatewayVolumeMount, gatewayVolume := getGatewayVolumes(apimanager, num)
	gatewaydeployports := getGatewayContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-gateway",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 8243",
	}

	initContainers := []corev1.Container{}

	if enableAnalytics {
		getInitContainers([]string{"init-apim-analytics", "init-km", "init-tm-1", "init-tm-2"}, &initContainers)
	} else {
		getInitContainers([]string{"init-km", "init-tm-1", "init-tm-2"}, &initContainers)
	}

	gatewaySecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, gatewaySecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "ENABLE_ANALYTICS",
			Value: strconv.FormatBool(enableAnalytics),
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-gateway-deployment-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas:        &z.Replicas,
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
							SecurityContext: gatewaySecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           gatewaydeployports,
							Env:             envVariables,
							VolumeMounts:    gatewayVolumeMount,
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

//km deployment
func KeyManagerDeployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.StatefulSet {

	kmVolumeMount, kmVolume := getKeyManagerVolumes(apimanager, num)
	kmdeployports := getKeyManagerContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-km",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9443",
	}

	initContainers := getMysqlInitContainers(apimanager, &kmVolume, &kmVolumeMount)

	keyManagerSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, keyManagerSecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "PROFILE_NAME",
			Value: "api-key-manager",
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       statefulsetKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-km-statefulset",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.StatefulSetSpec{
			Replicas:    &z.Replicas,
			ServiceName: "wso2-am-km-svc",
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
							Name:  "wso2-am-km",
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
							Env:             envVariables,
							VolumeMounts:    kmVolumeMount,
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

//tm deployment
func TrafficManagerDeployment(apimanager *apimv1alpha1.APIManager, z *configvalues, num int) *appsv1.StatefulSet {

	tmVolumeMount, tmVolume := getTrafficManagerVolumes(apimanager, num)
	tmdeployports := getTmContainerPorts()

	labels := map[string]string{
		"deployment": "wso2-tm",
	}

	cmdstring := []string{
		"/bin/sh",
		"-c",
		"nc -z localhost 9611",
	}

	initContainers := getMysqlInitContainers(apimanager, &tmVolume, &tmVolumeMount)

	tmSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(z.SecurityContext), ":")

	AssignSecurityContext(securityContextString, tmSecurityContext)

	envVariables := []corev1.EnvVar{
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
			Name:  "PROFILE_NAME",
			Value: "traffic-manager",
		},
	}

	for _, env := range z.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       statefulsetKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-tm-statefulset",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},

		Spec: appsv1.StatefulSetSpec{
			Replicas:    &z.Replicas,
			ServiceName: "wso2-am-tm-headless-svc",
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
							Name:  "wso2-am-tm",
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
							SecurityContext: tmSecurityContext,
							ImagePullPolicy: corev1.PullPolicy(z.Imagepull),
							Ports:           tmdeployports,
							Env:             envVariables,
							VolumeMounts:    tmVolumeMount,
						},
					},
					ServiceAccountName: z.ServiceAccountName,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: z.ImagePullSecret,
						},
					},
					Volumes: tmVolume,
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
		"deployment": "wso2-analytics-dashboard",
	}

	dashVolumeMount, dashVolume := getAnalyticsDashVolumes(apimanager, num)

	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &dashVolume, &dashVolumeMount)
	}

	dashboardSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, dashboardSecurityContext)

	envVariables := []corev1.EnvVar{}

	for _, env := range y.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: depAPIVersion,
			Kind:       deploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2-am-analytics-dashboard-deployment-" + apimanager.Name,
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
							Name:  "wso2am-analytics-dashboard",
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
							Env: envVariables,
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
func WorkerDeployment(apimanager *apimv1alpha1.APIManager, y *configvalues, num int) *appsv1.StatefulSet {
	useMysql := true
	if apimanager.Spec.UseMysql != "" {
		useMysql, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}
	workerVolMounts, workerVols := getAnalyticsWorkerVolumes(apimanager, num)

	labels := map[string]string{
		"deployment": "wso2-analytics-worker",
	}

	workerContainerPorts := getWorkerContainerPorts()
	initContainers := []corev1.Container{}
	if useMysql {
		initContainers = getMysqlInitContainers(apimanager, &workerVols, &workerVolMounts)
	}

	workerSecurityContext := &corev1.SecurityContext{}
	securityContextString := strings.Split(strings.TrimSpace(y.SecurityContext), ":")

	AssignSecurityContext(securityContextString, workerSecurityContext)

	envVariables := []corev1.EnvVar{
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
	}

	for _, env := range y.EnvironmentVariables {
		variable := strings.Split(strings.ReplaceAll(strings.TrimSpace(env), " ", ""), ":")
		envObject := corev1.EnvVar{}
		envObject.Name = variable[0]
		envObject.Value = variable[1]
		envVariables = append(envVariables, envObject)
	}

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
			Replicas:    &y.Replicas,
			ServiceName: "wso2-am-analytics-worker-headless-svc",
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
							Env:             envVariables,
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
