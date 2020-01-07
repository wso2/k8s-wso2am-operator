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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Apimanager is a specification for a Apimanager resource
type APIManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIManagerSpec   `json:"spec"`
	Status APIManagerStatus `json:"status"`
}

type PVCServer struct {
	SynapseConfigs string `json:"synapseConfigs"`
	ExecutionPlans string `json:"executionPlans"`
}

//type ApiManagerInstance struct {
//	//ApimDeployment   ApimDeployment   `json:"apim-deployment"`
//	//Deployment     Deployment         `json:"deployment"`
//	//DeploymentConfigmap string `json:"deployment-configmap"`
//	//PersistentVolumeClaimServer PVCServer `json:"persistent-volume-claim-server"`
//}

type Storage struct {
	Memory string `json:"memory"`
	CPU string `json:"cpu"`
}

type Resources struct {
	Requests Storage `json:"requests"`
	Limits Storage `json:"limits"`
}

//type ApimDeployment struct {
//
//}

type Probe struct {
	InitialDelaySeconds int32 		`json:"initialDelaySeconds"`
	FailureThreshold    int32		`json:"failureThreshold"`
	PeriodSeconds		int32		`json:"periodSeconds"`
}

type RollUpdate struct {
	MaxSurge int32 `json:"maxSurge"`
	MaxUnavailable int32 `json:"maxUnavailable"`
}

type Strategy struct {
	RollingUpdate RollUpdate `json:"rollingUpdate"`
}

//type Deployment struct {
//
//
//
//}

//type Profiles struct {
//	ApiManager1    ApiManagerInstance		  `json:"api-manager-1"`
//	ApiManager2    ApiManagerInstance		  `json:"api-manager-2"`
//
//}

type NewConfigmap struct {
	Name string `json:"name"`
	MountPath string `json:"mountPath"`
}

type Configmap struct {
	DeploymentConfigmap string `json:"deploymentConfigmap"`
	NewConfigmap []NewConfigmap `json:"newConfigmap"`

}

type Deployment struct {
	Resources Resources `json:"resources"`
	Replicas *int32 `json:"replicas"`
	MinReadySeconds int32 `json:"minReadySeconds"`
	Strategy Strategy `json:"strategy"`
	ImagePullPolicy string `json:"imagePullPolicy"`
	LivenessProbe Probe `json:"livenessProbe"`
	ReadinessProbe Probe `json:"readinessProbe"`
	Configmaps Configmap `json:"configmaps"`
}

type Service struct {
	Type string `json:"type"`
}

type Profile struct {
	Name  string `json:"name"`
	Deployment Deployment `json:"deployment"`


}

// APIManagerSpec is the spec for a APIManager resource
type APIManagerSpec struct {
	Pattern        string             `json:"pattern"`
	Replicas       *int32             `json:"replicas"`
	//Profiles 		Profiles 		  `json:"profiles"`
	Profiles 	[]Profile  				`json:"profiles"`
	Service Service `json:"service"`

}

// APIManagerStatus is the status for a APIManager resource
type APIManagerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIManagerList is a list of APIManager resources
type APIManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []APIManager `json:"items"`
}
