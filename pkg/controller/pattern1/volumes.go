package pattern1

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)


func getApim1Volumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	am1conf :=  "wso2am-pattern-1-am-1-conf"
	am1confFromYaml := apimanager.Spec.Profiles[0].Deployment.Configmaps.DeploymentConfigmap
	if am1confFromYaml != ""{
		am1conf = am1confFromYaml
	}
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	synapseConfFromYaml := apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.SynapseConfigs
	if synapseConfFromYaml != "" {
		synapseConf = synapseConfFromYaml
	}
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"
	execPlanFromYaml := apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.ExecutionPlans
	if execPlanFromYaml != "" {
		execPlan = execPlanFromYaml
	}
	//for newly created set of configmaps by user
	var am1volumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.Configmaps.NewConfigmap{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.NewClaim{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             am1conf,
		MountPath:        "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:          "deployment.toml",

	})
	//adding default synapseConfigs pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             synapseConf,
		MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",

	})
	//adding default executionPlans pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             execPlan,
		MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
	})
	//adding docker entrypoint
	//am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
	//	Name:             "wso2am-pattern-1-am-conf-entrypoint",
	//	MountPath:			"/home/wso2carbon/docker-entrypoint.sh",
	//	SubPath:			"docker-entrypoint.sh",
	//})

	var am1volume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.Configmaps.NewConfigmap{
		am1volume =append(am1volume, corev1.Volume{
			Name:         c.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: c.Name,
					},
				},
			},
		})
	}
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.NewClaim{
		am1volume =append(am1volume, corev1.Volume{
			Name:         c.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: c.Name,
					},
				},
			},
		})
	}
	am1volume =append(am1volume,corev1.Volume{
		Name: am1conf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: am1conf,
				},
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: synapseConf,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName:"pvc-synapse-configs",
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: execPlan,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "pvc-execution-plans",
			},
		},
	})
	//defaultmode := int32(0407)
	//am1volume =append(am1volume,corev1.Volume{
	//	Name: "wso2am-pattern-1-am-conf-entrypoint",
	//	VolumeSource: corev1.VolumeSource{
	//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//			LocalObjectReference: corev1.LocalObjectReference{
	//				Name: "wso2am-pattern-1-am-conf-entrypoint",
	//			},
	//			DefaultMode:&defaultmode,
	//		},
	//	},
	//})
	//

	return am1volumemounts, am1volume

}

func getApim2Volumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	am1conf :=  "wso2am-pattern-1-am-2-conf"
	am1confFromYaml := apimanager.Spec.Profiles[0].Deployment.Configmaps.DeploymentConfigmap
	if am1confFromYaml != ""{
		am1conf = am1confFromYaml
	}
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	synapseConfFromYaml := apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.SynapseConfigs
	if synapseConfFromYaml != "" {
		synapseConf = synapseConfFromYaml
	}
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"
	execPlanFromYaml := apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.ExecutionPlans
	if execPlanFromYaml != "" {
		execPlan = execPlanFromYaml
	}
	//for newly created set of configmaps by user
	var am1volumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.Configmaps.NewConfigmap{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.NewClaim{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             am1conf,
		MountPath:        "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:          "deployment.toml",

	})
	//adding default synapseConfigs pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             synapseConf,
		MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",

	})
	//adding default executionPlans pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             execPlan,
		MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
	})
	//adding docker entrypoint
	//am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
	//	Name:             "wso2am-pattern-1-am-conf-entrypoint",
	//	MountPath:			"/home/wso2carbon/docker-entrypoint.sh",
	//	SubPath:			"docker-entrypoint.sh",
	//})

	var am1volume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.Configmaps.NewConfigmap{
		am1volume =append(am1volume, corev1.Volume{
			Name:         c.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: c.Name,
					},
				},
			},
		})
	}
	for _,c:= range apimanager.Spec.Profiles[0].Deployment.PersistentVolumeClaim.NewClaim{
		am1volume =append(am1volume, corev1.Volume{
			Name:         c.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: c.Name,
					},
				},
			},
		})
	}
	am1volume =append(am1volume,corev1.Volume{
		Name: am1conf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: am1conf,
				},
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: synapseConf,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName:"pvc-synapse-configs",
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: execPlan,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "pvc-execution-plans",
			},
		},
	})
	//defaultmode := int32(0407)
	//am1volume =append(am1volume,corev1.Volume{
	//	Name: "wso2am-pattern-1-am-conf-entrypoint",
	//	VolumeSource: corev1.VolumeSource{
	//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//			LocalObjectReference: corev1.LocalObjectReference{
	//				Name: "wso2am-pattern-1-am-conf-entrypoint",
	//			},
	//			DefaultMode:&defaultmode,
	//		},
	//	},
	//})
	//

	return am1volumemounts, am1volume

}
