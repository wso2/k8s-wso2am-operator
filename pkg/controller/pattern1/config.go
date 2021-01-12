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

package pattern1

import (
	"strconv"
	"strings"

	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

type configvalues struct {
	Livedelay            int32
	Liveperiod           int32
	Livethres            int32
	Readydelay           int32
	Readyperiod          int32
	Readythres           int32
	Minreadysec          int32
	Maxsurge             int32
	Maxunavail           int32
	Replicas             int32
	Imagepull            string
	Image                string
	Reqcpu               resource.Quantity
	Reqmem               resource.Quantity
	Limitcpu             resource.Quantity
	Limitmem             resource.Quantity
	APIMVersion          string
	ImagePullSecret      string
	ServiceAccountName   string
	SecurityContext      string
	JvmMemOpts           string
	Enablenalytics       string
	EnvironmentVariables []string
}

type ingressConfigvalues struct {
	Annotations   map[string]string
	Hostname      string
	TransportMode string
	IngressName   string
}

func AssignApimConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data
	totalProfiles := len(apimanager.Spec.Profiles)

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	apimVersion := ControlConfigData["api-manager-version"]
	replicas, _ := strconv.ParseInt(ControlConfigData["apim-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["apim-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["apim-deployment-maxUnavailable"], 10, 32)
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	amImages := ControlConfigData["apim-deployment-image"]
	imagePull, _ := ControlConfigData["apim-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["apim-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	memXmx := ControlConfigData["apim-deployment-resources-jvm-heap-memory-xmx"]
	memXms := ControlConfigData["apim-deployment-resources-jvm-heap-memory-xms"]
	memOpts := "-Xms" + memXms + " -Xmx" + memXmx

	var envVariables []string

	if totalProfiles > 0 {
		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if replicasFromYaml != nil {
			replicas = int64(*replicasFromYaml)
		}

		minReadySecFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if minReadySecFromYaml != 0 {
			minReadySec = int64(minReadySecFromYaml)
		}

		imagePullFromYaml := apimanager.Spec.Profiles[num].Deployment.ImagePullPolicy
		if imagePullFromYaml != "" {
			imagePull = imagePullFromYaml
		}

		reqCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU)
		lenreqcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU), ""
		if lenreqcpu != 0 {
			reqCPU = reqCPUFromYaml
		}

		reqMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory)
		lenreqmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory), ""
		if lenreqmem != 0 {
			reqMem = reqMemFromYaml
		}

		limitCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU)
		lenlimcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU), ""
		if lenlimcpu != 0 {
			limitCPU = limitCPUFromYaml
		}

		limitMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory)
		lenlimmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory), ""
		if lenlimmem != 0 {
			limitMem = limitMemFromYaml
		}

		liveDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.InitialDelaySeconds
		if liveDelayFromYaml != 0 {
			liveDelay = int64(liveDelayFromYaml)
		}
		livePeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.PeriodSeconds
		if livePeriodFromYaml != 0 {
			livePeriod = int64(livePeriodFromYaml)
		}

		liveThresFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.FailureThreshold
		if liveThresFromYaml != 0 {
			liveThres = int64(liveThresFromYaml)
		}

		readyDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.InitialDelaySeconds
		if readyDelayFromYaml != 0 {
			readyDelay = int64(readyDelayFromYaml)
		}

		readyPeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.PeriodSeconds
		if readyPeriodFromYaml != 0 {
			readyPeriod = int64(readyPeriodFromYaml)
		}

		readyThresFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.FailureThreshold
		if readyThresFromYaml != 0 {
			readyThres = int64(readyThresFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}

		envVariablesFromYaml := apimanager.Spec.Profiles[num].Deployment.EnvironmentVariables
		if len(envVariablesFromYaml) > 0 {
			envVariables = envVariablesFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:            int32(liveDelay),
		Liveperiod:           int32(livePeriod),
		Livethres:            int32(liveThres),
		Readydelay:           int32(readyDelay),
		Readyperiod:          int32(readyPeriod),
		Readythres:           int32(readyThres),
		Minreadysec:          int32(minReadySec),
		Maxsurge:             int32(maxSurges),
		Maxunavail:           int32(maxUnavail),
		Imagepull:            imagePull,
		Image:                amImages,
		Reqcpu:               reqCPU,
		Reqmem:               reqMem,
		Limitcpu:             limitCPU,
		Limitmem:             limitMem,
		Replicas:             int32(replicas),
		APIMVersion:          apimVersion,
		ImagePullSecret:      imagePullSecret,
		ServiceAccountName:   serviceAccountName,
		SecurityContext:      securityContext,
		JvmMemOpts:           memOpts,
		EnvironmentVariables: envVariables,
	}

	return cmvalues

}

func AssignApimAnalyticsDashboardConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	replicas, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	amImages := ControlConfigData["apim-analytics-deployment-dashboard-image"]
	imagePull, _ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)

	var envVariables []string

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-dashboard" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if replicasFromYaml != nil {
			replicas = int64(*replicasFromYaml)
		}

		minReadySecFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if minReadySecFromYaml != 0 {
			minReadySec = int64(minReadySecFromYaml)
		}

		imagePullFromYaml := apimanager.Spec.Profiles[num].Deployment.ImagePullPolicy
		if imagePullFromYaml != "" {
			imagePull = imagePullFromYaml
		}

		reqCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU)
		lenreqcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU), ""
		if lenreqcpu != 0 {
			reqCPU = reqCPUFromYaml
		}

		reqMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory)
		lenreqmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory), ""
		if lenreqmem != 0 {
			reqMem = reqMemFromYaml
		}

		limitCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU)
		lenlimcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU), ""
		if lenlimcpu != 0 {
			limitCPU = limitCPUFromYaml
		}

		limitMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory)
		lenlimmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory), ""
		if lenlimmem != 0 {
			limitMem = limitMemFromYaml
		}

		liveDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.InitialDelaySeconds
		if liveDelayFromYaml != 0 {
			liveDelay = int64(liveDelayFromYaml)
		}
		livePeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.PeriodSeconds
		if livePeriodFromYaml != 0 {
			livePeriod = int64(livePeriodFromYaml)
		}

		liveThresFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.FailureThreshold
		if liveThresFromYaml != 0 {
			liveThres = int64(liveThresFromYaml)
		}

		readyDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.InitialDelaySeconds
		if readyDelayFromYaml != 0 {
			readyDelay = int64(readyDelayFromYaml)
		}

		readyPeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.PeriodSeconds
		if readyPeriodFromYaml != 0 {
			readyPeriod = int64(readyPeriodFromYaml)
		}

		readyThresFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.FailureThreshold
		if readyThresFromYaml != 0 {
			readyThres = int64(readyThresFromYaml)
		}

		// Get maxSurge value from the YAML file.
		maxSurgesFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxSurge
		if maxSurgesFromYaml != 0 {
			maxSurges = int64(maxSurgesFromYaml)
		}

		//Get maxUnavailable value from the YAML file.
		maxUnavailFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxUnavail
		if maxUnavailFromYaml != 0 {
			maxUnavail = int64(maxUnavailFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}

		envVariablesFromYaml := apimanager.Spec.Profiles[num].Deployment.EnvironmentVariables
		if len(envVariablesFromYaml) > 0 {
			envVariables = envVariablesFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:            int32(liveDelay),
		Liveperiod:           int32(livePeriod),
		Livethres:            int32(liveThres),
		Readydelay:           int32(readyDelay),
		Readyperiod:          int32(readyPeriod),
		Readythres:           int32(readyThres),
		Minreadysec:          int32(minReadySec),
		Maxsurge:             int32(maxSurges),
		Maxunavail:           int32(maxUnavail),
		Imagepull:            imagePull,
		Image:                amImages,
		Reqcpu:               reqCPU,
		Reqmem:               reqMem,
		Limitcpu:             limitCPU,
		Limitmem:             limitMem,
		Replicas:             int32(replicas),
		ImagePullSecret:      imagePullSecret,
		ServiceAccountName:   serviceAccountName,
		SecurityContext:      securityContext,
		EnvironmentVariables: envVariables,
	}
	return cmvalues

}

func AssignApimAnalyticsWorkerConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	replicas, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	amImages := ControlConfigData["apim-analytics-deployment-worker-image"]
	imagePull, _ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["apim-analytics-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)

	var envVariables []string

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-worker" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if replicasFromYaml != nil {
			replicas = int64(*replicasFromYaml)
		}
		minReadySecFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if minReadySecFromYaml != 0 {
			minReadySec = int64(minReadySecFromYaml)
		}

		imagePullFromYaml := apimanager.Spec.Profiles[num].Deployment.ImagePullPolicy
		if imagePullFromYaml != "" {
			imagePull = imagePullFromYaml
		}

		reqCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU)
		lenreqcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.CPU), ""
		if lenreqcpu != 0 {
			reqCPU = reqCPUFromYaml
		}

		reqMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory)
		lenreqmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Requests.Memory), ""
		if lenreqmem != 0 {
			reqMem = reqMemFromYaml
		}

		limitCPUFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU)
		lenlimcpu, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.CPU), ""
		if lenlimcpu != 0 {
			limitCPU = limitCPUFromYaml
		}

		limitMemFromYaml, _ := resource.ParseQuantity(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory)
		lenlimmem, _ := len(apimanager.Spec.Profiles[num].Deployment.Resources.Limits.Memory), ""
		if lenlimmem != 0 {
			limitMem = limitMemFromYaml
		}

		liveDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.InitialDelaySeconds
		if liveDelayFromYaml != 0 {
			liveDelay = int64(liveDelayFromYaml)
		}
		livePeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.PeriodSeconds
		if livePeriodFromYaml != 0 {
			livePeriod = int64(livePeriodFromYaml)
		}

		liveThresFromYaml := apimanager.Spec.Profiles[num].Deployment.LivenessProbe.FailureThreshold
		if liveThresFromYaml != 0 {
			liveThres = int64(liveThresFromYaml)
		}

		readyDelayFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.InitialDelaySeconds
		if readyDelayFromYaml != 0 {
			readyDelay = int64(readyDelayFromYaml)
		}

		readyPeriodFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.PeriodSeconds
		if readyPeriodFromYaml != 0 {
			readyPeriod = int64(readyPeriodFromYaml)
		}

		readyThresFromYaml := apimanager.Spec.Profiles[num].Deployment.ReadinessProbe.FailureThreshold
		if readyThresFromYaml != 0 {
			readyThres = int64(readyThresFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}

		envVariablesFromYaml := apimanager.Spec.Profiles[num].Deployment.EnvironmentVariables
		if len(envVariablesFromYaml) > 0 {
			envVariables = envVariablesFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:            int32(liveDelay),
		Liveperiod:           int32(livePeriod),
		Livethres:            int32(liveThres),
		Readydelay:           int32(readyDelay),
		Readyperiod:          int32(readyPeriod),
		Readythres:           int32(readyThres),
		Minreadysec:          int32(minReadySec),
		Maxsurge:             int32(maxSurges),
		Maxunavail:           int32(maxUnavail),
		Imagepull:            imagePull,
		Image:                amImages,
		Reqcpu:               reqCPU,
		Reqmem:               reqMem,
		Limitcpu:             limitCPU,
		Limitmem:             limitMem,
		Replicas:             int32(replicas),
		ImagePullSecret:      imagePullSecret,
		ServiceAccountName:   serviceAccountName,
		SecurityContext:      securityContext,
		EnvironmentVariables: envVariables,
	}
	return cmvalues

}

func AssignAPIMIngressConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap) *ingressConfigvalues {
	klog.Info("APIM CONFIGMAPS")
	ControlConfigData := configMap.Data

	annotations := ControlConfigData["ingress.properties"]
	ingressName := ControlConfigData["ingressResourceName"]
	transportMode := ControlConfigData["ingressTransportMode"]
	hostName := ControlConfigData["ingressHostName"]

	annotationsArray := strings.Split(annotations, "\n")
	klog.Info("LENGTH: ", len(annotationsArray))
	annotationsMap := make(map[string]string)
	for _, c := range annotationsArray {
		mapObj := strings.Split(strings.TrimSpace(c), ":")
		klog.Info(len(mapObj))
		if len(mapObj) > 1 {
			annotationsMap[mapObj[0]] = mapObj[1]
		}
	}

	klog.Info("Annotations: ", annotationsMap)
	klog.Info("IngressName: ", ingressName)
	klog.Info("TransportMode: ", transportMode)
	klog.Info("HostName: ", hostName)

	ingressVals := &ingressConfigvalues{
		Annotations:   annotationsMap,
		IngressName:   ingressName,
		TransportMode: transportMode,
		Hostname:      hostName,
	}

	return ingressVals
}

func AssignGatewayIngressConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap) *ingressConfigvalues {

	ControlConfigData := configMap.Data

	annotations := ControlConfigData["ingress.properties"]
	ingressName := ControlConfigData["ingressResourceName"]
	transportMode := ControlConfigData["ingressTransportMode"]
	hostName := ControlConfigData["ingressHostName"]

	annotationsArray := strings.Split(annotations, "\n")
	klog.Info("LENGTH: ", len(annotationsArray))
	annotationsMap := make(map[string]string)
	for _, c := range annotationsArray {
		mapObj := strings.Split(strings.TrimSpace(c), ":")
		klog.Info(len(mapObj))
		if len(mapObj) > 1 {
			annotationsMap[mapObj[0]] = mapObj[1]
		}
	}

	klog.Info("Annotations: ", annotationsMap)
	klog.Info("IngressName: ", ingressName)
	klog.Info("TransportMode: ", transportMode)
	klog.Info("HostName: ", hostName)

	ingressVals := &ingressConfigvalues{
		Annotations:   annotationsMap,
		IngressName:   ingressName,
		TransportMode: transportMode,
		Hostname:      hostName,
	}

	return ingressVals
}

func AssignDashboardIngressConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap) *ingressConfigvalues {

	ControlConfigData := configMap.Data

	annotations := ControlConfigData["ingress.properties"]
	ingressName := ControlConfigData["ingressResourceName"]
	transportMode := ControlConfigData["ingressTransportMode"]
	hostName := ControlConfigData["ingressHostName"]

	annotationsArray := strings.Split(annotations, "\n")
	klog.Info("LENGTH: ", len(annotationsArray))
	annotationsMap := make(map[string]string)
	for _, c := range annotationsArray {
		mapObj := strings.Split(strings.TrimSpace(c), ":")
		klog.Info(len(mapObj))
		if len(mapObj) > 1 {
			annotationsMap[mapObj[0]] = mapObj[1]
		}
	}

	klog.Info("Annotations: ", annotationsMap)
	klog.Info("IngressName: ", ingressName)
	klog.Info("TransportMode: ", transportMode)
	klog.Info("HostName: ", hostName)

	ingressVals := &ingressConfigvalues{
		Annotations:   annotationsMap,
		IngressName:   ingressName,
		TransportMode: transportMode,
		Hostname:      hostName,
	}

	return ingressVals
}

func AssignMysqlConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap) *configvalues {

	ControlConfigData := configMap.Data

	replicas, _ := strconv.ParseInt(ControlConfigData["mysql-replicas"], 10, 32)

	amImages := ControlConfigData["mysql-image"]

	imagePull, _ := ControlConfigData["mysql-imagePullPolicy"]

	cmvalues := &configvalues{

		Imagepull: imagePull,
		Image:     amImages,
		Replicas:  int32(replicas),
	}
	return cmvalues

}

func MakeConfigMap(apimanager *apimv1alpha1.APIManager, configMap *corev1.ConfigMap) *corev1.ConfigMap {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name + "-" + apimanager.Name,
			Namespace: apimanager.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Data: configMap.Data,
	}
}
