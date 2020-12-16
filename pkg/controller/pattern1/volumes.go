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

package pattern1

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getApim1Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p1-apim-1-conf"

	var am1volumemounts []corev1.VolumeMount
	var am1volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-manager-1" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
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
	}

	//adding default deploymentConfigmap
	am1volumemounts = append(am1volumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	am1volume = append(am1volume, corev1.Volume{
		Name: deployconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return am1volumemounts, am1volume

}

func getApim2Volumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p1-apim-2-conf"

	var am2volumemounts []corev1.VolumeMount
	var am2volume []corev1.Volume

	totalProfiles := len(apimanager.Spec.Profiles)

	if totalProfiles > 0 && apimanager.Spec.Profiles[num].Name == "api-manager-2" {

		deployConfFromYaml := apimanager.Spec.Profiles[num].Deployment.Configmaps.DeploymentConfigmap
		if deployConfFromYaml != "" {
			deployconfigmap = deployConfFromYaml
		}

		//for newly created set of configmaps by user
		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am2volumemounts = append(am2volumemounts, corev1.VolumeMount{
				Name:      c.Name,
				MountPath: c.MountPath,
			})
		}

		for _, c := range apimanager.Spec.Profiles[num].Deployment.Configmaps.NewConfigmap {
			am2volume = append(am2volume, corev1.Volume{
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
	am2volumemounts = append(am2volumemounts, corev1.VolumeMount{
		Name:      deployconfigmap,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	am2volume = append(am2volume, corev1.Volume{
		Name: deployconfigmap,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployconfigmap,
				},
			},
		},
	})

	return am2volumemounts, am2volume

}

func getAnalyticsDashVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	dashconfigmap := "wso2am-p1-analytics-conf"
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

	return dashvolumemounts, dashvolume

}

func getAnalyticsWorkerVolumes(apimanager *apimv1alpha1.APIManager, num int) ([]corev1.VolumeMount, []corev1.Volume) {

	deployconfigmap := "wso2am-p1-analytics-conf"

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
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker",
		//SubPath:"deployment.yaml",

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

	return workervolumemounts, workervolume

}

func GetMysqlVolumes(apimanager *apimv1alpha1.APIManager) ([]corev1.VolumeMount, []corev1.Volume) {

	//for newly created set of configmaps by user
	var mysqlvolumemounts []corev1.VolumeMount

	//adding default deploymentConfigmap
	mysqlvolumemounts = append(mysqlvolumemounts, corev1.VolumeMount{
		Name:      "wso2am-p1-mysql-dbscripts",
		MountPath: "/docker-entrypoint-initdb.d",
	})

	var mysqlvolume []corev1.Volume

	mysqlvolume = append(mysqlvolume, corev1.Volume{
		Name: "wso2am-p1-mysql-dbscripts",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wso2am-p1-mysql-dbscripts-" + apimanager.Name,
				},
			},
		},
	})

	return mysqlvolumemounts, mysqlvolume

}
