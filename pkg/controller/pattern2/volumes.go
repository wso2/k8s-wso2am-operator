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

package pattern2

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getDevPubTm1Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p2-am-1-conf"

	var devpubtm1volumemounts []corev1.VolumeMount
	var devpubtm1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-pub-dev-tm-1" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			devpubtm1volumemounts = append(devpubtm1volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			devpubtm1volume = append(devpubtm1volume, corev1.Volume{
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
	devpubtm1volumemounts = append(devpubtm1volumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	devpubtm1volume = append(devpubtm1volume, corev1.Volume{
		Name: deployconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return devpubtm1volumemounts, devpubtm1volume

}

func getDevPubTm2Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p2-am-2-conf"
	var devpubtm2volumemounts []corev1.VolumeMount
	var devpubtm2volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-pub-dev-tm-2" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			devpubtm2volumemounts = append(devpubtm2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			devpubtm2volume = append(devpubtm2volume, corev1.Volume{
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
	devpubtm2volumemounts = append(devpubtm2volumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	devpubtm2volume = append(devpubtm2volume, corev1.Volume{
		Name: deployconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return devpubtm2volumemounts, devpubtm2volume

}

func getKeyManagerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p2-am-km-conf"
	var kmvolumemounts []corev1.VolumeMount
	var kmvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-key-manager" {

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
	}

	kmvolumemounts = append(kmvolumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	kmvolume = append(kmvolume, corev1.Volume{
		Name: deployconfigmap,
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

func getgatewayVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p2-am-gateway-conf"

	var gwvolumemounts []corev1.VolumeMount
	var gwvolume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-gateway" {

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
	}

	gwvolumemounts = append(gwvolumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	gwvolume = append(gwvolume, corev1.Volume{
		Name: deployconfigmap,
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

	dashconfigmap := "wso2am-p2-analytics-dash-conf"
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
		Name:      dashconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
		SubPath:   "deployment.yaml",
	})
	dashvolumemounts = append(dashvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-p2-analytics-bin",
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
		Name: "wso2am-p2-analytics-bin",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2am-p2-analytics-bin",
				},
			},
		},
	})

	return dashvolumemounts, dashvolume

}

func getAnalyticsWorkerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p2-analytics-worker-conf"

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
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker/deployment.yaml",
		SubPath:   "deployment.yaml",
	})

	workervolumemounts = append(workervolumemounts, corev1.VolumeMount{
		Name:      "wso2am-p2-analytics-bin",
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
		Name: "wso2am-p2-analytics-bin",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2am-p2-analytics-bin",
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

	return mysqlvolumemounts, mysqlvolume

}
