/*
 *
 *  * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
    "strings"
    "fmt"

	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PvcConfig struct {
	Name string
	// NamePrefix   string
	AccessModes        []v1.PersistentVolumeAccessMode
	ClaimSize          string
	StorageClass       *string
	PVCConfigmap       string
	StorageClassEnable string
	// Selector         *metav1.LabelSelector
}

// GetAccessModesFromString returns an array of AccessModes from a string created by GetAccessModesAsString
func GetAccessModesFromString(modes string) []v1.PersistentVolumeAccessMode {
	strmodes := strings.Split(modes, ",")
	accessModes := []v1.PersistentVolumeAccessMode{}
	for _, s := range strmodes {
		s = strings.Trim(s, " ")
		switch {
		case s == "ReadWriteOnce":
			accessModes = append(accessModes, v1.ReadWriteOnce)
		case s == "ReadOnlyMany":
			accessModes = append(accessModes, v1.ReadOnlyMany)
		case s == "ReadWriteMany":
			accessModes = append(accessModes, v1.ReadWriteMany)
		}
	}
	return accessModes
}

func AssignConfigMapValuesForSynapseConfigsPvc(apimanager *apimv1alpha1.APIManager, pvcConfWso2 *v1.ConfigMap) *PvcConfig {

	ControlConfigData := pvcConfWso2.Data

    name, ok := ControlConfigData["wso2amP1AmSynapseConfigsPvcName"]
    if !ok {
        fmt.Println("fyhfryry")
    }
	accessModes, _ := ControlConfigData["wso2amPvcAccessmode"]
	claimSize, _ := ControlConfigData["wso2amPvcSynapseConfigsStorage"]
	storageClass, _ := ControlConfigData["storageClassName"]
	storageClassEnable, _ := ControlConfigData["storageClassEnable"]
	// selector, _ := 	ControlConfigData["synapseConfigsPvcLabel"]
	pvcConfigMap := "pvc-config"

	pvc1cmvalues := &PvcConfig{
		Name:               name,
		AccessModes:        GetAccessModesFromString(accessModes),
		ClaimSize:          claimSize,
		StorageClass:       &storageClass,
		PVCConfigmap:       pvcConfigMap,
		StorageClassEnable: storageClassEnable,
		// Selector: ,
	}
	return pvc1cmvalues
}

func AssignConfigMapValuesForExecutionPlansPvc(apimanager *apimv1alpha1.APIManager, pvcConfWso2 *v1.ConfigMap) *PvcConfig {

	ControlConfigData := pvcConfWso2.Data

	name, _ := ControlConfigData["wso2amP1AmExecutionPlansPvcName"]
	accessModes, _ := ControlConfigData["wso2amPvcAccessmode"]
	claimSize, _ := ControlConfigData["wso2amPvcExecutionPlansStorage"]
	storageClass, _ := ControlConfigData["storageClassName"]
	storageClassEnable, _ := ControlConfigData["storageClassEnable"]
	// selector, _ := 	ControlConfigData["executionPlansPvcLabel"]

	pvcConfigMap := "pvc-config"

	pvc2cmvalues := &PvcConfig{
		Name:               name,
		AccessModes:        GetAccessModesFromString(accessModes),
		ClaimSize:          claimSize,
		StorageClass:       &storageClass,
		PVCConfigmap:       pvcConfigMap,
		StorageClassEnable: storageClassEnable,
		// Selector: ,
	}
	return pvc2cmvalues
}

func AssignConfigMapValuesForMysqlPvc(apimanager *apimv1alpha1.APIManager, pvcConfWso2 *v1.ConfigMap) *PvcConfig {

	ControlConfigData := pvcConfWso2.Data

	name, _ := ControlConfigData["wso2amP1AmMysqlPvcName"]
	accessModes, _ := ControlConfigData["wso2amPvcAccessmode"]
	claimSize, _ := ControlConfigData["wso2amPvcMysqlStorage"]
	storageClass, _ := ControlConfigData["storageClassName"]
	storageClassEnable, _ := ControlConfigData["storageClassEnable"]
	// selector, _ := 	ControlConfigData["mySqlPvcLabel"]

	pvcConfigMap := "pvc-config"

	pvc3cmvalues := &PvcConfig{
		Name:               name,
		AccessModes:        GetAccessModesFromString(accessModes),
		ClaimSize:          claimSize,
		StorageClass:       &storageClass,
		PVCConfigmap:       pvcConfigMap,
		StorageClassEnable: storageClassEnable,
		// Selector: ,
	}
	return pvc3cmvalues
}

func MakeSynapseConfigsPvc(apimanager *apimv1alpha1.APIManager, sconf *PvcConfig) *corev1.PersistentVolumeClaim {

	// labels := map[string]string{
	// 	"name": "wso2am-p1-apim-synapse-configs",
	// }

	if len(sconf.AccessModes) == 0 {
		sconf.AccessModes = append(sconf.AccessModes, v1.ReadWriteMany)
	}

	if len(sconf.ClaimSize) == 0 {
		sconf.ClaimSize = "1Gi"
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sconf.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			// Selector: &metav1.LabelSelector{
			// 	MatchLabels: labels,
			// },
			AccessModes: sconf.AccessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(sconf.ClaimSize),
				},
			},
			StorageClassName: sconf.StorageClass,
		},
	}
}

func MakeExecutionPlansPvc(apimanager *apimv1alpha1.APIManager, epconf *PvcConfig) *corev1.PersistentVolumeClaim {
	// labels := map[string]string{
	// 	"name": "wso2am-p1-apim-execution-plans",
	// }

	if len(epconf.AccessModes) == 0 {
		epconf.AccessModes = append(epconf.AccessModes, v1.ReadWriteMany)
	}

	if len(epconf.ClaimSize) == 0 {
		epconf.ClaimSize = "1Gi"
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: epconf.Name,
			// Name:      "wso2am-p1-am-execution-plans",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			// Selector: &metav1.LabelSelector{
			// 	MatchLabels: labels,
			// },
			AccessModes: epconf.AccessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(epconf.ClaimSize),
				},
			},
			StorageClassName: epconf.StorageClass,
		},
	}
}

func MakeMysqlPvc(apimanager *apimv1alpha1.APIManager, sqlconf *PvcConfig) *corev1.PersistentVolumeClaim {
	// labels := map[string]string{
	// 	"name": "wso2am-p1-apim-mysql",
	// }

	if len(sqlconf.AccessModes) == 0 {
		sqlconf.AccessModes = append(sqlconf.AccessModes, v1.ReadWriteOnce)
	}

	if len(sqlconf.ClaimSize) == 0 {
		sqlconf.ClaimSize = "20Gi"
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: sqlconf.Name,
			// Name:      "wso2am-p1-am-synapse-configs",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			// Selector: &metav1.LabelSelector{
			// 	MatchLabels: labels,
			// },
			AccessModes: sqlconf.AccessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(sqlconf.ClaimSize),
				},
			},
			StorageClassName: sqlconf.StorageClass,
		},
	}
}
