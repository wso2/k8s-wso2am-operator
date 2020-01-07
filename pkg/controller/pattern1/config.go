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
	Imagepull   string
	Amimage     string
	Reqcpu      resource.Quantity
	Reqmem      resource.Quantity
	Limitcpu    resource.Quantity
	Limitmem    resource.Quantity
	Am1Configmap string

}


func AssignConfigMapValues(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *configvalues{

	ControlConfigData := configMap.Data

	liveDelay,_ := strconv.ParseInt(ControlConfigData["amProbeInitialDelaySeconds"], 10, 32)
	liveDelayFromYaml := apimanager.Spec.Profiles[0].Deployment.LivenessProbe.InitialDelaySeconds
	if liveDelayFromYaml != 0{
		liveDelay = int64(liveDelayFromYaml)
		//utilruntime.HandleError(fmt.Errorf("Live delay not present"))
	}
	livePeriod,_ := strconv.ParseInt(ControlConfigData["periodSeconds"], 10, 32)
	livePeriodFromYaml := apimanager.Spec.Profiles[0].Deployment.LivenessProbe.PeriodSeconds
	if livePeriodFromYaml != 0{
		livePeriod = int64(livePeriodFromYaml)
	}

	liveThres,_ := strconv.ParseInt(ControlConfigData["failureThreshold"], 10, 32)
	liveThresFromYaml := apimanager.Spec.Profiles[0].Deployment.LivenessProbe.FailureThreshold
	if liveThresFromYaml != 0{
		liveThres = int64(liveThresFromYaml)
	}

	readyDelay,_ := strconv.ParseInt(ControlConfigData["amProbeInitialDelaySeconds"], 10, 32)
	readyDelayFromYaml := apimanager.Spec.Profiles[0].Deployment.ReadinessProbe.InitialDelaySeconds
	if readyDelayFromYaml != 0{
		readyDelay = int64(readyDelayFromYaml)
	}

	readyPeriod,_ := strconv.ParseInt(ControlConfigData["periodSeconds"], 10, 32)
	readyPeriodFromYaml := apimanager.Spec.Profiles[0].Deployment.ReadinessProbe.PeriodSeconds
	if readyPeriodFromYaml != 0{
		readyPeriod = int64(readyPeriodFromYaml)
	}

	readyThres,_ := strconv.ParseInt(ControlConfigData["failureThreshold"], 10, 32)
	readyThresFromYaml := apimanager.Spec.Profiles[0].Deployment.ReadinessProbe.FailureThreshold
	if readyThresFromYaml != 0{
		readyThres = int64(readyThresFromYaml)
	}

	minReadySec,_ := strconv.ParseInt(ControlConfigData["amMinReadySeconds"], 10, 32)
	minReadySecFromYaml := apimanager.Spec.Profiles[0].Deployment.MinReadySeconds
	if minReadySecFromYaml != 0{
		minReadySec = int64(minReadySecFromYaml)
	}

	maxSurges,_ := strconv.ParseInt(ControlConfigData["maxSurge"], 10, 32)
	maxSurgeFromYaml := apimanager.Spec.Profiles[0].Deployment.Strategy.RollingUpdate.MaxSurge
	if maxSurgeFromYaml !=0{
		maxSurges = int64(maxSurgeFromYaml)
	}

	maxUnavail,_ := strconv.ParseInt(ControlConfigData["maxUnavailable"], 10, 32)
	maxUnavailFromYaml := apimanager.Spec.Profiles[0].Deployment.Strategy.RollingUpdate.MaxUnavailable
	if maxUnavailFromYaml !=0{
		maxUnavail = int64(maxUnavailFromYaml)
	}

	imagePull,_ := ControlConfigData["imagePullPolicy"]
	imagePullFromYaml := apimanager.Spec.Profiles[0].Deployment.ImagePullPolicy
	if imagePullFromYaml != ""{
		imagePull = imagePullFromYaml
	}

	amImages:=ControlConfigData["amImage"]


	reqCPU := resource.MustParse(ControlConfigData["amRequestsCPU"])
	reqCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles[0].Deployment.Resources.Requests.CPU)
	lenreqcpu,_ := len(apimanager.Spec.Profiles[0].Deployment.Resources.Requests.CPU),""
	if lenreqcpu != 0{
		reqCPU = reqCPUFromYaml
	}

	reqMem := resource.MustParse(ControlConfigData["amRequestsMemory"])
	reqMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles[0].Deployment.Resources.Requests.Memory)
	lenreqmem,_ := len(apimanager.Spec.Profiles[0].Deployment.Resources.Requests.Memory),""
	if lenreqmem != 0{
		reqMem = reqMemFromYaml
	}

	limitCPU := resource.MustParse(ControlConfigData["amLimitsCPU"])
	limitCPUFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles[0].Deployment.Resources.Limits.CPU)
	lenlimcpu,_ := len(apimanager.Spec.Profiles[0].Deployment.Resources.Limits.CPU),""
	if lenlimcpu != 0{
		limitCPU = limitCPUFromYaml
	}

	limitMem := resource.MustParse(ControlConfigData["amLimitsMemory"])
	limitMemFromYaml,_:=resource.ParseQuantity(apimanager.Spec.Profiles[0].Deployment.Resources.Limits.Memory)
	lenlimmem,_ := len(apimanager.Spec.Profiles[0].Deployment.Resources.Limits.Memory),""
	if lenlimmem != 0{
		limitMem = limitMemFromYaml
	}

	am1ConfigMap := "wso2am-pattern-1-am-1-conf"
	am1ConfigMapFromYaml := apimanager.Spec.Profiles[0].Deployment.Configmaps.DeploymentConfigmap
	if am1ConfigMapFromYaml != ""{
		am1ConfigMap = am1ConfigMapFromYaml
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
		Amimage:     amImages,
		Reqcpu:      reqCPU,
		Reqmem:      reqMem,
		Limitcpu:    limitCPU,
		Limitmem:    limitMem,
		Am1Configmap: am1ConfigMap,

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
			Name:      configMap.Name,
			Namespace: apimanager.Namespace,
			Labels:   labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Data: configMap.Data,
	}
}