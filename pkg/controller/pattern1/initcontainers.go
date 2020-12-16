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
	"strconv"

	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// getMysqlInitContainers returns init containers for mysql deployment
func getMysqlInitContainers(apimanager *apimv1alpha1.APIManager, vols *[]corev1.Volume, volMounts *[]corev1.VolumeMount) []corev1.Container {
	var initContainers []corev1.Container

	// UseMysql - default to true
	useMysqlPod := true
	if apimanager.Spec.UseMysql != "" {
		// the error has already
		useMysqlPod, _ = strconv.ParseBool(apimanager.Spec.UseMysql)
	}

	if useMysqlPod {
		// Downloading mysql connector
		// init container
		mysqlConnectorContainer := corev1.Container{}
		mysqlConnectorContainer.Name = "init-mysql-connector-download"
		mysqlConnectorContainer.Image = "busybox:1.32"
		downloadCmdStr := `set -e
              connector_version=8.0.17
              wget https://repo1.maven.org/maven2/mysql/mysql-connector-java/${connector_version}/mysql-connector-java-${connector_version}.jar -P /mysql-connector-jar/`
		mysqlConnectorContainer.Command = []string{"/bin/sh", "-c", downloadCmdStr}
		mysqlConnectorContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "mysql-connector-jar",
				MountPath: "/mysql-connector-jar",
			},
		}
		initContainers = append(initContainers, mysqlConnectorContainer)
		// volume for downloaded mysql connector
		*vols = append(*vols, corev1.Volume{
			Name: "mysql-connector-jar",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// volume mount for downloaded mysql connector
		*volMounts = append(*volMounts, corev1.VolumeMount{
			Name:      "mysql-connector-jar",
			MountPath: "/home/wso2carbon/wso2-artifact-volume/lib",
		})

		// Checking for the availability of MySQL Server deployment
		// init container
		mysqlWaitContainer := corev1.Container{}
		mysqlWaitContainer.Name = "init-mysql"
		mysqlWaitContainer.Image = "busybox:1.32"
		executionStr := "echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"mysql-svc\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";"
		mysqlWaitContainer.Command = []string{"/bin/sh", "-c", executionStr}
		initContainers = append(initContainers, mysqlWaitContainer)
	}

	return initContainers
}

func getAnalyticsWorkerInitContainers() corev1.Container {

	analyticsWorkerInitContainer := corev1.Container{}
	analyticsWorkerInitContainer.Name = "init-am-analytics-worker"
	analyticsWorkerInitContainer.Image = "busybox:1.32"
	analyticsExecutionStr := "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\"; while ! nc -z wso2-am-analytics-worker-svc 7712; do sleep 1; printf \"-\"; done; echo -e \"  >> WSO2 API Manager Analytics Worker has started\";"
	analyticsWorkerInitContainer.Command = []string{"sh", "-c", analyticsExecutionStr}

	return analyticsWorkerInitContainer
}
