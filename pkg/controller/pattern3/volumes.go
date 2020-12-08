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

package pattern3

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getDevVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-devportal-conf"

	var dev1volumemounts []corev1.VolumeMount
	var dev1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-devportal-1" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dev1volumemounts = append(dev1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			dev1volumemounts = append(dev1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dev1volume = append(dev1volume, corev1.Volume{
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
			dev1volume = append(dev1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	//adding default deploymentConfigmap
	dev1volumemounts = append(dev1volumemounts, corev1.VolumeMount{
		Name:      "wso2am-devportal-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	dev1volume = append(dev1volume, corev1.Volume{
		Name: "wso2am-devportal-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return dev1volumemounts, dev1volume

}

func getDev2Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-devportal-conf"

	var dev2volumemounts []corev1.VolumeMount
	var dev2volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-devportal-2" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dev2volumemounts = append(dev2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			dev2volumemounts = append(dev2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			dev2volume = append(dev2volume, corev1.Volume{
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
			dev2volume = append(dev2volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	//adding default deploymentConfigmap
	dev2volumemounts = append(dev2volumemounts, corev1.VolumeMount{
		Name:      "wso2am-devportal-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	dev2volume = append(dev2volume, corev1.Volume{
		Name: "wso2am-devportal-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return dev2volumemounts, dev2volume

}

func getPub1Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-publisher-conf"

	var pub1volumemounts []corev1.VolumeMount
	var pub1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-publisher-1" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			pub1volumemounts = append(pub1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			pub1volumemounts = append(pub1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			pub1volume = append(pub1volume, corev1.Volume{
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
			pub1volume = append(pub1volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	//adding default deploymentConfigmap
	pub1volumemounts = append(pub1volumemounts, corev1.VolumeMount{
		Name:      "wso2-am-publisher-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	pub1volume = append(pub1volume, corev1.Volume{
		Name: "wso2am-publisher-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return pub1volumemounts, pub1volume

}

func getPub2Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-publisher-conf"

	var pub1volumemounts []corev1.VolumeMount
	var pub1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-publisher-2" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			pub2volumemounts = append(pub2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			pub2volumemounts = append(pub2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			pub2volume = append(pub2volume, corev1.Volume{
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
			pub2volume = append(pub2volume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	//adding default deploymentConfigmap
	pub2volumemounts = append(pub2volumemounts, corev1.VolumeMount{
		Name:      "wso2am-publisher-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	pub2volume = append(pub2volume, corev1.Volume{
		Name: "wso2am-publisher-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return pub2volumemounts, pub2volume
}

func getTMVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-tm-conf"

	var tmvolumemounts []corev1.VolumeMount
	var tmvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "traffic-manager" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			tmvolumemounts = append(tmvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}
		//for newly created set of PVCs by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			tmvolumemounts = append(tmvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			tmvolume = append(tmvolume, corev1.Volume{
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
			tmvolume = append(tmvolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	//adding default deploymentConfigmap
	tmvolumemounts = append(tmvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-tm-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	tmvolumemounts = append(tmvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-tm-conf-entrypoint",
		MountPath: "/home/wso2carbon/docker-entrypoint.sh",
		SubPath:   "docker-entrypoint.sh",
	})

	tmvolume = append(tmvolume, corev1.Volume{
		Name: "wso2am-tm-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})
	tmvolume = append(tmvolume, corev1.Volume{
		Name: "wwso2am-tm-conf-entrypoint",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2-am-tm-conf-entrypoint",
				},
			},
		},
	})

	return tmvolumemounts, tmvolume
}

func getKeyManagerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-km-conf"
	var kmvolumemounts []corev1.VolumeMount
	var kmvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-key-manager" {

		deployconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployconfFromYaml != "" {
			deployconfigmap = deployconfFromYaml
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			kmvolumemounts = append(kmvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			kmvolumemounts = append(kmvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			kmvolume = append(kmvolume, corev1.Volume{
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
			kmvolume = append(kmvolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	kmvolumemounts = append(kmvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-km-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	kmvolume = append(kmvolume, corev1.Volume{
		Name: "wso2am-km-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return kmvolumemounts, kmvolume

}

func getGatewayVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-gateway-conf"
	var gwvolumemounts []corev1.VolumeMount
	var gwvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "apim-gateway" {

		gwconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if gwconfFromYaml != "" {
			deployconfigmap = gwconfFromYaml
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			gwvolumemounts = append(gwvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.PersistentVolumeClaim.NewClaim {
			gwvolumemounts = append(gwvolumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			gwvolume = append(gwvolume, corev1.Volume{
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
			gwvolume = append(gwvolume, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: c.Name,
					},
				},
			})
		}
	}

	gwvolumemounts = append(gwvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-gateway-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf",
	})

	gwvolume = append(gwvolume, corev1.Volume{
		Name: "wso2am-gateway-conf",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return gwvolumemounts, gwvolume

}

func getAnalyticsDashVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	dashconfigmap := "wso2-am-analytics-dashboard-conf"
	var dashvolumemounts []corev1.VolumeMount
	var dashvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-dashboard" {

		dashconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if dashconfFromYaml != "" {
			dashconfigmap = dashconfFromYaml
		}

		//for newly created set of configmaps by user

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
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
	}

	dashvolumemounts = append(dashvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-analytics-dashboard-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
		SubPath:   "deployment.yaml",
	})
	dashvolumemounts = append(dashvolumemounts, corev1.VolumeMount{
		Name:      " wso2am-analytics-dashboard-bin",
		MountPath: "/home/wso2carbon/wso2-config-volume/wso2/dashboard/bin/carbon.sh",
		SubPath:   "carbon.sh",
	})

	dashvolume = append(dashvolume, corev1.Volume{
		Name: dashconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: dashconfigmap,
				},
			},
		},
	})

	dashvolume = append(dashvolume, corev1.Volume{
		Name: " wso2am-analytics-dashboard-bin",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2-am-analytics-dashboard-bin",
				},
			},
		},
	})

	return dashvolumemounts, dashvolume

}

func getAnalyticsWorkerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2-am-analytics-worker-conf"

	var workervolumemounts []corev1.VolumeMount
	var workervolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "analytics-worker" {

		workerconfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if workerconfFromYaml != "" {
			deployconfigmap = workerconfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
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
	}

	//adding default deploymentConfigmap
	workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
		Name:      "wso2am-analytics-worker-conf",
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker/deployment.yaml",
		SubPath:   "deployment.yaml",
	})

	workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
		Name:      "wso2am-analytics-worker-bin",
		MountPath: "/home/wso2carbon/wso2-config-volume/wso2/worker/bin/carbon.sh",
		SubPath:   "carbon.sh",
	})

	workervolume = append(workervolume, corev1.Volume{
		Name: deployconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	workervolume = append(workervolume, corev1.Volume{
		Name: "wso2am-analytics-worker-bin",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2-am-analytics-worker-bin",
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
	mysqlvolumemounts = append(mysqlvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-p2-mysql-dbscripts",
		MountPath: "/docker-entrypoint-initdb.d",
	})

	// mysqlvolumemounts = append(mysqlvolumemounts, corev1.VolumeMount{
	// 	Name:      "wso2am-p2-mysql-exec",
	// 	MountPath: "/entrypoint.sh",
	// 	SubPath:   "entrypoint.sh",
	// })
	// mysqlvolumemounts = append(mysqlvolumemounts, corev1.VolumeMount{
	// 	Name:      "apim-rdbms-persistent-storage",
	// 	MountPath: "/var/lib/mysql",
	// })

	var mysqlvolume []corev1.Volume

	mysqlvolume = append(mysqlvolume, corev1.Volume{
		Name: "wso2am-p2-mysql-dbscripts",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2am-p2-mysql-dbscripts",
				},
			},
		},
	})
	// mysqlvolume = append(mysqlvolume, corev1.Volume{
	// 	Name: "wso2am-p2-mysql-exec",
	// 	VolumeSource: corev1.VolumeSource{
	// 		ConfigMap: &corev1.ConfigMapVolumeSource{
	// 			LocalObjectReference: corev1.LocalObjectReference{
	// 				Name: "wso2am-p2-mysql-exec",
	// 			},
	// 		},
	// 	},
	// })

	// mysqlvolume = append(mysqlvolume, corev1.Volume{
	// 	Name: "apim-rdbms-persistent-storage",
	// 	VolumeSource: corev1.VolumeSource{
	// 		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	// 			ClaimName: "wso2am-mysql",
	// 		},
	// 	},
	// })

	return mysqlvolumemounts, mysqlvolume

}
