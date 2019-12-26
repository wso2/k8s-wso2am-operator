package worker

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	"strconv"

	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// dashboardDeployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func WorkerDeployment(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *appsv1.Deployment {

	ControlConfigData := configMap.Data

	analyticsMinReadySeconds:= "analyticsMinReadySeconds"
	maxSurge:= "maxSurge"
	maxUnavailable:= "maxUnavailable"
	analyticsProbeInitialDelaySeconds := "analyticsProbeInitialDelaySeconds"
	periodSeconds:= "periodSeconds"
	imagePullPolicy:= "imagePullPolicy"
	//analyticsCPU:= "analyticsCPU"
	//analyticsMemory:= "analyticsMemory"


	minReadySec,_ := strconv.ParseInt(ControlConfigData[analyticsMinReadySeconds], 10, 32)
	maxSurges,_ := strconv.ParseInt(ControlConfigData[maxSurge], 10, 32)
	maxUnavail,_ := strconv.ParseInt(ControlConfigData[maxUnavailable], 10, 32)
	liveDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	livePeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	readyDelay,_ := strconv.ParseInt(ControlConfigData[analyticsProbeInitialDelaySeconds], 10, 32)
	readyPeriod,_ := strconv.ParseInt(ControlConfigData[periodSeconds], 10, 32)
	imagePull,_ := ControlConfigData[imagePullPolicy]
	//reqCPU := ControlConfigData[analyticsCPU]
	//reqMem := ControlConfigData[analyticsMemory]
	//limitCPU := ControlConfigData[analyticsCPU]
	//limitMem := ControlConfigData[analyticsMemory]




	//workercpu,_ := resource.ParseQuantity("2000m")
	//workermem,_ := resource.ParseQuantity("4Gi")
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-analytics-worker",
	}
	runasuser := int64(802)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.DeploymentName,
			Name:      "analytics-worker-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:int32(minReadySec),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxSurges),
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(maxUnavail),
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

					//InitContainers: []corev1.Container{
					//	{
					//		Name: "init-apim-analytics-db",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
					//			"sh",
					//			"-c",
					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
					//		},
					//
					//	},
					//	{
					//		Name: "init-am-analytics-worker",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
					//		},
					//
					//	},
					//},
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
								InitialDelaySeconds: int32(liveDelay),
								PeriodSeconds:       int32(livePeriod),
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

							InitialDelaySeconds: int32(readyDelay),
							PeriodSeconds:       int32(readyPeriod),

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
							//		corev1.ResourceCPU:workercpu,
							//		corev1.ResourceMemory:workermem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:workercpu,
							//		corev1.ResourceMemory:workermem,
							//	},
							//},

							ImagePullPolicy: corev1.PullPolicy(imagePull),

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

							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-analytics-worker-conf",
									MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
									//SubPath:"deployment.yaml",
								},

							},
						},
					},

					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-analytics-worker-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "worker-conf",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}










































