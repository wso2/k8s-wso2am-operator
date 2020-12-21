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
package patternX

import (
	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getApimXVolumes(apimanager *apimv1alpha1.APIManager, r apimv1alpha1.Profile, x *configvalues) ([]corev1.VolumeMount, []corev1.Volume) {

	var amXvolumemounts []corev1.VolumeMount
	var amXvolume []corev1.Volume
	defaultdeployConf := r.Deployment.Configmaps.DeploymentConfigmap

	// adding default deploymentConfigmap
	amXvolumemounts = append(amXvolumemounts, corev1.VolumeMount{
		Name:      defaultdeployConf,
		MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
		SubPath:   "deployment.toml",
	})

	amXvolume = append(amXvolume, corev1.Volume{
		Name: defaultdeployConf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: defaultdeployConf,
				},
			},
		},
	})

	return amXvolumemounts, amXvolume
}

func getDashboardXVolumes(apimanager *apimv1alpha1.APIManager, r apimv1alpha1.Profile) ([]corev1.VolumeMount, []corev1.Volume) {

	var dashxvolumemounts []corev1.VolumeMount
	var dashxvolume []corev1.Volume

	defaultdashconf := r.Deployment.Configmaps.DeploymentConfigmap

	dashxvolumemounts = append(dashxvolumemounts, corev1.VolumeMount{
		Name:      defaultdashconf,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
		SubPath:   "deployment.yaml",
	})

	dashxvolume = append(dashxvolume, corev1.Volume{
		Name: defaultdashconf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: defaultdashconf,
				},
			},
		},
	})

	return dashxvolumemounts, dashxvolume
}

func getWorkerXVolumes(apimanager *apimv1alpha1.APIManager, r apimv1alpha1.Profile) ([]corev1.VolumeMount, []corev1.Volume) {

	var workerxvolumemounts []corev1.VolumeMount
	var workerxvolume []corev1.Volume

	defaultdeployConf := r.Deployment.Configmaps.DeploymentConfigmap
	// adding default deploymentConfigmap
	workerxvolumemounts = append(workerxvolumemounts, corev1.VolumeMount{
		Name:      defaultdeployConf,
		MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker/deployment.yaml",
		SubPath:   "deployment.yaml",
	})

	workerxvolume = append(workerxvolume, corev1.Volume{
		Name: defaultdeployConf,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: defaultdeployConf,
				},
			},
		},
	})

	return workerxvolumemounts, workerxvolume
}
