package patternX

import (
	v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
)

type configvalues struct {
	Livedelay   int32
	Liveperiod  int32
	Livethres   int32
	Readydelay  int32
	Readyperiod int32
	Readythres  int32
	Minreadysec int32
	Maxsurge    int32
	Maxunavail  int32
	Replicas 	int32
	Imagepull   string
	Image       string
	Reqcpu      resource.Quantity
	Reqmem      resource.Quantity
	Limitcpu    resource.Quantity
	Limitmem    resource.Quantity
	APIMVersion string
	ImagePullSecret string
}

func AssignApimXConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap,  r apimv1alpha1.Profile) *configvalues{

	ControlConfigData := configMap.Data

	apimVersion := ControlConfigData["api-manager-version"]
	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	replicas,_ := strconv.ParseInt(ControlConfigData["apim-deployment-replicas"], 10, 32)
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySeconds"], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxSurge"], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxUnavailable"], 10, 32)
	amImages:=ControlConfigData["pX-apim-deployment-image"]
	imagePull,_ := ControlConfigData["apim-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["pX-apim-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["pX-apim-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["pX-apim-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["pX-apim-deployment-resources-limits-memory"])
	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)

	if r.Deployment.MinReadySeconds !=0{
		minReadySec = int64(r.Deployment.MinReadySeconds)
	}
	if r.Deployment.ImagePullPolicy != ""{
		imagePull = r.Deployment.ImagePullPolicy
	}
	reqMemFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Requests.Memory)
	lenreqmem, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenreqmem !=0{
		reqMem = reqMemFromYaml
	}
	reqCPUFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Requests.CPU)
	lenreqcpu, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenreqcpu !=0{
		reqCPU = reqCPUFromYaml
	}
	limMemFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Limits.Memory)
	lenlimmem, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenlimmem !=0{
		limitMem = limMemFromYaml
	}
	limCPUFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Limits.CPU)
	lenlimcpu, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenlimcpu !=0{
		limitCPU = limCPUFromYaml
	}

	if r.Deployment.LivenessProbe.InitialDelaySeconds !=0{
		liveDelay = int64(r.Deployment.LivenessProbe.InitialDelaySeconds)
	}
	if r.Deployment.LivenessProbe.PeriodSeconds !=0{
		livePeriod = int64(r.Deployment.LivenessProbe.PeriodSeconds)
	}
	if r.Deployment.LivenessProbe.FailureThreshold !=0{
		liveThres = int64(r.Deployment.LivenessProbe.FailureThreshold)
	}
	if r.Deployment.ReadinessProbe.InitialDelaySeconds !=0{
		readyDelay = int64(r.Deployment.LivenessProbe.InitialDelaySeconds)
	}
	if r.Deployment.ReadinessProbe.PeriodSeconds !=0{
		readyPeriod = int64(r.Deployment.LivenessProbe.PeriodSeconds)
	}
	if r.Deployment.ReadinessProbe.FailureThreshold !=0{
		readyThres = int64(r.Deployment.LivenessProbe.FailureThreshold)
	}


	cmvalues := &configvalues{
		Livedelay:   int32(liveDelay),
		Liveperiod:  int32(livePeriod),
		Livethres:   int32(liveThres),
		Readydelay:  int32(readyDelay),
		Readyperiod: int32(readyPeriod),
		Readythres:  int32(readyThres),
		Minreadysec: int32(minReadySec),
		Maxsurge:    int32(maxSurges),
		Maxunavail:  int32(maxUnavail),
		Imagepull:   imagePull,
		Image:     	 amImages,
		Reqcpu:      reqCPU,
		Reqmem:      reqMem,
		Limitcpu:    limitCPU,
		Limitmem:    limitMem,
		Replicas: int32(replicas),
		APIMVersion: apimVersion,
		ImagePullSecret: imagePullSecret,
	}

	return cmvalues

}

func AssignApimAnalyticsConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap,r apimv1alpha1.Profile) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)
	amImages :=""
	if r.Type == "analytics-dashboard"{
		amImages=ControlConfigData["pX-apim-analytics-deployment-image"]
	}
	if r.Type == "analytics-worker"{
		amImages=ControlConfigData["pX-apim-analytics-deployment-image"]
	}

	imagePullSecret := ControlConfigData["image-pull-secret-name"]
	imagePull,_ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	reqCPU := resource.MustParse(ControlConfigData["pX-apim-analytics-deployment-resources-requests-cpu"])
	reqMem := resource.MustParse(ControlConfigData["pX-apim-analytics-deployment-resources-requests-memory"])
	limitCPU := resource.MustParse(ControlConfigData["pX-apim-analytics-deployment-resources-limits-cpu"])
	limitMem := resource.MustParse(ControlConfigData["pX-apim-analytics-deployment-resources-limits-memory"])
	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)


	if r.Deployment.MinReadySeconds !=0{
		minReadySec = int64(r.Deployment.MinReadySeconds)
	}
	if r.Deployment.ImagePullPolicy != ""{
		imagePull = r.Deployment.ImagePullPolicy
	}
	reqMemFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Requests.Memory)
	lenreqmem, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenreqmem !=0{
		reqMem = reqMemFromYaml
	}
	reqCPUFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Requests.CPU)
	lenreqcpu, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenreqcpu !=0{
		reqCPU = reqCPUFromYaml
	}
	limMemFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Limits.Memory)
	lenlimmem, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenlimmem !=0{
		limitMem = limMemFromYaml
	}
	limCPUFromYaml,_ := resource.ParseQuantity(r.Deployment.Resources.Limits.CPU)
	lenlimcpu, _ := len(r.Deployment.Resources.Requests.Memory), ""
	if lenlimcpu !=0{
		limitCPU = limCPUFromYaml
	}

	if r.Deployment.LivenessProbe.InitialDelaySeconds !=0{
		liveDelay = int64(r.Deployment.LivenessProbe.InitialDelaySeconds)
	}
	if r.Deployment.LivenessProbe.PeriodSeconds !=0{
		livePeriod = int64(r.Deployment.LivenessProbe.PeriodSeconds)
	}
	if r.Deployment.LivenessProbe.FailureThreshold !=0{
		liveThres = int64(r.Deployment.LivenessProbe.FailureThreshold)
	}
	if r.Deployment.ReadinessProbe.InitialDelaySeconds !=0{
		readyDelay = int64(r.Deployment.LivenessProbe.InitialDelaySeconds)
	}
	if r.Deployment.ReadinessProbe.PeriodSeconds !=0{
		readyPeriod = int64(r.Deployment.LivenessProbe.PeriodSeconds)
	}
	if r.Deployment.ReadinessProbe.FailureThreshold !=0{
		readyThres = int64(r.Deployment.LivenessProbe.FailureThreshold)
	}


	cmvalues := &configvalues{
		Livedelay:   int32(liveDelay),
		Liveperiod:  int32(livePeriod),
		Livethres:   int32(liveThres),
		Readydelay:  int32(readyDelay),
		Readyperiod: int32(readyPeriod),
		Readythres:  int32(readyThres),
		Minreadysec: int32(minReadySec),
		Maxsurge:    int32(maxSurges),
		Maxunavail:  int32(maxUnavail),
		Imagepull:   imagePull,
		Image:     	 amImages,
		Reqcpu:      reqCPU,
		Reqmem:      reqMem,
		Limitcpu:    limitCPU,
		Limitmem:    limitMem,
		Replicas: int32(replicas),
		ImagePullSecret: imagePullSecret,
	}
	return cmvalues

}



func MakeConfigMap(apimanager *apimv1alpha1.APIManager, configMap *corev1.ConfigMap) *corev1.ConfigMap {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap.Name+"-"+apimanager.Name,
			Namespace: apimanager.Namespace,
			Labels:   labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Data: configMap.Data,
	}
}