////package worker
////
////import (
////	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
////	appsv1 "k8s.io/api/apps/v1"
////	corev1 "k8s.io/api/core/v1"
////	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
////	//"k8s.io/apimachinery/pkg/api/resource"
////	"k8s.io/apimachinery/pkg/util/intstr"
////)
////
////// dashboardDeployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
////// the appropriate OwnerReferences on the resource so handleObject can discover
////// the Apimanager resource that 'owns' it.
////func WorkerDeployment(apimanager *apimv1alpha1.APIManager) *appsv1.Deployment {
////	//workercpu,_ := resource.ParseQuantity("2000m")
////	//workermem,_ := resource.ParseQuantity("4Gi")
////	labels := map[string]string{
////		"deployment": "wso2am-pattern-1-analytics-worker",
////	}
////	runasuser := int64(802)
////	return &appsv1.Deployment{
////		ObjectMeta: metav1.ObjectMeta{
////			// Name: apimanager.Spec.DeploymentName,
////			Name:      "analytics-worker-deploy",
////			Namespace: apimanager.Namespace,
////			OwnerReferences: []metav1.OwnerReference{
////				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
////			},
////		},
////		Spec: appsv1.DeploymentSpec{
////			Replicas: apimanager.Spec.Replicas,
////			MinReadySeconds:30,
////			Strategy: appsv1.DeploymentStrategy{
////				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
////				RollingUpdate: &appsv1.RollingUpdateDeployment{
////					MaxSurge: &intstr.IntOrString{
////						Type:   intstr.Int,
////						IntVal: 1,
////					},
////					MaxUnavailable: &intstr.IntOrString{
////						Type:   intstr.Int,
////						IntVal: 0,
////					},
////				},
////			},
////
////			Selector: &metav1.LabelSelector{
////				MatchLabels: labels,
////			},
////			Template: corev1.PodTemplateSpec{
////				ObjectMeta: metav1.ObjectMeta{
////					Labels: labels,
////				},
////				Spec: corev1.PodSpec{
////
////					//InitContainers: []corev1.Container{
////					//	{
////					//		Name: "init-apim-analytics-db",
////					//		Image: "busybox:1.31",
////					//		Command: []string {
////					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
////					//			"sh",
////					//			"-c",
////					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
////					//		},
////					//
////					//	},
////					//	{
////					//		Name: "init-am-analytics-worker",
////					//		Image: "busybox:1.31",
////					//		Command: []string {
////					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
////					//		},
////					//
////					//	},
////					//},
////					Containers: []corev1.Container{
////						{
////							Name:  "wso2am-pattern-1-analytics-worker",
////							Image: "wso2/wso2am-analytics-worker:3.0.0",
////							//LivenessProbe: &corev1.Probe{
////							//	Handler: corev1.Handler{
////							//		Exec:&corev1.ExecAction{
////							//			Command:[]string{
////							//				"/bin/sh",
////							//				"-c",
////							//				"nc -z localhost 9444",
////							//			},
////							//		},
////							//	},
////							//	InitialDelaySeconds: 20,
////							//	PeriodSeconds:       10,
////							//},
////							//ReadinessProbe: &corev1.Probe{
////							//	Handler: corev1.Handler{
////							//		Exec:&corev1.ExecAction{
////							//			Command:[]string{
////							//				"/bin/sh",
////							//				"-c",
////							//				"nc -z localhost 9444",
////							//			},
////							//		},
////							//	},
////							//
////							//	InitialDelaySeconds: 20,
////							//	PeriodSeconds:       10,
////							//
////							//},
////
////							Lifecycle: &corev1.Lifecycle{
////								PreStop:&corev1.Handler{
////									Exec:&corev1.ExecAction{
////										Command:[]string{
////											"sh",
////											"-c",
////											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
////										},
////									},
////								},
////							},
////
////							//Resources:corev1.ResourceRequirements{
////							//	Requests:corev1.ResourceList{
////							//		corev1.ResourceCPU:workercpu,
////							//		corev1.ResourceMemory:workermem,
////							//	},
////							//	Limits:corev1.ResourceList{
////							//		corev1.ResourceCPU:workercpu,
////							//		corev1.ResourceMemory:workermem,
////							//	},
////							//},
////
////							ImagePullPolicy: "Always",
////
////							SecurityContext: &corev1.SecurityContext{
////								RunAsUser:&runasuser,
////							},
////
////							Ports: []corev1.ContainerPort{
////								{
////									ContainerPort: 9764,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 9444,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 7612,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 7712,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 9091,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 7071,
////									Protocol:      "TCP",
////								},
////								{
////									ContainerPort: 7444,
////									Protocol:      "TCP",
////								},
////							},
////
////							VolumeMounts: []corev1.VolumeMount{
////								{
////									Name: "wso2am-pattern-1-am-analytics-worker-conf",
////									MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
////									//SubPath:"deployment.yaml",
////								},
////
////							},
////						},
////					},
////
////					ServiceAccountName: "wso2am-pattern-1-svc-account",
////					ImagePullSecrets:[]corev1.LocalObjectReference{
////						{
////							Name:"wso2am-pattern-1-creds",
////						},
////					},
////
////					Volumes: []corev1.Volume{
////						{
////							Name: "wso2am-pattern-1-am-analytics-worker-conf",
////							VolumeSource: corev1.VolumeSource{
////								ConfigMap: &corev1.ConfigMapVolumeSource{
////									LocalObjectReference: corev1.LocalObjectReference{
////										Name: "worker-conf",
////									},
////								},
////							},
////						},
////					},
////				},
////			},
////		},
////	}
////}
//
//
//
//
//package worker
//
//import (
//	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
//	appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	//"k8s.io/apimachinery/pkg/api/resource"
//	"k8s.io/apimachinery/pkg/util/intstr"
//)
//
//// dashboardDeployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
//// the appropriate OwnerReferences on the resource so handleObject can discover
//// the Apimanager resource that 'owns' it.
//func WorkerDeployment(apimanager *apimv1alpha1.APIManager) *appsv1.Deployment {
//	//workercpu,_ := resource.ParseQuantity("2000m")
//	//workermem,_ := resource.ParseQuantity("4Gi")
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-worker",
//	}
//	runasuser := int64(802)
//	return &appsv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.DeploymentName,
//			Name:      "analytics-worker-deploy",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: appsv1.DeploymentSpec{
//			Replicas: apimanager.Spec.Replicas,
//			MinReadySeconds:30,
//			Strategy: appsv1.DeploymentStrategy{
//				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
//				RollingUpdate: &appsv1.RollingUpdateDeployment{
//					MaxSurge: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 1,
//					},
//					MaxUnavailable: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 0,
//					},
//				},
//			},
//
//			Selector: &metav1.LabelSelector{
//				MatchLabels: labels,
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: labels,
//				},
//				Spec: corev1.PodSpec{
//
//					//InitContainers: []corev1.Container{
//					//	{
//					//		Name: "init-apim-analytics-db",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
//					//			"sh",
//					//			"-c",
//					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
//					//		},
//					//
//					//	},
//					//	{
//					//		Name: "init-am-analytics-worker",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
//					//		},
//					//
//					//	},
//					//},
//					Containers: []corev1.Container{
//						{
//							Name:  "wso2am-pattern-1-analytics-worker",
//							Image: "wso2/wso2am-analytics-worker:3.0.0",
//							//LivenessProbe: &corev1.Probe{
//							//	Handler: corev1.Handler{
//							//		Exec:&corev1.ExecAction{
//							//			Command:[]string{
//							//				"/bin/sh",
//							//				"-c",
//							//				"nc -z localhost 9444",
//							//			},
//							//		},
//							//	},
//							//	InitialDelaySeconds: 20,
//							//	PeriodSeconds:       10,
//							//},
//							//ReadinessProbe: &corev1.Probe{
//							//	Handler: corev1.Handler{
//							//		Exec:&corev1.ExecAction{
//							//			Command:[]string{
//							//				"/bin/sh",
//							//				"-c",
//							//				"nc -z localhost 9444",
//							//			},
//							//		},
//							//	},
//							//
//							//	InitialDelaySeconds: 20,
//							//	PeriodSeconds:       10,
//							//
//							//},
//
//							Lifecycle: &corev1.Lifecycle{
//								PreStop:&corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"sh",
//											"-c",
//											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
//										},
//									},
//								},
//							},
//
//							//Resources:corev1.ResourceRequirements{
//							//	Requests:corev1.ResourceList{
//							//		corev1.ResourceCPU:workercpu,
//							//		corev1.ResourceMemory:workermem,
//							//	},
//							//	Limits:corev1.ResourceList{
//							//		corev1.ResourceCPU:workercpu,
//							//		corev1.ResourceMemory:workermem,
//							//	},
//							//},
//
//							ImagePullPolicy: "Always",
//
//							SecurityContext: &corev1.SecurityContext{
//								RunAsUser:&runasuser,
//							},
//
//							Ports: []corev1.ContainerPort{
//								{
//									ContainerPort: 9764,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9444,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7612,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7712,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9091,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7071,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7444,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7575,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7576,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7577,
//									Protocol:      "TCP",
//								},
//							},
//
//							VolumeMounts: []corev1.VolumeMount{
//								{
//									Name: "wso2am-pattern-1-am-analytics-worker-conf",
//									MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker/deployment.yaml",
//									SubPath:"deployment.yaml",
//								},
//
//							},
//						},
//					},
//
//					ServiceAccountName: "wso2am-pattern-1-svc-account",
//					ImagePullSecrets:[]corev1.LocalObjectReference{
//						{
//							Name:"wso2am-pattern-1-creds",
//						},
//					},
//
//					Volumes: []corev1.Volume{
//						{
//							Name: "wso2am-pattern-1-am-analytics-worker-conf",
//							VolumeSource: corev1.VolumeSource{
//								ConfigMap: &corev1.ConfigMapVolumeSource{
//									LocalObjectReference: corev1.LocalObjectReference{
//										Name: "worker-conf",
//									},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//}