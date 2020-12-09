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

package pattern4

import (
	"strconv"

	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

type configvalues struct {
	Livedelay          int32
	Liveperiod         int32
	Livethres          int32
	Readydelay         int32
	Readyperiod        int32
	Readythres         int32
	Minreadysec        int32
	Maxsurge           int32
	Maxunavail         int32
	Replicas           int32
	Imagepull          string
	Image              string
	Reqcpu             resource.Quantity
	Reqmem             resource.Quantity
	Limitcpu           resource.Quantity
	Limitmem           resource.Quantity
	APIMVersion        string
	ImagePullSecret    string
	ServiceAccountName string
	SecurityContext    string
	JvmMemOpts         string
}

//AssignDevPubTmConfigMapValues is to assign config-map values...
func AssignDevPubTmConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data
	totalProfiles := len(apimanager.Spec.Profiles)

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	apimVersion := ControlConfigData["api-manager-version"]
	replicas, _ := strconv.ParseInt(ControlConfigData["apim-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	amImages := ControlConfigData["p2-apim-deployment-image"]
	imagePull, _ := ControlConfigData["apim-deployment-imagePullPolicy"]
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-pubdevtm-deployment-livelinessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-pubdevtm-deployment-livelinessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-pubdevtm-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-pubdevtm-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	memXmx := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xmx"]
	memXms := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xms"]
	memOpts := "-Xms" + memXms + " -Xmx" + memXmx

	if totalProfiles > 0 && (apimanager.Spec.Profiles[num].Name == "api-pub-dev-tm-1" || apimanager.Spec.Profiles[num].Name == "api-pub-dev-tm-2") {
		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if replicasFromYaml != 0 {
			replicas = int64(replicasFromYaml)
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
	}
	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		Replicas:           int32(replicas),
		APIMVersion:        apimVersion,
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
		JvmMemOpts:         memOpts,
	}

	return cmvalues

}

// AssignApimExternalGatewayConfigMapValues is to assign config-map values...
func AssignApimExternalGatewayConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	apimVersion := ControlConfigData["api-manager-version"]
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	replicas, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-maxUnavailable"], 10, 32)
	amImages := ControlConfigData["p2-apim-deployment-image"]
	imagePull, _ := ControlConfigData["apim-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	memXmx := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xmx"]
	memXms := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xms"]
	memOpts := "-Xms" + memXms + " -Xmx" + memXmx

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-gateway" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if *replicasFromYaml != 0 {
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

		maxSurgesFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxSurge
		if maxSurgesFromYaml != 0 {
			maxSurges = int64(maxSurgesFromYaml)
		}

		maxUnavailFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxUnavail
		if maxUnavail != 0 {
			maxUnavail = int64(maxUnavailFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Maxsurge:           int32(maxSurges),
		Maxunavail:         int32(maxUnavail),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		APIMVersion:        apimVersion,
		Replicas:           int32(replicas),
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
		JvmMemOpts:         memOpts,
	}
	return cmvalues

}

// AssignApimInternalGatewayConfigMapValues is to assign config-map values...
func AssignApimInternalGatewayConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	apimVersion := ControlConfigData["api-manager-version"]
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	replicas, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-maxUnavailable"], 10, 32)
	amImages := ControlConfigData["p2-apim-deployment-image"]
	imagePull, _ := ControlConfigData["apim-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-gw-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	memXmx := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xmx"]
	memXms := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xms"]
	memOpts := "-Xms" + memXms + " -Xmx" + memXmx

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-gateway" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if *replicasFromYaml != 0 {
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

		maxSurgesFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxSurge
		if maxSurgesFromYaml != 0 {
			maxSurges = int64(maxSurgesFromYaml)
		}

		maxUnavailFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxUnavail
		if maxUnavail != 0 {
			maxUnavail = int64(maxUnavailFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Maxsurge:           int32(maxSurges),
		Maxunavail:         int32(maxUnavail),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		APIMVersion:        apimVersion,
		Replicas:           int32(replicas),
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
		JvmMemOpts:         memOpts,
	}
	return cmvalues

}

//AssignKeyManagerConfigMapValues is to assign config-map values...
func AssignKeyManagerConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	apimVersion := ControlConfigData["api-manager-version"]
	securityContext := ControlConfigData["apim-deployment-securityContext"]
	replicas, _ := strconv.ParseInt(ControlConfigData["p2-apim-km-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	amImages := ControlConfigData["p2-apim-deployment-image"]
	imagePull, _ := ControlConfigData["apim-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-km-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-km-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-km-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-km-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	memXmx := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xmx"]
	memXms := ControlConfigData["p2-apim-deployment-resources-jvm-heap-memory-xms"]
	memOpts := "-Xms" + memXms + " -Xmx" + memXmx

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-keymanager" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.Replicas
		if *replicasFromYaml != 0 {
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

	}

	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		APIMVersion:        apimVersion,
		Replicas:           int32(replicas),
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
		JvmMemOpts:         memOpts,
	}
	return cmvalues

}

//AssignApimAnalyticsDashboardConfigMapValues is ...
func AssignApimAnalyticsDashboardConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data
	klog.Info("ConfigMap Data: ", ControlConfigData["image-pull-secret-name"])

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	securityContext := ControlConfigData["apim-deployment-securityContext"]
	replicas, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)
	amImages := ControlConfigData["p2-apim-analytics-deployment-dashboard-image"]
	imagePull, _ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-dashboard" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if replicasFromYaml != 0 {
			replicas = int64(replicasFromYaml)
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

		maxSurgesFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxSurge
		if maxSurgesFromYaml != 0 {
			maxSurges = int64(maxSurgesFromYaml)
		}

		maxUnavailFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxUnavail
		if maxUnavailFromYaml != 0 {
			maxUnavail = int64(maxUnavailFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Maxsurge:           int32(maxSurges),
		Maxunavail:         int32(maxUnavail),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		Replicas:           int32(replicas),
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
	}
	return cmvalues

}

//AssignApimAnalyticsWorkerConfigMapValues is to assign config-map...
func AssignApimAnalyticsWorkerConfigMapValues(apimanager *apimv1alpha1.APIManager, configMap *v1.ConfigMap, num int) *configvalues {

	ControlConfigData := configMap.Data

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	serviceAccountName := ControlConfigData["service-account-name"]

	securityContext := ControlConfigData["apim-deployment-securityContext"]
	replicas, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	minReadySec, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	maxSurges, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)
	maxUnavail, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)
	amImages := ControlConfigData["p2-apim-analytics-deployment-worker-image"]
	imagePull, _ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["p2-apim-analytics-deployment-resources-limits-memory"])
	liveDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay, _ := strconv.ParseInt(ControlConfigData["p2-apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod, _ := strconv.ParseInt(ControlConfigData["p2apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres, _ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-worker" {

		replicasFromYaml := apimanager.Spec.Profiles[num].Deployment.MinReadySeconds
		if replicasFromYaml != 0 {
			replicas = int64(replicasFromYaml)
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

		maxSurgeFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxSurge
		if maxSurgeFromYaml != 0 {
			maxSurges = int64(maxSurgeFromYaml)
		}

		maxUnavailFromYaml := apimanager.Spec.Profiles[num].Deployment.Strategy.RollingUpdate.MaxUnavail
		if maxUnavailFromYaml != 0 {
			maxUnavail = int64(maxUnavailFromYaml)
		}

		securityContextFromYaml := apimanager.Spec.Profiles[num].Deployment.SecurityContext
		if securityContextFromYaml != "" {
			securityContext = securityContextFromYaml
		}
	}

	cmvalues := &configvalues{
		Livedelay:          int32(liveDelay),
		Liveperiod:         int32(livePeriod),
		Livethres:          int32(liveThres),
		Readydelay:         int32(readyDelay),
		Readyperiod:        int32(readyPeriod),
		Readythres:         int32(readyThres),
		Minreadysec:        int32(minReadySec),
		Maxsurge:           int32(maxSurges),
		Maxunavail:         int32(maxUnavail),
		Imagepull:          imagePull,
		Image:              amImages,
		Reqcpu:             reqCPU,
		Reqmem:             reqMem,
		Limitcpu:           limitCPU,
		Limitmem:           limitMem,
		Replicas:           int32(replicas),
		ImagePullSecret:    imagePullSecret,
		ServiceAccountName: serviceAccountName,
		SecurityContext:    securityContext,
	}
	return cmvalues

}

//AssignMysqlConfigMapValues is to assign config-map values...
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

//MakeConfigMap is to create config-map...
func MakeConfigMap(apimanager *apimv1alpha1.APIManager, configMap *corev1.ConfigMap) *corev1.ConfigMap {
	// klog.Info("ConfigMap Creating")
	labels := map[string]string{
		"deployment": "wso2am-pattern-2-am",
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name,
			Namespace: apimanager.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("APIManager")),
			},
		},
		Data: configMap.Data,
	}
}
