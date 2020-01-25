package pattern1

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)


func getApim1Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	//defaultdeployConf :=  "wso2am-p1-apim-1-conf-"+apimanager.Name
	deployconfigmap := "wso2am-p1-apim-1-conf-"+apimanager.Name
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"

	var am1volumemounts []corev1.VolumeMount
	var am1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles >0 &&  apimanager.Spec.Profiles[num].Name=="api-manager-1"{

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		synapseConfFromYaml := apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.SynapseConfigs
		if synapseConfFromYaml != "" {
			synapseConf = synapseConfFromYaml
		}
		execPlanFromYaml := apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.ExecutionPlans
		if execPlanFromYaml != "" {
			execPlan = execPlanFromYaml
		}
		//for newly created set of configmaps by user

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am1volumemounts = append(am1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			am1volumemounts = append(am1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am1volume = append(am1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			am1volume = append(am1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
	}




	//adding default deploymentConfigmap
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             deployconfigmap,
		MountPath:        "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:          "deployment.toml",

	})
	//adding default synapseConfigs pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             synapseConf,
		MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/synapse-configs",

	})
	//adding default executionPlans pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             execPlan,
		MountPath:"/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/executionplans",
	})

	am1volume =append(am1volume,corev1.Volume{
		Name: deployconfigmap,
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
				ClaimName:"wso2am-p1-am-synapse-configs",
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: execPlan,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "wso2am-p1-am-execution-plans",
			},
		},
	})





	return am1volumemounts, am1volume

}

func getApim2Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	//defaultdeployConf :=  "wso2am-p1-apim-1-conf-"+apimanager.Name
	deployconfigmap := "wso2am-p1-apim-2-conf-"+apimanager.Name
	synapseConf := "wso2am-pattern-1-am-volume-claim-synapse-configs"
	execPlan := "wso2am-pattern-1-am-volume-claim-executionplans"

	var am1volumemounts []corev1.VolumeMount
	var am1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles >0 &&  apimanager.Spec.Profiles[num].Name=="api-manager-2"{

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
			//fmt.Print("aaaaaaaaaaaa11111111 there is a deployment configmap")
		}
		// else {
		//	deployconfigmap = "wso2am-p1-apim-1-conf-"+apimanager.Name
		//	fmt.Print("aaaaaaaaaaaa111111111 there is Noooooo deployment configmap")
		//
		//}
		synapseConfFromYaml := apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.SynapseConfigs
		if synapseConfFromYaml != "" {
			synapseConf = synapseConfFromYaml
		}
		execPlanFromYaml := apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.ExecutionPlans
		if execPlanFromYaml != "" {
			execPlan = execPlanFromYaml
		}
		//for newly created set of configmaps by user

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am1volumemounts = append(am1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			am1volumemounts = append(am1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am1volume = append(am1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			am1volume = append(am1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
	}


	//adding default deploymentConfigmap
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             deployconfigmap,
		MountPath:        "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:          "deployment.toml",

	})
	//adding default synapseConfigs pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             synapseConf,
		MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/synapse-configs",

	})
	//adding default executionPlans pvc
	am1volumemounts=append(am1volumemounts,corev1.VolumeMount{
		Name:             execPlan,
		MountPath:"/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/executionplans",
	})

	am1volume =append(am1volume,corev1.Volume{
		Name: deployconfigmap,
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
				ClaimName:"wso2am-p1-am-synapse-configs",
			},
		},
	})
	am1volume =append(am1volume,corev1.Volume{
		Name: execPlan,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "wso2am-p1-am-execution-plans",
			},
		},
	})





	return am1volumemounts, am1volume

}

func getAnalyticsDashVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {


	dashconfigmap := "wso2am-p1-analytics-dash-conf-"+apimanager.Name
	var dashvolumemounts []corev1.VolumeMount
	var dashvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles >0 && apimanager.Spec.Profiles[num].Name=="analytics-dashboard"{

		dashconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if dashconfFromYaml != "" {
			dashconfigmap = dashconfFromYaml
		}
		//else {
		//	dashconfigmap = "wso2am-p1-analytics-dash-conf-"+apimanager.Name
		//}

		//for newly created set of configmaps by user

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dashvolumemounts = append(dashvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			dashvolumemounts = append(dashvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//adding default deploymentConfigmap

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dashvolume = append(dashvolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			dashvolume = append(dashvolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}

	}



	dashvolumemounts=append(dashvolumemounts,corev1.VolumeMount{
		Name: dashconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
		SubPath:"deployment.yaml",
	})

	dashvolume =append(dashvolume,corev1.Volume{
		Name: dashconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: dashconfigmap,
				},
			},
		},
	})



	return dashvolumemounts, dashvolume

}

func getAnalyticsWorkerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p1-analytics-worker-conf-"+apimanager.Name

	var workervolumemounts []corev1.VolumeMount
	var workervolume []corev1.Volume


	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles >0  && apimanager.Spec.Profiles[num].Name=="analytics-worker"{

		workerconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if workerconfFromYaml != "" {
			deployconfigmap = workerconfFromYaml
		}
		//else {
		//	deployconfigmap = "wso2am-p1-analytics-worker-conf-"+apimanager.Name
		//}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}


		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			workervolume = append(workervolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			workervolume = append(workervolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: c.Name,
						},
					},
				},
			})
		}

	}


	//adding default deploymentConfigmap
	workervolumemounts=append(workervolumemounts,corev1.VolumeMount{
		Name:         deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
		//SubPath:"deployment.yaml",

	})

	workervolume =append(workervolume,corev1.Volume{
		Name: deployconfigmap,
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

func GetMysqlVolumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	//for newly created set of configmaps by user
	var mysqlvolumemounts []corev1.VolumeMount

	//adding default deploymentConfigmap
	mysqlvolumemounts=append(mysqlvolumemounts,corev1.VolumeMount{
		Name: "wso2am-p1-mysql-dbscripts",
		MountPath: "/docker-entrypoint-initdb.d",

	})
	mysqlvolumemounts=append(mysqlvolumemounts,corev1.VolumeMount{
		Name: "apim-rdbms-persistent-storage",
		MountPath: "/var/lib/mysql",

	})

	var mysqlvolume []corev1.Volume

	mysqlvolume =append(mysqlvolume,corev1.Volume{
		Name: "wso2am-p1-mysql-dbscripts",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2am-p1-mysql-dbscripts-"+apimanager.Name,
				},
			},
		},
	})

	mysqlvolume =append(mysqlvolume,corev1.Volume{
		Name: "apim-rdbms-persistent-storage",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName:"wso2am-p1-mysql",
			},
		},
	})

	return mysqlvolumemounts, mysqlvolume

}

