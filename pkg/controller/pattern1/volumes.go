package pattern1

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"

)


func getApim1Volumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	defaultdeployConf :=  "wso2am-pattern-1-am-1-conf"
	deployconfigmap := "wso2am-pattern-1-am-1-conf"
	deployConfFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.Configmaps.DeploymentConfigmap
	if deployConfFromYaml != ""{
		deployconfigmap = deployConfFromYaml
	}
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	synapseConfFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.PersistentVolumeClaim.SynapseConfigs
	if synapseConfFromYaml != "" {
		synapseConf = synapseConfFromYaml
	}
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"
	execPlanFromYaml := apimanager.Spec.Profiles.Apimanager1.Deployment.PersistentVolumeClaim.ExecutionPlans
	if execPlanFromYaml != "" {
		execPlan = execPlanFromYaml
	}
	//for newly created set of configmaps by user
	var am1volumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles.Apimanager1.Deployment.Configmaps.NewConfigmap{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles.Apimanager1.Deployment.PersistentVolumeClaim.NewClaim{
		am1volumemounts =append(am1volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             defaultdeployConf,
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


	var am1volume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles.Apimanager1.Deployment.Configmaps.NewConfigmap{
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
	for _,c:= range apimanager.Spec.Profiles.Apimanager1.Deployment.PersistentVolumeClaim.NewClaim{
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
		Name: defaultdeployConf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
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

	return am1volumemounts, am1volume

}

func getApim2Volumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	defaultdeployConf :=  "wso2am-pattern-1-am-2-conf"
	deployconfigmap := "wso2am-pattern-1-am-2-conf"
	deployConfFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.Configmaps.DeploymentConfigmap
	if deployConfFromYaml != ""{
		deployconfigmap = deployConfFromYaml
	}
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	synapseConfFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.PersistentVolumeClaim.SynapseConfigs
	if synapseConfFromYaml != "" {
		synapseConf = synapseConfFromYaml
	}
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"
	execPlanFromYaml := apimanager.Spec.Profiles.Apimanager2.Deployment.PersistentVolumeClaim.ExecutionPlans
	if execPlanFromYaml != "" {
		execPlan = execPlanFromYaml
	}
	//for newly created set of configmaps by user
	var am2volumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles.Apimanager2.Deployment.Configmaps.NewConfigmap{
		am2volumemounts =append(am2volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles.Apimanager2.Deployment.PersistentVolumeClaim.NewClaim{
		am2volumemounts =append(am2volumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	am2volumemounts=append(am2volumemounts,corev1.VolumeMount{
		Name:             defaultdeployConf,
		MountPath:        "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:          "deployment.toml",

	})
	//adding default synapseConfigs pvc
	am2volumemounts=append(am2volumemounts,corev1.VolumeMount{
		Name:             synapseConf,
		MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",

	})
	//adding default executionPlans pvc
	am2volumemounts=append(am2volumemounts,corev1.VolumeMount{
		Name:             execPlan,
		MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
	})


	var am2volume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles.Apimanager2.Deployment.Configmaps.NewConfigmap{
		am2volume =append(am2volume, corev1.Volume{
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
	for _,c:= range apimanager.Spec.Profiles.Apimanager2.Deployment.PersistentVolumeClaim.NewClaim{
		am2volume =append(am2volume, corev1.Volume{
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
	am2volume =append(am2volume,corev1.Volume{
		Name: defaultdeployConf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})
	am2volume =append(am2volume,corev1.Volume{
		Name: synapseConf,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName:"pvc-synapse-configs",
			},
		},
	})
	am2volume =append(am2volume,corev1.Volume{
		Name: execPlan,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "pvc-execution-plans",
			},
		},
	})

	return am2volumemounts, am2volume

}

func getAnalyticsDashVolumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	defaultdashconf :=  "wso2am-pattern-1-am-analytics-dashboard-conf"
	dashconfigmap := "dash-conf"
	dashconfFromYaml := apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Configmaps.DeploymentConfigmap
	if dashconfFromYaml != ""{
		dashconfigmap = dashconfFromYaml
	}

	//for newly created set of configmaps by user
	var dashvolumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Configmaps.NewConfigmap{
		dashvolumemounts =append(dashvolumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.PersistentVolumeClaim.NewClaim{
		dashvolumemounts =append(dashvolumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	dashvolumemounts=append(dashvolumemounts,corev1.VolumeMount{
		Name: defaultdashconf,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
		SubPath:"deployment.yaml",


	})

	////adding docker entrypoint
	//dashvolumemounts=append(am1volumemounts,corev1.VolumeMount{
	//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
	//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
	//	SubPath:"docker-entrypoint.sh",
	//})

	var dashvolume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.Configmaps.NewConfigmap{
		dashvolume =append(dashvolume, corev1.Volume{
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
	for _,c:= range apimanager.Spec.Profiles.AnalyticsDashboard.Deployment.PersistentVolumeClaim.NewClaim{
		dashvolume =append(dashvolume, corev1.Volume{
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
	dashvolume =append(dashvolume,corev1.Volume{
		Name: defaultdashconf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: dashconfigmap,
				},
			},
		},
	})

	//defaultmode := int32(0407)
	//dashvolume =append(dashvolume,corev1.Volume{
	//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
	//	VolumeSource: corev1.VolumeSource{
	//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//			LocalObjectReference: corev1.LocalObjectReference{
	//				Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
	//			},
	//			DefaultMode:&defaultmode,
	//		},
	//	},
	//})


	return dashvolumemounts, dashvolume

}

func getAnalyticsWorkerVolumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	defaultdeployConf :=  "worker-conf"
	deployconfigmap := "worker-conf"
	workerconfFromYaml := apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Configmaps.DeploymentConfigmap
	if workerconfFromYaml != ""{
		deployconfigmap = workerconfFromYaml
	}

	//for newly created set of configmaps by user
	var workervolumemounts []corev1.VolumeMount
	for _,c:= range apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Configmaps.NewConfigmap{
		workervolumemounts =append(workervolumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//for newly created set of PVCs by user
	for _,c:= range apimanager.Spec.Profiles.AnalyticsWorker.Deployment.PersistentVolumeClaim.NewClaim{
		workervolumemounts =append(workervolumemounts, corev1.VolumeMount{
			Name:             c.Name,
			MountPath: 		  c.MountPath,
		})
	}
	//adding default deploymentConfigmap
	workervolumemounts=append(workervolumemounts,corev1.VolumeMount{
		Name:             defaultdeployConf,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
		//SubPath:"deployment.yaml",

	})

	//adding docker entrypoint
	//am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
	//	Name:             "wso2am-pattern-1-am-conf-entrypoint",
	//	MountPath:			"/home/wso2carbon/docker-entrypoint.sh",
	//	SubPath:			"docker-entrypoint.sh",
	//})

	var workervolume []corev1.Volume
	for _,c:= range apimanager.Spec.Profiles.AnalyticsWorker.Deployment.Configmaps.NewConfigmap{
		workervolume =append(workervolume, corev1.Volume{
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
	for _,c:= range apimanager.Spec.Profiles.AnalyticsWorker.Deployment.PersistentVolumeClaim.NewClaim{
		workervolume =append(workervolume, corev1.Volume{
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
	workervolume =append(workervolume,corev1.Volume{
		Name: defaultdeployConf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})



	return workervolumemounts, workervolume

}

func getMysqlVolumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	//for newly created set of configmaps by user
	var mysqlvolumemounts []corev1.VolumeMount

	//adding default deploymentConfigmap
	mysqlvolumemounts=append(mysqlvolumemounts,corev1.VolumeMount{
		Name: "mysql-dbscripts",
		MountPath: "/docker-entrypoint-initdb.d",

	})
	mysqlvolumemounts=append(mysqlvolumemounts,corev1.VolumeMount{
		Name: "apim-rdbms-persistent-storage",
		MountPath: "/var/lib/mysql",

	})

	var mysqlvolume []corev1.Volume

	mysqlvolume =append(mysqlvolume,corev1.Volume{
		Name: "mysql-dbscripts",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "mysql-dbscripts",
				},
			},
		},
	})

	mysqlvolume =append(mysqlvolume,corev1.Volume{
		Name: "apim-rdbms-persistent-storage",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName:"pvc-mysql",
			},
		},
	})

	return mysqlvolumemounts, mysqlvolume

}

