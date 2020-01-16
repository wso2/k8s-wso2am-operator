package pattern1

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
	Image     string
	Reqcpu      resource.Quantity
	Reqmem      resource.Quantity
	Limitcpu    resource.Quantity
	Limitmem    resource.Quantity
}

func AssignApim1ConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["apim-deployment-replicas"], 10, 32)
	replicasFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.MinReadySeconds
	if replicasFromYaml != 0{
		replicas = int64(replicasFromYaml)
	}
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySecondss"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxSurge"], 10, 32)

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxUnavailable"], 10, 32)

	amImages:=ControlConfigData["p1-apim-deployment-image"]

	imagePull,_ := ControlConfigData["apim-deployment-imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}

	reqCPU := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-requests-cpu"])
	reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Requests.CPU)
	lenreqcpu,_ := len(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Requests.CPU),""
	if lenreqcpu != 0{
		reqCPU = reqCPUFromYaml
	}

	reqMem := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-requests-memory"])
	reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Requests.Memory)
	lenreqmem,_ := len(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Requests.Memory),""
	if lenreqmem != 0{
		reqMem = reqMemFromYaml
	}

	limitCPU := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-limits-cpu"])
	limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Limits.CPU)
	lenlimcpu,_ := len(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Limits.CPU),""
	if lenlimcpu != 0{
		limitCPU = limitCPUFromYaml
	}

	limitMem := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-limits-memory"])
	limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Limits.Memory)
	lenlimmem,_ := len(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Limits.Memory),""
	if lenlimmem != 0{
		limitMem = limitMemFromYaml
	}

	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
	}
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(livePeriodFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}

	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
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
	}
	return cmvalues

}

func AssignApim2ConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["apim-deployment-replicas"], 10, 32)
	replicasFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.MinReadySeconds
	if replicasFromYaml != 0{
		replicas = int64(replicasFromYaml)
	}
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-deployment-minReadySecondss"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxSurge"], 10, 32)

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-deployment-maxUnavailable"], 10, 32)

	amImages:=ControlConfigData["p1-apim-deployment-image"]

	imagePull,_ := ControlConfigData["apim-deployment-imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}

	reqCPU := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-requests-cpu"])
	reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Requests.CPU)
	lenreqcpu,_ := len(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Requests.CPU),""
	if lenreqcpu != 0{
		reqCPU = reqCPUFromYaml
	}

	reqMem := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-requests-memory"])
	reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Requests.Memory)
	lenreqmem,_ := len(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Requests.Memory),""
	if lenreqmem != 0{
		reqMem = reqMemFromYaml
	}

	limitCPU := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-limits-cpu"])
	limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager1.Deployment.Resources.Limits.CPU)
	lenlimcpu,_ := len(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Limits.CPU),""
	if lenlimcpu != 0{
		limitCPU = limitCPUFromYaml
	}

	limitMem := resource.MustParse(ControlConfigData["p1-apim-deployment-resources-limits-memory"])
	limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Limits.Memory)
	lenlimmem,_ := len(apimanager.Spec.Profiles.Apimanager2.Deployment.Resources.Limits.Memory),""
	if lenlimmem != 0{
		limitMem = limitMemFromYaml
	}

	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
	}
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(livePeriodFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-livenessProbe-failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}

	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-deployment-readinessProbe-failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
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
	}
	return cmvalues

}

func AssignApimAnalyticsDashboardConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	replicasFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.MinReadySeconds
	if replicasFromYaml != 0{
		replicas = int64(replicasFromYaml)
	}
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)

	amImages:=ControlConfigData["p1-apim-analytics-deployment-image"]

	imagePull,_ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}

	reqCPU := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-requests-cpu"])
	reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Requests.CPU)
	lenreqcpu,_ := len(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Requests.CPU),""
	if lenreqcpu != 0{
		reqCPU = reqCPUFromYaml
	}

	reqMem := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-requests-memory"])
	reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Requests.Memory)
	lenreqmem,_ := len(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Requests.Memory),""
	if lenreqmem != 0{
		reqMem = reqMemFromYaml
	}

	limitCPU := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-limits-cpu"])
	limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Limits.CPU)
	lenlimcpu,_ := len(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Limits.CPU),""
	if lenlimcpu != 0{
		limitCPU = limitCPUFromYaml
	}

	limitMem := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-limits-memory"])
	limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Limits.Memory)
	lenlimmem,_ := len(apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Resources.Limits.Memory),""
	if lenlimmem != 0{
		limitMem = limitMemFromYaml
	}

	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
	}
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(livePeriodFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}

	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
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
	}
	return cmvalues

}

func AssignApimAnalyticsWorkerConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-replicas"], 10, 32)
	replicasFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.MinReadySeconds
	if replicasFromYaml != 0{
		replicas = int64(replicasFromYaml)
	}
	minReadySec,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-minReadySeconds"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}
	maxSurges,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxSurge"], 10, 32)

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-maxUnavailable"], 10, 32)

	amImages:=ControlConfigData["p1-apim-analytics-deployment-image"]

	imagePull,_ := ControlConfigData["apim-analytics-deployment-imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}

	reqCPU := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-requests-cpu"])
	reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Requests.CPU)
	lenreqcpu,_ := len(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Requests.CPU),""
	if lenreqcpu != 0{
		reqCPU = reqCPUFromYaml
	}

	reqMem := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-requests-memory"])
	reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Requests.Memory)
	lenreqmem,_ := len(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Requests.Memory),""
	if lenreqmem != 0{
		reqMem = reqMemFromYaml
	}

	limitCPU := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-limits-cpu"])
	limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Limits.CPU)
	lenlimcpu,_ := len(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Limits.CPU),""
	if lenlimcpu != 0{
		limitCPU = limitCPUFromYaml
	}

	limitMem := resource.MustParse(ControlConfigData["p1-apim-analytics-deployment-resources-limits-memory"])
	limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Limits.Memory)
	lenlimmem,_ := len(apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Resources.Limits.Memory),""
	if lenlimmem != 0{
		limitMem = limitMemFromYaml
	}

	liveDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-initialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
	}
	livePeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(livePeriodFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-livenessProbe-failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}

	readyDelay,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-initialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["apim-analytics-deployment-readinessProbe-failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
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
	}
	return cmvalues

}

func AssignMysqlConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	replicas,_ := strconv.ParseInt(ControlConfigData["p1-mysql-replicas"], 10, 32)

	amImages:=ControlConfigData["p1-mysql-image"]

	imagePull,_ := ControlConfigData["p1-mysql-imagePullPolicy"]

	cmvalues := &configvalues{

		Imagepull:   imagePull,
		Image:     	 amImages,
		Replicas: int32(replicas),
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